import { put, call } from "redux-saga/effects";
import * as Action from "../actions";
import GameContainer from "../common";

export const loadImages = (game: GameContainer) => 
    new Promise<GameContainer>((resolve, reject) => {
        game.app.loader.load(() => {
            game.init();
            return resolve(game)
        });
        game.app.loader.onError = () => reject(game);
    });

export function* generatePIXI(action: ReturnType<typeof Action.initPIXI.request>) {
    try {
        let container = yield call(loadImages, action.payload);
        return yield put(Action.initPIXI.success(container));
    } catch (e) {
        return yield put(Action.initPIXI.failure(e))
    }
}
