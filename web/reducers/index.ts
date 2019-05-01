import { fetchMap, diffMap, ActionPayload } from "../actions";
import { RushHourStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: ActionPayload}) => {
    if (action.payload === undefined || !action.payload.status) {
        return state;
    }
    switch (action.type) {
        case fetchMap.success.toString():
        case diffMap.success.toString():
            return Object.assign({}, state, {timestamp: action.payload.timestamp, map: action.payload.results});
        default:
            return state;
    }
};
