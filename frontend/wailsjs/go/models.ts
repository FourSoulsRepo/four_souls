export namespace carddb {
	
	export class Footnote {
	    title: string;
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new Footnote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.text = source["text"];
	    }
	}
	export class Reward {
	    count: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new Reward(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.kind = source["kind"];
	    }
	}
	export class Stat {
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Stat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}

}

export namespace legal {
	
	export class Part {
	    text: string;
	    url?: string;
	
	    static createFrom(source: any = {}) {
	        return new Part(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.url = source["url"];
	    }
	}
	export class Line {
	    parts: Part[];
	
	    static createFrom(source: any = {}) {
	        return new Line(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.parts = this.convertValues(source["parts"], Part);
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

export namespace main {
	
	export class CardInfo {
	    ref: string;
	    set: string;
	    type: string;
	    name: string;
	    effects: string[];
	    footnotes: carddb.Footnote[];
	    stats: carddb.Stat[];
	    rewards: carddb.Reward[];
	    image: string;
	
	    static createFrom(source: any = {}) {
	        return new CardInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ref = source["ref"];
	        this.set = source["set"];
	        this.type = source["type"];
	        this.name = source["name"];
	        this.effects = source["effects"];
	        this.footnotes = this.convertValues(source["footnotes"], carddb.Footnote);
	        this.stats = this.convertValues(source["stats"], carddb.Stat);
	        this.rewards = this.convertValues(source["rewards"], carddb.Reward);
	        this.image = source["image"];
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
	export class HostInfo {
	    running: boolean;
	    port: number;
	    addresses: string[];
	
	    static createFrom(source: any = {}) {
	        return new HostInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.port = source["port"];
	        this.addresses = source["addresses"];
	    }
	}

}

export namespace settings {
	
	export class Settings {
	    nickname: string;
	    recent: string[];
	    animation: string;
	    fullscreen: boolean;
	    theme: string;
	    mat: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nickname = source["nickname"];
	        this.recent = source["recent"];
	        this.animation = source["animation"];
	        this.fullscreen = source["fullscreen"];
	        this.theme = source["theme"];
	        this.mat = source["mat"];
	    }
	}

}

export namespace version {
	
	export class Info {
	    app: string;
	    engine: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.app = source["app"];
	        this.engine = source["engine"];
	    }
	}

}

