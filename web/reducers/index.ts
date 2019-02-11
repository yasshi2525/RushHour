import { ActionType } from "../actions";
import { RushHourStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: {status: boolean, results: any}}) => {
    if (action.payload !== undefined && !action.payload.status) {
        return state;
    }
    switch (action.type) {
        case ActionType.FETCH_MAP_SUCCEEDED:
            return Object.assign({}, state, {map: action.payload.results});
        case ActionType.FETCH_MAP_REQUESTED:
        default:
            return state;
    }
};
