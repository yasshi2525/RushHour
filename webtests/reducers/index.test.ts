import reducer from "@/reducers";
import { defaultState } from "@/state";
import { fetchMap } from "@/actions";

test("fetches map", () => {
    let time = new Date().getTime();
    const actual = reducer(defaultState({ my: undefined, isAdmin: false, inOperation: true }), {
        type: fetchMap.success.toString(),
        payload: {status: true, timestamp: time, results: {foo: "bar"}}
    });
    expect(actual.timestamp).toEqual(time);
});
