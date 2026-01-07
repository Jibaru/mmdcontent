export namespace entities {
	
	export class Model {
	    id: string;
	    name: string;
	    screenshots: string[];
	    description: string;
	    originalPath: string;
	    embedding?: number[];
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.screenshots = source["screenshots"];
	        this.description = source["description"];
	        this.originalPath = source["originalPath"];
	        this.embedding = source["embedding"];
	    }
	}
	export class Pagination_MMDContent_internal_entities_Model_ {
	    data: Model[];
	    total: number;
	    page: number;
	    perPage: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new Pagination_MMDContent_internal_entities_Model_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], Model);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.perPage = source["perPage"];
	        this.totalPages = source["totalPages"];
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
	export class Stage {
	    id: string;
	    name: string;
	    screenshots: string[];
	    description: string;
	    originalPath: string;
	    embedding?: number[];
	
	    static createFrom(source: any = {}) {
	        return new Stage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.screenshots = source["screenshots"];
	        this.description = source["description"];
	        this.originalPath = source["originalPath"];
	        this.embedding = source["embedding"];
	    }
	}
	export class Pagination_MMDContent_internal_entities_Stage_ {
	    data: Stage[];
	    total: number;
	    page: number;
	    perPage: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new Pagination_MMDContent_internal_entities_Stage_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], Stage);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.perPage = source["perPage"];
	        this.totalPages = source["totalPages"];
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

