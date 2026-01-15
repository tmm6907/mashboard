package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

func (a *App) fetchRSSFeed(feed Feed) error {
	rssParser := gofeed.NewParser()
	rssFeed, err := rssParser.ParseURL(feed.Link)
	if err != nil {
		log.Println(err)
		return err
	}
	var feedImage, feedAlt, feedCategories, mediaType string

	a.dbMutex.Lock()
	if rssFeed.Image != nil && feed.Image == "" {
		feedImage = rssFeed.Image.URL
		feedAlt = rssFeed.Image.Title
		if _, err := a.db.Exec("UPDATE feeds SET image = ?, alt_text = ? WHERE feed_id = ?;", feedImage, feedAlt, feed.FeedID); err != nil {
			a.dbMutex.Unlock()
			log.Println(err)
			return err
		}
	}
	if len(rssFeed.Categories) > 0 && feed.Categories == "" {
		feedCategories = strings.Join(rssFeed.Categories, ", ")
		if _, err := a.db.Exec("UPDATE feeds SET categories = ? WHERE feed_id = ?;", feedCategories, feed.FeedID); err != nil {
			a.dbMutex.Unlock()
			log.Println(err)
			return err
		}
	}
	a.dbMutex.Unlock()

	if feed.MediaType != "" {
		mediaType = feed.MediaType
	} else if rssFeed.FeedType != "" {
		mediaType = rssFeed.FeedType
	}
	for _, item := range rssFeed.Items {
		image := feedImage
		alt := feedAlt
		categories := ""
		media := ""

		if len(item.Enclosures) > 0 {
			media = item.Enclosures[0].Type
		}
		if item.Image != nil {
			image = item.Image.URL
			alt = item.Image.Title
		} else if item.ITunesExt != nil && item.ITunesExt.Image != "" {
			image = item.ITunesExt.Image
		}
		if len(item.Categories) > 0 {
			categories = strings.Join(item.Categories, ", ")
			if item.ITunesExt != nil && item.ITunesExt.Keywords != "" {
				log.Println(item.ITunesExt.Keywords)
			}
		} else {
			if item.ITunesExt != nil && item.ITunesExt.Keywords != "" {
				log.Println(item.ITunesExt.Keywords)
			}
		}

		if image == "" {
			image, _ = GetOGImage(item.Link)
			if image == "" && feedImage != "" {
				image = feedImage
			}
		}
		if categories == "" && feedCategories != "" {
			categories = feedCategories
		}
		if media == "" && mediaType != "" {
			media = mediaType
		}

		// Lock for the entire check-and-insert/update operation
		a.dbMutex.Lock()

		var feedItem FeedItem
		err := a.db.Get(&feedItem, "SELECT * FROM feed_items WHERE guid = ?;", item.GUID)
		if err != nil {
			// Item doesn't exist, insert it
			pubDate, err := ParseTimeStr(item.Published)
			if err != nil {
				a.dbMutex.Unlock()
				log.Println(err)
				return err
			}
			if _, err = a.db.Exec(
				`
				INSERT OR IGNORE INTO feed_items 
				(feed_id, feed_name, title, link, description, image, alt_text, categories, guid, pub_date, media_type) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
				`,
				feed.FeedID, feed.Title, item.Title, item.Link, item.Description, image, alt, categories, item.GUID, Timestamp(pubDate), media,
			); err != nil {
				a.dbMutex.Unlock()
				log.Println(err)
				return err
			}
		} else {
			// Item exists, update if needed
			if feedItem.Image == "" && image != "" {
				if _, err := a.db.Exec("UPDATE feed_items SET image = ?, alt_text = ? WHERE id = ?", image, alt, feedItem.ID); err != nil {
					a.dbMutex.Unlock()
					log.Println(err)
					return err
				}
			}
			if feedItem.Categories == "" && categories != "" {
				if _, err := a.db.Exec("UPDATE feed_items SET categories = ? WHERE id = ?", categories, feedItem.ID); err != nil {
					a.dbMutex.Unlock()
					log.Println(err)
					return err
				}
			}
		}
		a.dbMutex.Unlock()
	}

	return nil
}

func (a *App) FetchRSSFeeds() {
	var feeds []Feed

	err := a.db.Select(&feeds, "SELECT * FROM feeds;")
	if err != nil {
		log.Println(err)
	}
	workers := 10
	feedChan := make(chan Feed, len(feeds))
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for f := range feedChan {
				log.Println("Fetching RSS Feed for url: ", f.Link, " ...")
				if err := a.fetchRSSFeed(f); err != nil {
					log.Println("Error ", err)
				}
			}
		})
	}
	for _, feed := range feeds {
		feedChan <- feed
	}
	close(feedChan)
	wg.Wait()
}

func (a *App) StartRSSFetcher(interval *time.Duration) {
	a.FetchRSSFeeds()
	if interval == nil {
		defaultDuration := 1 * time.Minute
		interval = &defaultDuration
	}
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for range ticker.C {
		log.Println("Fetching RSS Feeds...")
		a.FetchRSSFeeds()
		log.Println("Finished fetching...")
	}
}
