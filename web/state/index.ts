export interface Identifiable {
    id: string
}

export interface Locatable extends Identifiable {
    x: number,
    y: number
}

export interface GameMap {
    [key: string]: Locatable[],
    "residences": Locatable[],
    "companies": Locatable[]
}

export interface RushHourStatus {
    [key: string]: any,
    readOnly: boolean,
    map: GameMap
}

export const defaultState: RushHourStatus = {
    readOnly: true,
    map: {
        "residences": [],
        "companies": []
    }
};