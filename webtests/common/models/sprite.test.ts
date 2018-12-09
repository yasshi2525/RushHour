import SpriteModel from "@/common/models/sprite";
import * as PIXI from "pixi.js"; 

let instance: SpriteModel;
const app = new PIXI.Application();
const testValue = 100;

const updateSpriteModel = (model : SpriteModel, testValue: number) => {
    model.merge("x", testValue);
    model.merge("y", testValue);
    model.merge("alpha", testValue);
    model.merge("scale", testValue);
    model.merge("anchor", {x: testValue, y: testValue});
};

const expectSprite = (sprite: PIXI.Sprite, testValue: number) => {
    expect(sprite.x).toBe(testValue);
    expect(sprite.y).toBe(testValue);
    expect(sprite.alpha).toBe(testValue);
    expect(sprite.scale.x).toBe(testValue);
    expect(sprite.scale.y).toBe(testValue);
    expect(sprite.anchor.x).toBe(testValue);
    expect(sprite.anchor.y).toBe(testValue);
}

beforeEach(() => {
    instance = new SpriteModel({
        name: "test",
        container: app.stage,
        loader: app.loader
    });
    instance.setupUpdateCallback();
    instance.setupAfterCallback();
    instance.setupDefaultValues();
});

test("update sprite properties when payload is changed", () => {
    let sprite: PIXI.Sprite | undefined;
    instance.setupBeforeCallback(); // create sprite instance
    instance.begin();

    updateSpriteModel(instance, testValue);
    sprite = instance.getSprite();

    expect(sprite).toBeDefined();
    if (sprite !== undefined) {
        expectSprite(sprite, testValue);
    }

    instance.end();
    expect(instance.getSprite()).toBeUndefined();
});

test("do nothing when sprite creation is failed", () => {
    instance.begin();

    updateSpriteModel(instance, testValue);
    expect(instance.getSprite()).toBeUndefined();
    
    instance.end();
});