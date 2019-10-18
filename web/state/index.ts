import { Entry, UserInfo } from "../common/interfaces";
import { Point } from "../common/interfaces/gamemap";

export interface Coordinates {
    cx: number,
    cy: number,
    scale: number,
    latest?: number
};

export interface Identifiable {
    id: string
};

export interface Locatable extends Identifiable {
    x: number,
    y: number
};

export enum MenuStatus {
    IDLE,
    SEEK_DEPARTURE,
    EXTEND_RAIL,
    DESTROY
};

export interface GameMap {
    [key: string]: Locatable[],
    "companies": Locatable[],
    "gates": Locatable[],
    "humans": Locatable[],
    "line_tasks": Locatable[],
    "platforms": Locatable[],
    "rail_edges": Locatable[],
    "rail_lines": Locatable[],
    "rail_nodes": Locatable[],
    "residences": Locatable[],
    "stations": Locatable[],
    "trains": Locatable[],
};

export interface AnchorStatus {
    pos: Point, 
    type: string, 
    cid: number
}

export interface AccountSettings {
    oauth_name: string,
    custom_name: string,
    use_cname: boolean,
    oauth_image: string,
    custom_image: string,
    use_cimage: boolean,
    auth_type: string
}

export interface RushHourStatus {
    [key: string]: any,
    timestamp: number,
    menu: MenuStatus,
    isLoginSucceeded: boolean,
    isLoginFailed: boolean,
    isRegisterSucceeded: boolean,
    isRegisterFailed: boolean,
    isFetchRequired: boolean,
    isPlayerFetched: boolean,
    readOnly: boolean,
    my: UserInfo | undefined,
    settings: AccountSettings | undefined,
    waitingFor: Entry | undefined,
    maintenance: boolean,
    isAdmin: boolean
};

export interface DefaultProp {
    my: UserInfo | undefined,
    isAdmin: boolean,
    maintenance: boolean
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
        maintenance: opts.maintenance,
        isAdmin: opts.isAdmin
    };
};