export type FeedItemData = {
  id: number;
  feedId: string;
  feedName: string;
  title: string;
  link: string;
  description: string;
  image: string;
  altText: string;
  mediaType: string;
  categories: string;
  pubDate: string;
  guid: string;
  saved: boolean;
  createdAt: string;
};

export type FeedsState = {
  items: FeedItemData[];
  currentItem?: FeedItemData;
  offset: number;
  newItems: FeedItemData[];
};

class FeedStore {
  items = $state<FeedItemData[]>([]);
  newItems = $state<FeedItemData[]>([]);
  offset = $state(0);
  currentItem = $state<FeedItemData | undefined>(undefined);
  isLoading = $state(false);
  category = $state("");
  saved = $state(false);

  setFeeds(data: FeedItemData[]) {
    this.items = data ? data : [];
  }

  setNewFeeds(data: FeedItemData[]) {
    this.newItems = data ? data : [];
  }

  addMoreFeeds(data: FeedItemData[]) {
    this.items.push(...data);
  }

  refresh() {
    this.items = [...this.newItems];
    this.newItems = [];
  }
}

export const feed = new FeedStore();
