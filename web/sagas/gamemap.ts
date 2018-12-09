import { aync } from ".";
import { ActionType } from "../actions";

const url = "http://5be50b8595e4340013f89011.mockapi.io/gamemap/1/";

const requestGameMap = () => 
    fetch(url).then(response => {
        if (!response.ok) {
            throw Error(response.statusText);
        }
        return response;
    }).then(response => response.json())
    .catch(error => error);

export const fetchMap = () => aync(requestGameMap, ActionType.FETCH_MAP_SUCCEEDED, ActionType.FETCH_MAP_FAILED, "payload");