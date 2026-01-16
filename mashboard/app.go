package main

import (
	"context"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tmm6907/mashboard/handler"
	"github.com/tmm6907/mashboard/serializers"
	"github.com/tmm6907/mashboard/worker"
)

// App struct
type App struct {
	ctx     context.Context
	handler handler.AppHandler
	worker  worker.AppWorker
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		handler: handler.NewAppHandler(),
		worker:  worker.NewAppWorker(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	interval := time.Minute * 5
	a.worker.StartRSSFetcher(&interval)
}

/*
* ROUTES
 */

func (a *App) GetFeeds(req serializers.GetFeedsRequest) handler.Response {
	return a.handler.GetFeeds(req)
}

func (a *App) GetFeedItems(req serializers.GetFeedItemsRequest) handler.Response {
	return a.handler.GetFeedItems(req)
}

func (a *App) GetFeedItem(req serializers.GetFeedItemRequest) handler.Response {
	return a.handler.GetFeedItem(req)
}

func (a *App) SearchForFeed(query string) handler.Response {
	return a.handler.SearchForFeed(query)
}

func (a *App) CreateFeed(req serializers.CreateFeedRequest) handler.Response {
	return a.handler.CreateFeed(req)
}

func (a *App) FollowFeed(req serializers.FollowRequest) handler.Response {
	return a.handler.FollowFeed(req)
}

func (a *App) HandleSaveFeedItem(req serializers.HandleSaveFeedItemRequest) handler.Response {
	return a.handler.HandleSaveFeedItem(req)
}

func (a *App) SetFeedItemAsRead(id int) handler.Response {
	return a.handler.SetFeedItemAsRead(id)
}

func (a *App) CreateNewCollection(name string) handler.Response {
	return a.handler.CreateNewCollection(name)
}

func (a *App) GetCollections() handler.Response {
	return a.handler.GetCollections()
}
