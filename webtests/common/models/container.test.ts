import Container from "@/common/models/container";
import Model from "@/common/models/base";
import { Monitorable } from "@/common/interfaces/monitor";

let instance: Container<Model>;

class SimpleModel extends Model implements Monitorable {
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues({testprop: "original"});
    }
}

beforeEach(() => {
    instance = new Container(SimpleModel);
    instance.setupDefaultValues();
});

describe("mergeChild", () => {
    let testId = "test";

    test("add child when new instance is specified", () => {
        instance.mergeChild({id: testId});
        expect(instance.existsChild(testId)).toBe(true);
        expect(instance.isChanged()).toBe(true);
    });

    test("update child when registered instance is specified", () => {
        instance.mergeChild({ id: testId, testprop: "initial" });
        instance.mergeChild({ id: testId, testprop: "changed"});
        expect(instance.getChild(testId).get("testprop")).toBe("changed");
        expect(instance.isChanged()).toBe(true);
    });

    test("don't change when child property isn't changed", () => {
        instance.mergeChild({ id: testId, testprop: "hoge" });
        instance.reset();
        instance.mergeChild({ id: testId, testprop: "hoge" });
        expect(instance.isChanged()).toBe(false);
    });

    afterEach(() => {
        instance.removeChild(testId);
        instance.reset();
    });
});

describe("mergeChildren", () => {
    let testId = "test";
    test("remove child when no property is specified", () => {
        instance.mergeChild({ id: testId });
        instance.mergeChildren([], {});
        expect(instance.existsChild(testId)).toBe(false);
    });
});

describe("removeChild", () => {
    test("do nothing when unregisted child is specified", () => {
        instance.removeChild("unregistered");
        expect(instance.isChanged()).toBe(false);
    });
});

