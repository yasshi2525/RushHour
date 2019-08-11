import { fetchMap, diffMap, cancelEditting, startDeparture, ActionPayload } from "../actions";
import { RushHourStatus, MenuStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: ActionPayload}) => {
    switch (action.type) {
        case fetchMap.success.toString():
        case diffMap.success.toString():
            if (action.payload === undefined || !action.payload.status) {
                        return state;
            }
            return Object.assign({}, state, { timestamp: action.payload.timestamp, map: action.payload.results});
        case startDeparture().type:
            return Object.assign({}, state, { menu: MenuStatus.SEEK_DEPARTURE });
        case cancelEditting().type:
            return Object.assign({}, state, { menu: MenuStatus.IDLE });
        default:
            return state;
    }
};
