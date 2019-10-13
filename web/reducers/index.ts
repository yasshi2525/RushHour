import * as Actions from "../actions";
import { ActionPayload } from "..//common/interfaces";
import { RushHourStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: ActionPayload}) => {
    switch (action.type) {
        case Actions.login.success.toString():
            return Object.assign({}, state, { isLoginSucceeded: true });
        case Actions.login.failure.toString():
            return Object.assign({}, state, { isLoginFailed: true });
        case Actions.resetLoginError.toString():
            return Object.assign({}, state, { isLoginFailed: false });
        case Actions.register.success.toString():
            return Object.assign({}, state, { isRegisterSucceeded: true });
        case Actions.register.failure.toString():
            return Object.assign({}, state, { isRegisterFailed: true });
        case Actions.setMenu.success.toString():
            return Object.assign({}, state, { menu: action.payload });
        case Actions.fetchMap.success.toString():
            return Object.assign({}, state, { 
                timestamp: action.payload.timestamp,
                isPlayerFetched: !action.payload.hasUnresolvedOwner,
                isFetchRequired: false
            });
        case Actions.destroy.success.toString(): {
            return Object.assign({}, state, {
                timestamp: action.payload.timestamp,
                isFetchRequired: true
            });
        }
        default:
            return state;
    }
};
