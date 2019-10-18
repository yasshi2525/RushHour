import * as Action from "../actions";
import { generateRequest, http, Method } from ".";

const statusURL = "api/v1/game";
const startURL = "api/v1/game/start";
const stopURL = "api/v1/game/stop";

async function status() {
    return await http(statusURL);
}

async function start() {
    return await http(startURL, Method.POST);
}

async function stop() {
    return await http(stopURL, Method.POST);
}

export function* generateStatus(action: ReturnType<typeof Action.gameStatus.request>) {
    return yield generateRequest(() => status(), action, Action.gameStatus);
}

export function* generateStart(action: ReturnType<typeof Action.startGame.request>) {
    return yield generateRequest(() => start(), action, Action.startGame);
}

export function* generateStop(action: ReturnType<typeof Action.stopGame.request>) {
    return yield generateRequest(() => stop(), action, Action.stopGame);
}