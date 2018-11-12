import { initialState } from "../consts";
import { ActionType } from "../actions";

export default (state = initialState, action) => {
    switch (action.type) {
        case ActionType.FETCH_MAP_SUCCEEDED:
            return Object.assign({}, state, {gamemap: action.gamemap});
        case ActionType.FETCH_MAP_REQUESTED:
        default:
            return state;
    }
};
