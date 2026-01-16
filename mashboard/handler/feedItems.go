package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/tmm6907/mashboard/models"
	"github.com/tmm6907/mashboard/serializers"
)

func (a *AppHandler) GetFeedItems(req serializers.GetFeedItemsRequest) Response {
	var feedItems []models.FeedItem
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
		return a.SendError(err)
	}
	return a.SendResponse(feedItems)
}

func (a *AppHandler) GetFeedItem(req serializers.GetFeedItemRequest) Response {
	var feedItem models.FeedItem
	query := fmt.Sprintf("SELECT * FROM feed_items WHERE id = %d", req.ID)
	if req.Saved {
		query = query + " AND saved = 1;"
	}
	if err := a.db.Get(&feedItem, query); err != nil {
		return a.SendError(err)
	}
	return a.SendResponse(feedItem)
}

func (a *AppHandler) HandleSaveFeedItem(req serializers.HandleSaveFeedItemRequest) Response {
	val := 0
	if req.Val {
		val = 1
	}
	query := fmt.Sprintf("UPDATE feed_items SET saved = %d WHERE id = %d;", val, req.ID)
	log.Println(query)
	if _, err := a.db.Exec(query); err != nil {
		return a.SendError(err)
	}
	return a.SendResponse("Success")
}

func (a *AppHandler) SetFeedItemAsRead(id int) Response {
	query := fmt.Sprintf("UPDATE feed_items SET read = %d WHERE id = %d;", 1, id)
	log.Println(query)
	if _, err := a.db.Exec(query); err != nil {
		return a.SendError(err)
	}
	return a.SendResponse("Success")
}
