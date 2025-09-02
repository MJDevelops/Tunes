export namespace events {
	
	export enum Event {
	    QUEUE_STARTED = "tunes:dqueue:started",
	    QUEUE_DONE = "tunes:dqueue:done",
	    DOWNLOAD_STARTED = "tunes:dqueue:downloadStarted",
	    DOWNLOAD_INTERRUPT = "tunes:dqueue:downloadInterrupt",
	    DOWNLOAD_FINISHED = "tunes:dqueue:downloadFinished",
	}

}

export namespace ytdlp {
	
	export class Download {
	    ID: string;
	    Url: string;
	    Progress: number;
	
	    static createFrom(source: any = {}) {
	        return new Download(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Url = source["Url"];
	        this.Progress = source["Progress"];
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
	export class ThumbnailJson {
	    thumbnails: Thumbnail[];
	
	    static createFrom(source: any = {}) {
	        return new ThumbnailJson(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.thumbnails = this.convertValues(source["thumbnails"], Thumbnail);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

