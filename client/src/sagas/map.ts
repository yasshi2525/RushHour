import GameModel from "common/models";
import * as Action from "actions";
import { generateRequest, http } from ".";

const mapURL = "api/v1/gamemap";

function buildQuery(model: GameModel): string {
  let params = new URLSearchParams();
  params.set("cx", model.coord.cx.toString());
  params.set("cy", model.coord.cy.toString());
  params.set("scale", (model.coord.scale + 1).toString());
  params.set("delegate", model.delegate.toString());
  return params.toString();
}

export async function fetchMap(model: GameModel) {
  let json = await http(mapURL + "?" + buildQuery(model));
  let error = model.gamemap.mergeAll(json);
  model.timestamp = json.timestamp;
  if (model.gamemap.isChanged()) {
    model.gamemap.updateDisplayInfo();
  }
  return { ...json, ...error };
}

export function* generateMap(
  action: ReturnType<typeof Action.fetchMap.request>
) {
  return yield generateRequest(
    () => fetchMap(action.payload.model),
    action,
    Action.fetchMap
  );
}
