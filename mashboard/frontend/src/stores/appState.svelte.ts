import type { FeedItemData } from "./feeds.svelte.ts";

/*
""
"feeds/{id}"
"feeds/items/{id}"
"collections/{type}"
"saved"
"explore"
*/
export type DynamicID = number | null;
export type AppState = {
  currentPage: string[];
  collectionType: string;
  dynamicID: DynamicID;
  recentPosts: FeedItemData[];
  showNewCollectionForm: boolean;
};

export const appState: AppState = $state({
  currentPage: ["#/"],
  collectionType: "",
  dynamicID: null,
  recentPosts: [],
  showNewCollectionForm: false,
});
