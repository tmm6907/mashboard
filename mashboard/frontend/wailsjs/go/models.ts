export namespace main {
	
	export class CreateFeedRequest {
	    link: string;
	    title: string;
	    description: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateFeedRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.link = source["link"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.language = source["language"];
	    }
	}
	export class FollowRequest {
	    link: string;
	    title: string;
	    desc: string;
	    collection: string;
	
	    static createFrom(source: any = {}) {
	        return new FollowRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.link = source["link"];
	        this.title = source["title"];
	        this.desc = source["desc"];
	        this.collection = source["collection"];
	    }
	}
	export class GetFeedItemRequest {
	    id: number;
	    saved: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GetFeedItemRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.saved = source["saved"];
	    }
	}
	export class GetFeedItemsRequest {
	    category: string;
	    offset: number;
	    saved: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GetFeedItemsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.offset = source["offset"];
	        this.saved = source["saved"];
	    }
	}
	export class GetFeedsRequest {
	    offset: number;
	    filter: string;
	
	    static createFrom(source: any = {}) {
	        return new GetFeedsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.offset = source["offset"];
	        this.filter = source["filter"];
	    }
	}
	export class HandleSaveFeedItemRequest {
	    id: number;
	    value: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HandleSaveFeedItemRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.value = source["value"];
	    }
	}
	export class Response {
	    error: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.error = source["error"];
	        this.data = source["data"];
	    }
	}

}

