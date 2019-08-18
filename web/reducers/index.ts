import * as Actions from "../actions";
import { RushHourStatus, MenuStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: Actions.ActionPayload}) => {
    switch (action.type) {
        case Actions.fetchMap.success.toString():
        case Actions.diffMap.success.toString():
            return Object.assign({}, state, { timestamp: action.payload.timestamp });
        case Actions.depart.success.toString():
            return Object.assign({}, state, { 
                menu: MenuStatus.EXTEND_RAIL
            });
        case Actions.seekDept().type:
            return Object.assign({}, state, { menu: MenuStatus.SEEK_DEPARTURE });
        case Actions.seekDest().type:
            return Object.assign({}, state, { 
                menu: MenuStatus.EXTEND_RAIL
            });
        case Actions.cancelEditting().type:
            return Object.assign({}, state, { menu: MenuStatus.IDLE });
        default:
            return state;
    }
};
