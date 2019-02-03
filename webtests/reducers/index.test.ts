import { ActionType } from "@/actions";
import reducer from "@/reducers";
import { defaultState } from "@/state";

test("fetches map", () => {
    const actual = reducer(defaultState, {
        type: ActionType.FETCH_MAP_SUCCEEDED,
        payload: {status: true, results: {foo: "bar"}}
    });
    expect(actual.map).toEqual({foo: "bar"});
});
