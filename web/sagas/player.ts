import GameMap from "../common/models/map";
import * as Action from "../actions";
import { generateRequest, http, Method } from ".";

const playersURL = "api/v1/players";
const loginURL = "api/v1/login";

export async function fetchPlayers(map: GameMap) {
    let json = await http(playersURL)
    map.mergeChildren("players", json.results);
    map.resolve();
    return json;
}

async function login(opts: Action.LoginRequest) {
    return await http(loginURL, Method.POST, opts);
}

export function* generatePlayers(action: ReturnType<typeof Action.players.request>) {
    return yield generateRequest(() => fetchPlayers(action.payload.model.gamemap), action, Action.players);
}

export function* generateLogin(action: ReturnType<typeof Action.login.request>) {
    return yield generateRequest(() => login(action.payload), action, Action.login);
}