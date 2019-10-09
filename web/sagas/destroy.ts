import * as Action from "../actions";
import { generateRequest, http, Method } from ".";

const destroyURL = "api/v1";

export async function deleteObject(opts: Action.DestroyRequest) {
    let json = await http(`${destroyURL}/${opts.resource}`, Method.DELETE, { id: opts.cid })
    opts.model.gamemap.removeChild(opts.resource, `${opts.id}`);
    return json;
}

export function* generateDestroy(action: ReturnType<typeof Action.destroy.request>) {
    return yield generateRequest(() => deleteObject(action.payload), action, Action.destroy);
}
