import * as Actions from "../actions";
import { ActionPayload } from "..//common/interfaces";
import { RushHourStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: ActionPayload}) => {
    switch (action.type) {
        case Actions.fetchMap.success.toString():
            return Object.assign({}, state, { 
                timestamp: action.payload.timestamp,
                isPlayerFetched: !action.payload.hasUnresolvedOwner
            });
        default:
            return state;
    }
};
