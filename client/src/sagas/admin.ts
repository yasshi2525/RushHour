import * as Action from "actions";
import { generateRequest, http, Method } from ".";

const statusURL = "api/v1/game";
const startURL = "api/v1/game/start";
const stopURL = "api/v1/game/stop";
const purgeURL = "api/v1/game/purge";

async function status() {
  return await http(statusURL);
}

async function inOperation(value: boolean) {
  let url = value ? startURL : stopURL;
  return await http(url, Method.POST);
}

async function purgeUserData() {
  return await http(purgeURL, Method.DELETE);
}

export function* generateStatus(
  action: ReturnType<typeof Action.gameStatus.request>
) {
  return yield generateRequest(() => status(), action, Action.gameStatus);
}

export function* generateInOperation(
  action: ReturnType<typeof Action.inOperation.request>
) {
  return yield generateRequest(
    () => inOperation(action.payload.value),
    action,
    Action.inOperation
  );
}

export function* generatePurgeUserData(
  action: ReturnType<typeof Action.purgeUserData.request>
) {
  return yield generateRequest(
    () => purgeUserData(),
    action,
    Action.purgeUserData
  );
}
