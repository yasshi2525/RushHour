import { Entry } from "../common/interfaces";
import GameMap from "../common/models/map";
import { jwtToUserInfo } from "../state";
import * as Action from "../actions";
import { generateRequest, http, Method } from ".";

const playersURL = "api/v1/players";
const loginURL = "api/v1/login";
const signoutURL = "api/v1/signout";
const registerURL = "api/v1/register";
const settingsURL = "api/v1/settings";

export async function fetchPlayers(map: GameMap | undefined = undefined) {
    let json = await http(playersURL)
    if (map !== undefined) {
        map.mergeChildren("players", json);
        map.resolve();
    }
    return json;
}

async function fetchSettings() {
    return await http(settingsURL)
}

async function editSettings(entry: Entry) {
    let res = await http(`${settingsURL}/${entry.key}`, Method.POST, {value: entry.value});
    localStorage.setItem("jwt", res.jwt)
    let my = jwtToUserInfo(res.jwt)
    return my == undefined ? res : { ...res, my }
}

async function login(opts: Action.LoginRequest) {
    let json = await http(loginURL, Method.POST, opts);
    localStorage.setItem("jwt", json.jwt)
    return json;
}

async function signout(opts: Action.Request) {
    let json: string = ""
    try {
        json = await http(signoutURL, Method.POST, opts);
    } finally {
        localStorage.removeItem("jwt");
        location.href = "/";
    }
    return json
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

export function* generateSignOut(action: ReturnType<typeof Action.signout.request>) {
    return yield generateRequest(() => signout(action.payload), action, Action.signout);
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