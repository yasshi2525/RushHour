import { aync } from ".";
import { ActionType } from "../actions";

const url = "api/v1/gamemap";

const requestGameMap = () => 
    fetch(url).then(response => {
        if (!response.ok) {
            throw Error(response.statusText);
        }
        return response;
    }).then(response => response.json())
    .catch(error => error);

export const fetchMap = () => aync(requestGameMap, ActionType.FETCH_MAP_SUCCEEDED, ActionType.FETCH_MAP_FAILED, "payload");