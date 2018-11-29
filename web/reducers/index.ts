import { ActionType } from "../actions";
import { RushHourStatus } from "../state";

export default (state: RushHourStatus, action: {type: string, payload: any}) => {

    switch (action.type) {
        case ActionType.FETCH_MAP_SUCCEEDED:
            return Object.assign({}, state, {map: action.payload});

        case ActionType.MOVE_SPRITE: {
            const newState = {map: Object.assign({}, state.map)};
            let sprite = newState.map[action.payload.key][action.payload.id];
            sprite.x = action.payload.x;
            sprite.y = action.payload.y;

            return newState;
        }
        case ActionType.DESTROY_SPRITE: {
            const newState = {map: Object.assign({}, state.map)};
            delete newState.map[action.payload.key][action.payload.id];
            return newState;
        }
        case ActionType.FETCH_MAP_REQUESTED:
        default:
            return state;
    }
};
