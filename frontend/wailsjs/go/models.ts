export namespace audio {
	
	export class TrackMeta {
	    Title: string;
	    Artist: string;
	    Duration: number;
	    Album: string;
	    Genre: string;
	
	    static createFrom(source: any = {}) {
	        return new TrackMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Artist = source["Artist"];
	        this.Duration = source["Duration"];
	        this.Album = source["Album"];
	        this.Genre = source["Genre"];
	    }
	}
	export class AudioFile {
	    Path: string;
	    Metadata: TrackMeta;
	
	    static createFrom(source: any = {}) {
	        return new AudioFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Metadata = this.convertValues(source["Metadata"], TrackMeta);
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

