export namespace deadlock {
	
	export class HeroAssets {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new HeroAssets(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class Player {
	    account_id: number;
	    team: number;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.account_id = source["account_id"];
	        this.team = source["team"];
	    }
	}
	export class MatchInfo {
	    players: Player[];
	    winning_team: number;
	
	    static createFrom(source: any = {}) {
	        return new MatchInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
	        this.winning_team = source["winning_team"];
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
	export class MatchInfoResponse {
	    match_info: MatchInfo;
	
	    static createFrom(source: any = {}) {
	        return new MatchInfoResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.match_info = this.convertValues(source["match_info"], MatchInfo);
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
	
	export class PlayerMatchHistoryEntry {
	    account_id: number;
	    player_kills: number;
	    player_deaths: number;
	    player_assists: number;
	    hero_id: number;
	    match_result: number;
	    player_team: number;
	    match_duration_s: number;
	    start_time: number;
	    net_worth: number;
	    game_mode: number;
	    match_id: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerMatchHistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.account_id = source["account_id"];
	        this.player_kills = source["player_kills"];
	        this.player_deaths = source["player_deaths"];
	        this.player_assists = source["player_assists"];
	        this.hero_id = source["hero_id"];
	        this.match_result = source["match_result"];
	        this.player_team = source["player_team"];
	        this.match_duration_s = source["match_duration_s"];
	        this.start_time = source["start_time"];
	        this.net_worth = source["net_worth"];
	        this.game_mode = source["game_mode"];
	        this.match_id = source["match_id"];
	    }
	}

}

