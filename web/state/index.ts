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

export interface RushHourStatus {
    [key: string]: any,
    timestamp: number,
    menu: MenuStatus,
    isFetchRequired: boolean,
    isPlayerFetched: boolean
};


export const defaultState: RushHourStatus = {
    readOnly: true,
    timestamp: 0,
    menu: MenuStatus.IDLE,
    isFetchRequired: false,
    isPlayerFetched: false
};