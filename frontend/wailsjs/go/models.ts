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

