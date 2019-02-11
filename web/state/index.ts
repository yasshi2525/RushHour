export interface Identifiable {
    id: string
}

export interface Locatable extends Identifiable {
    x: number,
    y: number
}

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
}

export interface RushHourStatus {
    [key: string]: any,
    readOnly: boolean,
    map: GameMap
}

export const defaultState: RushHourStatus = {
    readOnly: true,
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
    }
};