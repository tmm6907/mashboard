package handler

import (
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/tmm6907/mashboard/models"
	"github.com/tmm6907/mashboard/serializers"
	"github.com/tmm6907/mashboard/utils"
)

func (a *AppHandler) GetFeeds(req serializers.GetFeedsRequest) Response {
	var feeds []models.Feed
	if req.Category != "" {
		if err := a.db.Select(&feeds, "select * from feeds where followed is 1 and category like ? limit 25", req.Category); err != nil {
			return a.SendError(err)
		}
		return a.SendResponse(feeds)
	}
	if err := a.db.Select(&feeds, "select * from feeds where followed is  limit 25"); err != nil {
		return a.SendError(err)
	}

	return a.SendResponse(feeds)
}

type CreateFeedOpts struct {
	Title       string
	Description string
	Language    string
}

func (a *AppHandler) createFeed(url string, opts *CreateFeedOpts) Response {
	var feed models.Feed
	feedURL := url
	feedTitle := ""
	feedDescription := ""
	feedLanguage := ""
	if opts != nil {
		if opts.Title != "" {
			feedTitle = opts.Title
		}
		if opts.Description != "" {
			feedDescription = opts.Description
		}
		if opts.Language != "" {
			feedLanguage = opts.Language
		}
	}
	if err := a.db.Get(&feed, "select * from feeds where link = ?", url); err != nil {
		isYoutube := utils.IsYoutubeChannelURL(url)
		if isYoutube {
			link, err := utils.GetYouTubeRSS(url)
			if err != nil {
				return a.SendError(err)
			}
			feedURL = link
		}
		rssParser := gofeed.NewParser()
		log.Println("Parsing: ", url)
		feedData, err := rssParser.ParseURL(feedURL)
		if err != nil {
			log.Println(err)

			return a.SendError(err)
		}
		if feedData.Image == nil {
			ogImage, err := utils.GetOGImage(feedURL)
			if err == nil {
				feedData.Image = &gofeed.Image{URL: ogImage}
			} else {
				feedData.Image = &gofeed.Image{URL: "not found"}
			}
		}

		if feedTitle == "" {
			feedTitle = feedData.Title
		}
		if feedDescription == "" {
			feedDescription = feedData.Description
		}
		if feedLanguage == "" {
			feedLanguage = feedData.Language
		}

		feedUUID := uuid.New()
		feedID := feedUUID[:]
		if isYoutube {
			if _, err := a.db.Exec("INSERT OR IGNORE INTO feeds(feed_id, title, link, description, language, categories, media_type, followed) VALUES (?, ?, ?, ?, ?, ?, ?, ?);",
				feedID, feedTitle, feedURL, feedDescription, feedLanguage, "youtube", "video", true); err != nil {
				return a.SendError(err)
			}
			return a.SendResponse(map[string]any{
				"title":       feedTitle,
				"description": feedDescription,
				"image":       feedData.Image.URL,
				"items":       feedData.Items,
			})
		}
		if _, err := a.db.Exec("INSERT OR IGNORE INTO feeds(feed_id, title, link, description, language, followed) VALUES (?, ?, ?, ?, ?, ?);",
			feedID, feedTitle, feedURL, feedDescription, feedLanguage, true); err != nil {
			return a.SendError(err)
		}
		return a.SendResponse(map[string]any{
			"title":       feedTitle,
			"description": feedDescription,
			"image":       feedData.Image.URL,
			"items":       feedData.Items,
		})
	}
	return a.SendError(errors.New("feed already exists"))
}

func (a *AppHandler) CreateFeed(req serializers.CreateFeedRequest) Response {
	if !utils.ValidateURL(req.Link) {
		return a.SendError(errors.New("invalid RSS feed link"))
	}
	return a.createFeed(req.Link, &CreateFeedOpts{
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
	})
}

func (a *AppHandler) searchNewFeed(url string) Response {
	if strings.Contains(url, "reddit.com") {
		url = strings.Trim(url, "/")
		url = strings.TrimSpace(url) + "/.rss"
	}
	return a.createFeed(url, nil)
}

func (a *AppHandler) SearchForFeed(query string) Response {
	if utils.IsURL(query) {
		return a.searchNewFeed(query)
	}
	var feeds []models.Feed
	if err := a.db.Select(&feeds, "SELECT * FROM feeds WHERE title LIKE ? or category LIKE ? ", query, query); err != nil {
		return a.SendError(err)
	}
	return a.SendResponse(feeds)
}

func (a *AppHandler) FollowFeed(req serializers.FollowRequest) Response {
	if req.Link == "" {
		return a.SendError(errors.New("url must not be empty"))
	}
	var feed models.Feed
	if followErr := a.db.Get(&feed, "SELECT * FROM feeds WHERE link = ?;", req.Link); followErr != nil {
		feedID := uuid.New()
		if _, err := a.db.Exec("INSERT OR IGNORE INTO feeds (feed_id, title, link, description, followed) VALUES (?, ?, ?, ?, ?);", utils.UUID(feedID[:]), req.Title, req.Link, req.Desc, true); err != nil {
			log.Println(err)
			return a.SendError(err)
		}
	}

	if req.Collection != "" {
		var collectionID int
		if err := a.db.Get(&collectionID, "SELECT id from collections WHERE name = ?", req.Collection); err != nil {
			if _, err = a.db.Exec("INSERT INTO collections (name) VALUES (?);", req.Collection); err != nil {
				return a.SendError(err)
			}
		}
	}
	return a.SendResponse("Success")
}
