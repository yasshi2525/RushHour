import * as PIXI from "pixi.js";
import GameModel from "models";
import Base from "models/base";

const app = new PIXI.Application();
const model = new GameModel({ app, cx: 0, cy: 0, scale: 10, zoom: 0 });
let instance: Base;

beforeEach(() => {
  instance = new Base({ model });
  instance.setupDefaultValues();
  instance.setupUpdateCallback();
});

describe("calls handler", () => {
  const testHandler = jest.fn(() => {});

  test("call update handler when registered key is specified", () => {
    instance.addUpdateCallback("id", testHandler);
    instance.merge("id", 100);
    expect(testHandler).toBeCalled();
  });

  test("call before handler", () => {
    instance.addBeforeCallback(testHandler);
    instance.begin();
    expect(testHandler).toBeCalled();
  });

  test("call after handler", () => {
    instance.addAfterCallback(testHandler);
    instance.end();
    expect(testHandler).toBeCalled();
  });
});

describe("setInitialValues", () => {
  test("set value when unregistered key is specified", () => {
    instance.setInitialValues({ unregistered: 100 });
    expect(instance.get("unregistered")).toBe(100);
  });
  test("update value when registered key is specified", () => {
    instance.setInitialValues({ id: "changed" });
    expect(instance.get("id")).toBe("changed");
  });
});

describe("merge", () => {
  test("set value when unregistered key is specified", () => {
    instance.merge("unregistered", 100);
    expect(instance.get("unregistered")).toBe(100);
    expect(instance.isChanged()).toBe(true);
  });

  test("don't change when same value is specified", () => {
    instance.merge("id", 0);
    expect(instance.isChanged()).toBe(false);
  });

  test("update value when registered key is specified", () => {
    instance.merge("id", "changed");
    expect(instance.get("id")).toBe("changed");
    expect(instance.isChanged()).toBe(true);
  });

  afterEach(() => {
    instance.reset();
  });
});
