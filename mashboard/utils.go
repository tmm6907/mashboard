package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	netURL "net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

type Response struct {
	Err  string `json:"error"`
	Data any    `json:"data"`
}

// Timestamp type wrapping time.Time
type Timestamp time.Time

func (t Timestamp) String() string {
	return string(time.Time(t).Format(internalTimeFormat))
}

// Standard SQLite date-time format
const internalTimeFormat = "2006-01-02 15:04:05"

var externalFormats = []string{
	"2006-01-02 15:04:05",
	"Mon, 02 Jan 2006 15:04:05 GMT",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"2006-01-02T15:04:05-0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	time.RFC3339,
}

func ParseTimeStr(timeStr string) (time.Time, error) {
	timeStr = strings.Replace(timeStr, "+00:00", "+0000", 1)
	for _, format := range externalFormats {
		parsedTime, err := time.Parse(format, timeStr)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid RSS time format: %s", timeStr)
}

// MarshalJSON controls how Timestamp is serialized to JSON
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(internalTimeFormat))
}

// UnmarshalJSON controls how Timestamp is deserialized from JSON
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Unmarshal JSON into a string first
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	// Parse the string into time.Time
	parsedTime, err := time.Parse(internalTimeFormat, str)
	if err != nil {
		return err
	}
	*t = Timestamp(parsedTime)
	return nil
}

// Scan converts database value into time.Time
func (t *Timestamp) Scan(value any) error {
	if value == nil {
		*t = Timestamp(time.Time{}) // Set to zero time
		return nil
	}
	switch v := value.(type) {
	case string:
		parsedTime, err := ParseTimeStr(v)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %v", err)
		}
		*t = Timestamp(parsedTime)
	case []byte:
		parsedTime, err := ParseTimeStr(string(v))
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %v", err)
		}
		*t = Timestamp(parsedTime)
	case time.Time:
		*t = Timestamp(v)
	default:
		return fmt.Errorf("unexpected type for timestamp: %T", value)
	}
	return nil
}

// Value converts Go time.Time into SQLite format
func (t Timestamp) Value() (driver.Value, error) {
	return time.Time(t).Format(internalTimeFormat), nil
}

type UUID []byte

func (u UUID) MarshalJSON() ([]byte, error) {
	id, err := uuid.FromBytes(u)
	if err != nil {
		return nil, err
	}
	return json.Marshal(id.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*u = id[:]
	return nil
}

func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		*u = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UUID: %v", value)
	}
	*u = bytes
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return []byte(u), nil
}

func ExtractChannelID(youtubeURL string) (string, error) {
	resp, err := http.Get(youtubeURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch page: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read page body: %v", err)
	}

	// Regex to find the channel ID in og:url meta tag
	re := regexp.MustCompile(`https://www\.youtube\.com/channel/([A-Za-z0-9_-]+)`)
	match := re.FindStringSubmatch(string(body))

	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("channel ID not found")
}

func IsYoutubeChannelURL(url string) bool {
	return strings.Contains(url, "www.youtube.com")
}

func getFeedFromChannelID(channelID string) string {
	return "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelID
}

func GetYouTubeRSS(channelURL string) (string, error) {
	if strings.Contains(channelURL, "https://www.youtube.com/feeds/") {
		return channelURL, nil
	}
	handleRegex := regexp.MustCompile(`(?i)(?:https?://)?(?:www\.)?youtube\.com/@([a-zA-Z0-9_-]+)(?:/|$)`)
	customRegex := regexp.MustCompile(`(?i)(?:https?://)?(www\.)?youtube\.com/c/[a-zA-Z0-9_-]+`)
	channelIDRegex := regexp.MustCompile(`(?i)(?:https?://)?(www\.)?youtube\.com/channel/[A-Za-z0-9_-]+`)
	userRegex := regexp.MustCompile(`(?i)(?:https?://)?(www\.)?youtube\.com/user/[a-zA-Z0-9_-]+`)
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return "", errors.New("missing api key")
	}
	var channelID string
	switch {
	case handleRegex.MatchString(channelURL):
		//https://www.googleapis.com/youtube/v3/channels?part=id&forUsername=@LegalEagle&key=
		identifier := handleRegex.FindStringSubmatch(channelURL)[1]
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=id&forHandle=@%s&key=%s", identifier, apiKey)
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		data := make(map[string]any)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return "", err
		}

		items, ok := data["items"].([]any)
		if !ok {
			return "", fmt.Errorf("Invalid response format: items is not a list, got %T", data["items"])
		}

		if len(items) != 1 {
			return "", fmt.Errorf("Unknown exception: %v", data)
		}
		item, ok := items[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("Invalid response format: item is not a map, got %T", items[0])
		}

		id, ok := item["id"].(string)
		if !ok {
			return "", fmt.Errorf("Invalid response format: id is not a string, got %T", item["id"])
		}

		channelID = id

	case customRegex.MatchString(channelURL):
		identifier, err := ExtractChannelID(channelURL)
		if err != nil {
			return "", err
		}
		channelID = identifier
	case channelIDRegex.MatchString(channelURL):
		channelID = channelIDRegex.FindStringSubmatch(channelURL)[2]

	case userRegex.MatchString(channelURL):
		identifier := userRegex.FindStringSubmatch(channelURL)[2]
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=id&forUsername=%s&key=%s", identifier, apiKey)
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		data := make(map[string]any)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return "", err
		}

		items := data["items"].([]map[string]any)

		if len(items) != 1 {
			return "", fmt.Errorf("Unknown exception: %v", data)
		}
		channelID = items[0]["id"].(string)

	default:
		return "", fmt.Errorf("Unrecognized channelURL: %s", channelURL)
	}

	if channelID == "" {
		return "", fmt.Errorf("Unrecognized channelURL: %s", channelURL)
	}
	return getFeedFromChannelID(channelID), nil
}

func ValidateURL(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	defer res.Body.Close() // Prevent resource leak

	// Ensure we have a successful response
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return false
	}
	contentType := res.Header.Get("Content-Type")
	validTypes := []string{"application/xml", "text/xml", "application/rss+xml", "application/atom+xml"}

	for _, validType := range validTypes {
		if strings.Contains(contentType, validType) {
			return true
		}
	}
	return false
}

func GetOGImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the HTML
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:
			// End of document
			return "", errors.New("og:image not found")
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" {
				// Check attributes for property="og:image"
				var content string
				for _, attr := range token.Attr {
					if attr.Key == "property" && attr.Val == "og:image" {
						// Once we find property="og:image", get the content attribute
						for _, attr := range token.Attr {
							if attr.Key == "content" {
								content = attr.Val
								return content, nil
							}
						}
					}
				}
			}
		}
	}
}

func IsURL(s string) bool {
	_, err := netURL.Parse(s)
	return err == nil
}
