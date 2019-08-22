import * as PIXI from "pixi.js";
import GameModel from "@/common/models";
import { SpriteModel } from "@/common/models/sprite";

let instance: SpriteModel;
const app = new PIXI.Application();
const model = new GameModel({ app: app, cx: 0, cy: 0, scale: 10, zoom: 0 });
const testValue = 100;

const updateSpriteModel = (model : SpriteModel, testValue: number) => {
    model.merge("pos", {x: testValue, y: testValue});
    model.merge("alpha", testValue);
    model.merge("spscale", testValue);
    model.merge("anchor", {x: testValue, y: testValue});
};

const expectSprite = (sprite: PIXI.Sprite, testValue: number) => {
    expect(sprite.alpha).toBe(testValue);
    expect(sprite.scale.x).toBe(testValue);
    expect(sprite.scale.y).toBe(testValue);
}

class SimpleSpriteModel extends SpriteModel {

}

beforeEach(() => {
    instance = new SimpleSpriteModel({
        model: model,
        texture: PIXI.Texture.EMPTY, 
        container: new PIXI.Container(), 
        app: app, cx: 0, cy: 0, scale: 10});
    instance.setupUpdateCallback();
    instance.setupAfterCallback();
    instance.setupDefaultValues();
});

test("update sprite properties when payload is changed", () => {
    let sprite: PIXI.Sprite;
    instance.setupBeforeCallback(); // create sprite instance
    instance.begin();

    updateSpriteModel(instance, testValue);
    sprite = instance.sprite;

    expect(sprite).toBeDefined();
    if (sprite !== undefined) {
        expectSprite(sprite, testValue);
    }

    instance.end();
});

test("do nothing when sprite creation is failed", () => {
    instance.begin();

    updateSpriteModel(instance, testValue);
    
    instance.end();
});