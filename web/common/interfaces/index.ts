import GameContainer from "..";
import GameModel from "../models";

export interface GameBoardProperty {
    [key: string]: any,
    readOnly: boolean,
    game: GameContainer,
    isPIXILoaded: boolean,
    isPlayersFetched: boolean
}

export interface GameComponentProperty {
    [key: string]: any,
    readOnly: boolean,
    model: GameModel
}