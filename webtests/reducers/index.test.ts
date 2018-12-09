import { ActionType } from "@/actions";
import reducer from "@/reducers";
import { defaultState } from "@/state";

test("fetches map", () => {
    const actual = reducer(defaultState, {
        type: ActionType.FETCH_MAP_SUCCEEDED,
        payload: "test"
    });
    expect(actual.map).toEqual("test");
});
