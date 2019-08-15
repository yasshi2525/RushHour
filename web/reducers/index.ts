import * as Actions from "../actions";
import { RushHourStatus, MenuStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: Actions.ActionPayload}) => {
    switch (action.type) {
        case Actions.fetchMap.success.toString():
        case Actions.diffMap.success.toString():
            if (action.payload === undefined || !action.payload.status) {
                return state;
            }
            return Object.assign({}, state, { 
                timestamp: action.payload.timestamp, 
                map: action.payload.results,
                needsFetch: false
            });
        case Actions.depart.success.toString():
            return Object.assign({}, state, { needsFetch: true });
        case Actions.startDeparture().type:
            return Object.assign({}, state, { menu: MenuStatus.SEEK_DEPARTURE });
        case Actions.cancelEditting().type:
            return Object.assign({}, state, { menu: MenuStatus.IDLE });
        default:
            return state;
    }
};
