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

export namespace lan {
	
	export class Game {
	    id: string;
	    host: string;
	    seats: number;
	    taken: number;
	    locked?: boolean;
	    started?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.host = source["host"];
	        this.seats = source["seats"];
	        this.taken = source["taken"];
	        this.locked = source["locked"];
	        this.started = source["started"];
	    }
	}
	export class Found {
	    address: string;
	    protocol: number;
	    app: string;
	    games: Game[];
	
	    static createFrom(source: any = {}) {
	        return new Found(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.protocol = source["protocol"];
	        this.app = source["app"];
	        this.games = this.convertValues(source["games"], Game);
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
	
	export class CardBack {
	    kind: string;
	    image: string;
	
	    static createFrom(source: any = {}) {
	        return new CardBack(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.image = source["image"];
	    }
	}
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
	export class DiceSkin {
	    id: string;
	    name: string;
	    set: string;
	    image: string;
	
	    static createFrom(source: any = {}) {
	        return new DiceSkin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.set = source["set"];
	        this.image = source["image"];
	    }
	}
	export class Fonts {
	    title: string;
	    body: string;
	    cyrillic: string;
	
	    static createFrom(source: any = {}) {
	        return new Fonts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.body = source["body"];
	        this.cyrillic = source["cyrillic"];
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
	export class Icons {
	    coin: string;
	    sword: string;
	    banner: string;
	    sky: string[];
	
	    static createFrom(source: any = {}) {
	        return new Icons(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.coin = source["coin"];
	        this.sword = source["sword"];
	        this.banner = source["banner"];
	        this.sky = source["sky"];
	    }
	}
	export class MatInfo {
	    id: string;
	    name: string;
	    image: string;
	
	    static createFrom(source: any = {}) {
	        return new MatInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.image = source["image"];
	    }
	}
	export class Sounds {
	    attackHit: string;
	    effectCast: string;
	    effectImpact: string;
	    diceRoll: string;
	
	    static createFrom(source: any = {}) {
	        return new Sounds(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.attackHit = source["attackHit"];
	        this.effectCast = source["effectCast"];
	        this.effectImpact = source["effectImpact"];
	        this.diceRoll = source["diceRoll"];
	    }
	}

}

export namespace protocol {
	
	export class SeatInfo {
	    live: boolean;
	    game?: string;
	    host?: string;
	    started?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SeatInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.live = source["live"];
	        this.game = source["game"];
	        this.host = source["host"];
	        this.started = source["started"];
	    }
	}

}

export namespace settings {
	
	export class HostSetup {
	    seats: number;
	    sets: string[];
	    detailed: boolean;
	    response_timer: number;
	    no_bonus_souls: boolean;
	    no_deck_ratios: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HostSetup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.seats = source["seats"];
	        this.sets = source["sets"];
	        this.detailed = source["detailed"];
	        this.response_timer = source["response_timer"];
	        this.no_bonus_souls = source["no_bonus_souls"];
	        this.no_deck_ratios = source["no_deck_ratios"];
	    }
	}
	export class RecentServer {
	    address: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new RecentServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.name = source["name"];
	    }
	}
	export class Rejoin {
	    address: string;
	    token: string;
	    game: string;
	    via: string;
	
	    static createFrom(source: any = {}) {
	        return new Rejoin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.token = source["token"];
	        this.game = source["game"];
	        this.via = source["via"];
	    }
	}
	export class Settings {
	    nickname: string;
	    recent: RecentServer[];
	    animation: string;
	    fullscreen: boolean;
	    mat: string;
	    dice: string;
	    no_end_turn_ask: boolean;
	    rejoin: Rejoin;
	    host: HostSetup;
	    lan_address: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nickname = source["nickname"];
	        this.recent = this.convertValues(source["recent"], RecentServer);
	        this.animation = source["animation"];
	        this.fullscreen = source["fullscreen"];
	        this.mat = source["mat"];
	        this.dice = source["dice"];
	        this.no_end_turn_ask = source["no_end_turn_ask"];
	        this.rejoin = this.convertValues(source["rejoin"], Rejoin);
	        this.host = this.convertValues(source["host"], HostSetup);
	        this.lan_address = source["lan_address"];
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

