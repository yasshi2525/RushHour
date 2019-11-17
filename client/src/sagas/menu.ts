import { put, call } from "redux-saga/effects";
import * as Action from "actions";
import { MenuStatus } from "state";

const setMenu = (opts: Action.MenuRequest) => {
  if (opts.menu !== MenuStatus.IDLE) {
    opts.model.setMenuState(MenuStatus.IDLE);
  }
  return opts.model.setMenuState(opts.menu);
};

export function* generateSetMenu(
  action: ReturnType<typeof Action.setMenu.request>
) {
  let menu = yield call(setMenu, action.payload);
  return yield put(Action.setMenu.success(menu));
}
