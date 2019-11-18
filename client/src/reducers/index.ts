import * as Actions from "actions";
import { defaultState } from "state";
import { AsyncStatus } from "common/interfaces";
import { RushHourStatus } from "state";

export default (
  state: RushHourStatus | undefined,
  action: { type: string; payload: any }
) => {
  if (state == undefined) {
    return defaultState({ my: undefined });
  }
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
    case Actions.settings.success.toString():
      return Object.assign({}, state, { settings: action.payload });
    case Actions.editSettings.request.toString():
      return Object.assign({}, state, { waitingFor: action.payload });
    case Actions.editSettings.success.toString():
      let val = Object.assign({}, state, { waitingFor: undefined });
      if (action.payload.my) {
        return Object.assign({}, val, { my: action.payload.my });
      } else {
        return val;
      }
    case Actions.editSettings.failure.toString():
      return Object.assign({}, state, { waitingFor: undefined });
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
    case Actions.gameStatus.request.toString():
    case Actions.inOperation.request.toString():
      var inOperation: AsyncStatus = Object.assign({}, state.inOperation, {
        waiting: true
      });
      return Object.assign({}, state, inOperation);
    case Actions.gameStatus.success.toString():
    case Actions.inOperation.success.toString():
      var inOperation: AsyncStatus = Object.assign({}, state.inOperation, {
        waiting: false,
        value: action.payload.status
      });
      return Object.assign({}, state, { inOperation });
    case Actions.gameStatus.failure.toString():
    case Actions.inOperation.failure.toString():
      var inOperation: AsyncStatus = Object.assign({}, state.inOperation, {
        waiting: false
      });
      return Object.assign({}, state, { inOperation });
    case Actions.purgeUserData.request.toString():
      var inPurge: AsyncStatus = Object.assign({}, state.inPurge, {
        waiting: true
      });
      return Object.assign({}, state, { inPurge });
    case Actions.purgeUserData.success.toString():
    case Actions.purgeUserData.failure.toString():
      var inPurge: AsyncStatus = Object.assign({}, state.inPurge, {
        waiting: false,
        value: false
      });
      return Object.assign({}, state, { inPurge });
    case Actions.playersPlain.request.toString():
      var players: AsyncStatus = Object.assign({}, state.players, {
        waiting: true
      });
      return Object.assign({}, state, { players });
    case Actions.playersPlain.success.toString():
    case Actions.playersPlain.failure.toString():
      var players: AsyncStatus = Object.assign({}, state.players, {
        waiting: false,
        value: action.payload
      });
      return Object.assign({}, state, { players });
    default:
      return state;
  }
};
