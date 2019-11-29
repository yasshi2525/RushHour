import { Entity } from "interfaces";
import { UserInfo } from "interfaces/user";
import { Point } from "interfaces/gamemap";

export interface Coordinates {
  cx: number;
  cy: number;
  scale: number;
  latest?: number;
}

export interface AnchorStatus {
  pos: Point;
  type: string;
  cid: number;
}

export interface AccountSettings {
  oauth_name: string;
  custom_name: string;
  use_cname: boolean;
  oauth_image: string;
  custom_image: string;
  use_cimage: boolean;
  auth_type: string;
}

export interface DefaultProp {
  my: UserInfo | undefined;
}

export function defaultState(opts: DefaultProp): RushHourStatus {
  return {
    timestamp: 0,
    menu: MenuStatus.IDLE,
    isLoginSucceeded: false,
    isLoginFailed: false,
    isRegisterSucceeded: false,
    isRegisterFailed: false,
    isFetchRequired: false,
    isPlayerFetched: false,
    my: opts.my,
    readOnly: opts.my === undefined,
    settings: undefined,
    waitingFor: undefined,
    inOperation: { waiting: false, value: true },
    isAdmin: opts.my !== undefined ? opts.my.admin : false,
    inPurge: { waiting: false, value: false },
    players: { waiting: false, value: [] }
  };
}
