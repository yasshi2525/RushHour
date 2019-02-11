import Base from "@/common/models/base";

let instance: Base;

beforeEach(() => {
    instance = new Base();
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
    test("don't set value when unregistered key is specified", () => {
        instance.setInitialValues({ unregistered: 100 });
        expect(instance.get("unregistered")).toBeUndefined();
    }); 
    test("update value when registered key is specified", () => {
        instance.setInitialValues({ id: "changed" });
        expect(instance.get("id")).toBe("changed");
    });

});

describe("merge", () => {
    test("don't set value when unregistered key is specified", () => {
        instance.merge("unregistered", 100);
        expect(instance.get("unregistered")).toBeUndefined();
        expect(instance.isChanged()).toBe(false);
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


