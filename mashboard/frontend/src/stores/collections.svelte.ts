export interface Collection {
  id: string;
  name: string;
  createdAt: string;
}

export type CollectionsState = {
  items: Collection[];
};

export const collections: CollectionsState = $state({
  items: [],
});
