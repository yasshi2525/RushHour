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
    SEEK_DEPARTURE
};

export interface GameMap {
    [key: string]: Locatable[],
    "companies": Locatable[],
    "gates": Locatable[],
    "humans": Locatable[],
    "line_tasks": Locatable[],
    "platforms": Locatable[],
    "players": Locatable[],
    "rail_edges": Locatable[],
    "rail_lines": Locatable[],
    "rail_nodes": Locatable[],
    "residences": Locatable[],
    "stations": Locatable[],
    "trains": Locatable[],
};

export interface RushHourStatus {
    [key: string]: any,
    readOnly: boolean,
    timestamp: number,
    map: GameMap,
    menu: MenuStatus,
    needsFetch: boolean
};

export const defaultState: RushHourStatus = {
    readOnly: true,
    timestamp: 0,
    map: {
        "companies": [],
        "gates": [],
        "humans": [],
        "line_tasks": [],
        "platforms": [],
        "players": [],
        "rail_edges": [],
        "rail_lines": [],
        "rail_nodes": [],
        "residences": [],
        "stations": [],
        "trains": [],
    },
    menu: MenuStatus.IDLE,
    needsFetch: true
};