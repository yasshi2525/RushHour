import GameContainer from "..";
import GameModel from "../models";

export interface GameBoardProperty {
    [key: string]: any,
    readOnly: boolean,
    displayName: string | undefined,
    image: string | undefined,
    game: GameContainer,
    isPIXILoaded: boolean,
    isPlayersFetched: boolean
}

export interface GameBarProperty {
    [key: string]: any,
    readOnly: boolean,
    displayName: string | undefined,
    image: string | undefined
}

export interface GameComponentProperty {
    [key: string]: any,
    readOnly: boolean,
    model: GameModel
}