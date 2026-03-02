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
	export class PlayerMatchHistoryEntry {
	    account_id: number;
	    player_kills: number;
	    player_deaths: number;
	    player_assists: number;
	    hero_id: number;
	    match_result: number;
	    match_duration_s: number;
	    start_time: number;
	    net_worth: number;
	    game_mode: number;
	
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
	        this.match_duration_s = source["match_duration_s"];
	        this.start_time = source["start_time"];
	        this.net_worth = source["net_worth"];
	        this.game_mode = source["game_mode"];
	    }
	}

}

