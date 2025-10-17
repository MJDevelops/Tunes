export namespace events {
	
	export enum Event {
	    QUEUE_STARTED = "tunes:dqueue:started",
	    QUEUE_DONE = "tunes:dqueue:done",
	    DOWNLOAD_STARTED = "tunes:dl:downloadStarted",
	    DOWNLOAD_INTERRUPT = "tunes:dl:downloadInterrupt",
	    DOWNLOAD_FINISHED = "tunes:dl:downloadFinished",
	    DOWNLOAD_PROGRESS = "tunes:dl:downloadProgress",
	    TRACK_PROGRESS = "tunes:track:progress",
	}

}

export namespace ytdlp {
	
	export class Download {
	    ID: string;
	    Url: string;
	    Options: string[];
	
	    static createFrom(source: any = {}) {
	        return new Download(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Url = source["Url"];
	        this.Options = source["Options"];
	    }
	}
	export class Thumbnail {
	    url: string;
	    height: string;
	    width: string;
	    resolution: string;
	
	    static createFrom(source: any = {}) {
	        return new Thumbnail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.height = source["height"];
	        this.width = source["width"];
	        this.resolution = source["resolution"];
	    }
	}

}

