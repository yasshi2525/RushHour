import * as PIXI from "pixi.js";
import GameModel from "models";
import GameMap from "models/map";

let instance: GameMap;

const app = new PIXI.Application();

const modelOptions = {
  app,
  cx: 0,
  cy: 0,
  scale: 0,
  zoom: 0,
  my: 0
};

beforeEach(() => {
  let game = new GameModel(modelOptions);
  //game.init();
  instance = game.gamemap;
  instance.updateDisplayInfo();
});

describe("get", () => {
  test("get nothing when unregistered key is specified", () => {
    expect(instance.get("unregisted", 1)).toBeUndefined();
  });
});
