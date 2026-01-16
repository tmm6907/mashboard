package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
)

type GetFeedsRequest struct {
	Offset   int    `json:"offset"`
	Category string `json:"filter"`
}

func (a *App) GetFeeds(req GetFeedsRequest) Response {
	var feeds []Feed
	if req.Category != "" {
		if err := a.db.Select(&feeds, "select * from feeds where followed is 1 and category like ? limit 25", req.Category); err != nil {
			return Response{err.Error(), nil}
		}
		return Response{Data: feeds}
	}
	if err := a.db.Select(&feeds, "select * from feeds where followed is 1 limit 25"); err != nil {
		return Response{err.Error(), nil}
	}
	return Response{Data: feeds}
}

type CreateFeedRequest struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Language    string `json:"language"`
}
type CreateFeedOpts struct {
	Title       string
	Description string
	Language    string
}

func (a *App) createFeed(url string, opts *CreateFeedOpts) Response {
	var feed Feed
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
		isYoutube := IsYoutubeChannelURL(url)
		if isYoutube {
			link, err := GetYouTubeRSS(url)
			if err != nil {
				return Response{err.Error(), nil}
			}
			feedURL = link
		}
		rssParser := gofeed.NewParser()
		log.Println("Parsing: ", url)
		feedData, err := rssParser.ParseURL(feedURL)
		if err != nil {
			log.Println(err)
			return Response{err.Error(), nil}
		}
		if feedData.Image == nil {
			ogImage, err := GetOGImage(feedURL)
			if err == nil {
				feedData.Image = &gofeed.Image{URL: ogImage}
			} else {
				log.Println(err)
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
				return Response{err.Error(), nil}
			}
			return Response{"", map[string]any{
				"title":       feedTitle,
				"description": feedDescription,
				"image":       feedData.Image.URL,
				"items":       feedData.Items,
			}}
		} else {
			if _, err := a.db.Exec("INSERT OR IGNORE INTO feeds(feed_id, title, link, description, language, followed) VALUES (?, ?, ?, ?, ?, ?);",
				feedID, feedTitle, feedURL, feedDescription, feedLanguage, true); err != nil {
				return Response{err.Error(), nil}
			}
			return Response{"", []map[string]any{
				{
					"title":       feed.Title,
					"description": feedData.Description,
					"image":       feedData.Image,
					"items":       feedData.Items,
				},
			}}

		}
	}
	return Response{"feed already exists", nil}
}

func (a *App) CreateFeed(req CreateFeedRequest) Response {
	if !ValidateURL(req.Link) {
		return Response{"invalid RSS feed link", nil}
	}
	return a.createFeed(req.Link, &CreateFeedOpts{
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
	})
}

func (a *App) searchNewFeed(url string) Response {
	if strings.Contains(url, "reddit.com") {
		url = strings.Trim(url, "/")
		url = strings.TrimSpace(url) + "/.rss"
	}
	return a.createFeed(url, nil)
}

func (a *App) SearchForFeed(query string) Response {
	if IsURL(query) {
		return a.searchNewFeed(query)
	}
	var feeds []Feed
	sql := fmt.Sprintf("SELECT * FROM feeds WHERE title LIKE %s OR category LIKE %s", query, query)
	if err := a.db.Select(&feeds, sql); err != nil {
		return Response{err.Error(), nil}
	}
	return Response{"", feeds}
}

type FollowRequest struct {
	Link       string `json:"link"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Collection string `json:"collection"`
}

func (a *App) FollowFeed(req FollowRequest) Response {
	if req.Link == "" {
		return Response{"url must not be empty", nil}
	}
	var feed Feed
	if followErr := a.db.Get(&feed, "SELECT * FROM feeds WHERE link = ?;", req.Link); followErr != nil {
		feedID := uuid.New()
		if _, err := a.db.Exec("INSERT OR IGNORE INTO feeds (feed_id, title, link, description, followed) VALUES (?, ?, ?, ?, ?);", UUID(feedID[:]), req.Title, req.Link, req.Desc, true); err != nil {
			log.Println(err)
			return Response{err.Error(), nil}
		}
	}

	if req.Collection != "" {
		var collectionID int
		if err := a.db.Get(&collectionID, "SELECT id from collections WHERE name = ?", req.Collection); err != nil {
			if _, err = a.db.Exec("INSERT INTO collections (name) VALUES (?);", req.Collection); err != nil {
				return Response{err.Error(), nil}
			}
		}
	}
	return Response{"", "Success"}
}
