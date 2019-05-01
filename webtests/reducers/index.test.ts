import reducer from "@/reducers";
import { defaultState } from "@/state";
import { fetchMap } from "@/actions";

test("fetches map", () => {
    const actual = reducer(defaultState, {
        type: fetchMap.success.toString(),
        payload: {status: true, timestamp: new Date().getTime(), results: {foo: "bar"}}
    });
    expect(actual.map.foo).toEqual("bar");
});
