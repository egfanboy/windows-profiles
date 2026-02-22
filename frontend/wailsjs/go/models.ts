export namespace main {
	
	export class AudioDevice {
	    id: string;
	    name: string;
	    isDefault: boolean;
	    isEnabled: boolean;
	    deviceType: string;
	    state: string;
	    selected: boolean;
	    nickname: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioDevice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.isDefault = source["isDefault"];
	        this.isEnabled = source["isEnabled"];
	        this.deviceType = source["deviceType"];
	        this.state = source["state"];
	        this.selected = source["selected"];
	        this.nickname = source["nickname"];
	    }
	}
	export class Rect {
	    x: number;
	    y: number;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new Rect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class Monitor {
	    deviceName: string;
	    displayName: string;
	    isPrimary: boolean;
	    isActive: boolean;
	    isEnabled: boolean;
	    monitorId: string;
	    bounds: Rect;
	    nickname: string;
	
	    static createFrom(source: any = {}) {
	        return new Monitor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deviceName = source["deviceName"];
	        this.displayName = source["displayName"];
	        this.isPrimary = source["isPrimary"];
	        this.isActive = source["isActive"];
	        this.isEnabled = source["isEnabled"];
	        this.monitorId = source["monitorId"];
	        this.bounds = this.convertValues(source["bounds"], Rect);
	        this.nickname = source["nickname"];
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
	export class Profile {
	    name: string;
	    monitors: Monitor[];
	    audioDevices: AudioDevice[];
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.monitors = this.convertValues(source["monitors"], Monitor);
	        this.audioDevices = this.convertValues(source["audioDevices"], AudioDevice);
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

