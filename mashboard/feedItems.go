package main

import (
	"fmt"
	"log"
	"strings"
)

type GetFeedItemsRequest struct {
	Category string `json:"category"`
	Offset   int    `json:"offset"`
	Saved    bool   `json:"saved"`
}

func (a *App) GetFeedItems(req GetFeedItemsRequest) Response {
	var feedItems []FeedItem
	category := strings.ToLower(req.Category)
	query := "SELECT * FROM feed_items"
	if category == "technology" {
		category = "tech"
	}
	if category != "" || req.Saved {
		if category != "" && req.Saved {
			query = query + fmt.Sprintf(" WHERE categories LIKE %s OR media_type LIKE %s and saved = 1", category, category)
		} else if category != "" {
			query = query + " WHERE categories LIKE ? OR media_type LIKE ?"
		} else {
			query = query + " WHERE saved = 1"
		}
	}
	query = fmt.Sprintf("%s ORDER BY pub_date DESC LIMIT 25 OFFSET %d;", query, req.Offset)
	log.Println(query, fmt.Sprintf("%d, %s, %v", req.Offset, req.Category, req.Saved))
	if err := a.db.Select(&feedItems, query); err != nil {
		return Response{err.Error(), nil}
	}
	return Response{"", feedItems}
}

type GetFeedItemRequest struct {
	ID    int  `json:"id"`
	Saved bool `json:"saved"`
}

func (a *App) GetFeedItem(req GetFeedItemRequest) Response {
	var feedItem FeedItem
	query := fmt.Sprintf("SELECT * FROM feed_items WHERE id = %d", req.ID)
	if req.Saved {
		query = query + " AND saved = 1;"
	}
	if err := a.db.Get(&feedItem, query); err != nil {
		return Response{err.Error(), nil}
	}
	return Response{"", feedItem}
}

type HandleSaveFeedItemRequest struct {
	ID  int  `json:"id"`
	Val bool `json:"value"`
}

func (a *App) HandleSaveFeedItem(req HandleSaveFeedItemRequest) Response {
	val := 0
	if req.Val {
		val = 1
	}
	query := fmt.Sprintf("UPDATE feed_items SET saved = %d WHERE id = %d;", val, req.ID)
	log.Println(query)
	if _, err := a.db.Exec(query); err != nil {
		return Response{err.Error(), nil}
	}
	return Response{"", "Success!"}
}
