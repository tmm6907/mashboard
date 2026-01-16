// Package models defines the data structures and types used throughout the application.
package models

import (
	"github.com/tmm6907/mashboard/utils"
)

type Feed struct {
	FeedID        utils.UUID      `db:"feed_id" json:"feedId"`
	Title         string          `db:"title" json:"title"`
	Link          string          `db:"link" json:"link"`
	Image         string          `db:"image" json:"image"`
	AltText       string          `db:"alt_text" json:"altText"`
	MediaType     string          `db:"media_type" json:"mediaType"`
	Categories    string          `db:"categories" json:"categories"`
	Description   string          `db:"description" json:"description"`
	Language      string          `db:"language" json:"language"`
	LastBuildDate utils.Timestamp `db:"last_build_date" json:"lastBuildDate"`
	Followed      bool            `db:"followed" json:"followed"`
	CreatedAt     utils.Timestamp `db:"created_at" json:"createdAt"`
}

type FeedItem struct {
	ID          uint            `db:"id" json:"id"`
	FeedID      utils.UUID      `db:"feed_id" json:"feedId"`
	FeedName    string          `db:"feed_name" json:"feedName"`
	Title       string          `db:"title" json:"title"`
	Link        string          `db:"link" json:"link"`
	Description string          `db:"description" json:"description"`
	Image       string          `db:"image" json:"image"`
	AltText     string          `db:"alt_text" json:"altText"`
	MediaType   string          `db:"media_type" json:"mediaType"`
	Categories  string          `db:"categories" json:"categories"`
	PubDate     utils.Timestamp `db:"pub_date" json:"pubDate"`
	GUID        string          `db:"guid" json:"guid"`
	Saved       bool            `db:"saved" json:"saved"`
	Read        bool            `db:"read" json:"read"`
	CreatedAt   utils.Timestamp `db:"created_at" json:"createdAt"`
}

type Collection struct {
	ID        uint            `db:"id" json:"id"`
	Name      string          `db:"name" json:"name"`
	CreatedAt utils.Timestamp `db:"created_at" json:"createdAt"`
}
