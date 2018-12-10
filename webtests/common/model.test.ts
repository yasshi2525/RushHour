import GameModel from "@/common/model";
import * as PIXI from "pixi.js";

let instance: GameModel;
const app = new PIXI.Application();

const testmap = {
    "residences": [{
        id: "1", x: 100, y: 100
    }],
    "companies": []
};

beforeEach(() => {
    instance = new GameModel({app: app});
});

describe("get", () => {
    test("get nothing when unregistered key is specified", () => {
        expect(instance.get("unregisted", "1")).toBeUndefined();
    });
});

describe("mergeAll", () => {
    test("do nothing when unregistered key is specified", () => {
        instance.mergeAll({
            "residences": [],
            "companies": [],
            "unregistered": [{
                id: "1", x: 100, y: 100
            }]
        });

        expect(instance.get("unregisted", "1")).toBeUndefined();
        expect(instance.isChanged()).toBe(false);
    });

    test("set when registered key is specified", () => {
        instance.mergeAll(testmap);

        expect(instance.get("residences", "1")).toBeDefined();
        expect(instance.isChanged()).toBe(true);
    });
});

describe("render", () => {
    test("reset all after rendered", () => {
        instance.mergeAll(testmap);
        instance.render();

        expect(instance.isChanged()).toBe(false);
    });
});