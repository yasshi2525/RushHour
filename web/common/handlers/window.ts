import GameModel from "../models";
import { fetchMap } from "../../actions";

export default class {
    model: GameModel;
    dispatch: any;

    constructor(model: GameModel, dispatch: any) {
        this.model = model;
        this.dispatch = dispatch;
        window.addEventListener("resize", () => this.onResize());
    }

    onResize() {
        let needsFetch = this.model.resize(window.innerWidth, window.innerHeight);
        if (needsFetch) {
            this.dispatch(fetchMap.request({ model: this.model }));
        }
    }
}