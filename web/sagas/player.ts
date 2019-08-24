import * as Action from "../actions";
import { requestURL, isOK as validateCode } from ".";

const url = "api/v1/players";

const request = (url: string, params: Action.GameMapRequest): Promise<any> => 

    fetch(url)
    .then(validateCode)
    .then(response => response.json())
    .then(response => {
        if (!response.status) {
            throw Error(response.results);
        }
        params.model.gamemap.mergeChildren("players", response.results);
        params.model.gamemap.resolve();
        return response;
    })
    .catch(error => error);

export function* players(action: ReturnType<typeof Action.players.request>) {
    return yield requestURL({ request, url, args: action, callbacks: Action.players });
}