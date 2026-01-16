// Package serializers handles encoding and decoding of data for API responses and requests.
package serializers

type FollowRequest struct {
	Link       string `json:"link"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Collection string `json:"collection"`
}
type GetFeedItemsRequest struct {
	Category string `json:"category"`
	Offset   int    `json:"offset"`
	Saved    bool   `json:"saved"`
}
type GetFeedItemRequest struct {
	ID    int  `json:"id"`
	Saved bool `json:"saved"`
}
type HandleSaveFeedItemRequest struct {
	ID  int  `json:"id"`
	Val bool `json:"value"`
}
type GetFeedsRequest struct {
	Offset   int    `json:"offset"`
	Category string `json:"filter"`
}
type CreateFeedRequest struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Language    string `json:"language"`
}
