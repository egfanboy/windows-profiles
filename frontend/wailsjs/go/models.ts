export namespace main {
	
	export class AudioDevice {
	    id: string;
	    name: string;
	    isDefault: boolean;
	    isEnabled: boolean;
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
	        this.selected = source["selected"];
	        this.nickname = source["nickname"];
	    }
	}
	export class AudioProfile {
	    defaultOutputDeviceId: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultOutputDeviceId = source["defaultOutputDeviceId"];
	    }
	}
	export class Monitor {
	    deviceName: string;
	    displayName: string;
	    isPrimary: boolean;
	    isActive: boolean;
	    isEnabled: boolean;
	    monitorId: string;
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
	        this.nickname = source["nickname"];
	    }
	}
	export class Profile {
	    name: string;
	    audio: AudioProfile;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.audio = this.convertValues(source["audio"], AudioProfile);
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
	export class SaveProfileRequest {
	    name: string;
	    defaultOutputDeviceId: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveProfileRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.defaultOutputDeviceId = source["defaultOutputDeviceId"];
	    }
	}

}

