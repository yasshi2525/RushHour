import { Point } from "../interfaces/gamemap";
import GameModel from "../models";
import BaseHandler from "./base";

/**
 * posのクライアント座標系を中心にscaleを ds 変更したとき、
 * 新たな画面中央のサーバ座標系の座標を取得します。
 * @param model
 * @param to 中心座標(クライアント座標系)
 * @param ds scaleの差分 (負のとき拡大、正のとき縮小)
 */
export function getZoomPos(model: GameModel, from: Point, fromScale: number, to: Point, ds: number) : Point {
    // 画面中心からの相対座標化
    let size = Math.max(
        model.renderer.width, 
        model.renderer.height
    );
    let zoom = Math.pow(2, fromScale);

    let dC = {
        x: to.x - model.renderer.width / 2,
        y: to.y - model.renderer.height / 2
    };

    let dS = {
        x: dC.x / size * zoom,
        y: dC.y / size * zoom,
    };

    // 中心座標の変更
    let fromLen = Math.sqrt(dS.x * dS.x + dS.y * dS.y);

    if (fromLen <= 0) {
        return to;
    }

    let toLen = fromLen * Math.pow(2, ds);

    let theta = Math.atan2(dS.y, dS.x);
    return {
        x: from.x - (toLen - fromLen) * Math.cos(theta),
        y: from.y - (toLen - fromLen) * Math.sin(theta)
    };
}

export abstract class PointHandler <T> extends BaseHandler<T> {
    /**
     * イベントが発生した画面上の座標を取得します。
     */
    protected abstract getClientXY(ev: T): Point;

    /**
     * posのクライアント座標系を中心にscaleを ds 変更したとき、
     * 新たな画面中央のサーバ座標系の座標を取得します。
     * @param pos 中心座標(クライアント座標系)
     * @param ds scaleの差分 (負のとき拡大、正のとき縮小)
     */
    protected zoom(pos: Point, ds: number) : Point {
        return getZoomPos(this.model, this.server.from, this.scale.from, pos, ds);
    }
}