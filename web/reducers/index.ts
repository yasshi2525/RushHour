import * as Actions from "../actions";
import { RushHourStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: Actions.ActionPayload}) => {
    switch (action.type) {
        case Actions.initPIXI.success.toString():
            return Object.assign({}, state, { isLoaded: true });
        case Actions.fetchMap.success.toString():
        case Actions.diffMap.success.toString():
            return Object.assign({}, state, { timestamp: action.payload.timestamp });
        case Actions.depart.success.toString():
            return state;
        default:
            return state;
    }
};
