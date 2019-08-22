import { put, call } from "redux-saga/effects";
import * as Action from "../actions";
import GameContainer from "../common";

const loadImages = (game: GameContainer) => 
    new Promise<any>((resolve, reject) => {
        game.app.loader.load(() => {
                game.initModel();
                return resolve(game)
        });
        game.app.loader.onError = () => reject(game);
    });

export function* initPIXI(action: ReturnType<typeof Action.initPIXI.request>) {
    try {
        let container = yield call(loadImages, action.payload);
        return yield put(Action.initPIXI.success(container));
    } catch (e) {
        return yield put(Action.initPIXI.failure(e))
    }
}
