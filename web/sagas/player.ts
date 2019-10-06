import GameMap from "../common/models/map";
import * as Action from "../actions";
import { generateRequest, httpGET } from ".";

const playersURL = "api/v1/players";

export async function fetchPlayers(map: GameMap) {
    let json = await httpGET(playersURL)
    map.mergeChildren("players", json.results);
    map.resolve();
    return json;
}

export function* generatePlayers(action: ReturnType<typeof Action.players.request>) {
    return yield generateRequest(() => fetchPlayers(action.payload.model.gamemap), action, Action.players);
}