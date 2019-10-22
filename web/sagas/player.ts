import { Entry } from "../common/interfaces";
import GameMap from "../common/models/map";
import * as Action from "../actions";
import { generateRequest, http, Method } from ".";

const playersURL = "api/v1/players";
const loginURL = "api/v1/login";
const registerURL = "api/v1/register";
const settingsURL = "api/v1/settings";

export async function fetchPlayers(map: GameMap | undefined = undefined) {
    let json = await http(playersURL)
    if (map !== undefined) {
        map.mergeChildren("players", json.results);
        map.resolve();
    }
    return json;
}

async function fetchSettings() {
    return await http(settingsURL)
}

async function editSettings(entry: Entry) {
    return await http(`${settingsURL}/${entry.key}`, Method.POST, {value: entry.value});
}

async function login(opts: Action.LoginRequest) {
    return await http(loginURL, Method.POST, opts);
}

async function register(opts: Action.RegisterRequest) {
    return await http(registerURL, Method.POST, opts);
}

export function* generatePlayers(action: ReturnType<typeof Action.players.request>) {
    return yield generateRequest(() => fetchPlayers(action.payload.model.gamemap), action, Action.players);
}

export function* generatePlayersPlain(action: ReturnType<typeof Action.playersPlain.request>) {
    return yield generateRequest(() => fetchPlayers(), action, Action.playersPlain);
}

export function* generateLogin(action: ReturnType<typeof Action.login.request>) {
    return yield generateRequest(() => login(action.payload), action, Action.login);
}

export function* generateRegister(action: ReturnType<typeof Action.register.request>) {
    return yield generateRequest(() => register(action.payload), action, Action.register);
}

export function* generateSettings(action: ReturnType<typeof Action.settings.request>) {
    return yield generateRequest(() => fetchSettings(), action, Action.settings);
}

export function* generateEditSettings(action: ReturnType<typeof Action.editSettings.request>) {
    return yield generateRequest(() => editSettings(action.payload), action, Action.editSettings);
}