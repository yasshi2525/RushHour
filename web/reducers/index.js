import * as Type from '../actions';

const initialState = {
    payload: {
        timestamp: 0
    }
};

export default (state = initialState, action) => {
    switch (action.type) {
        case Type.FETCH_MAP_REQUESTED:
        case Type.FETCH_MAP_SUCCEEDED:
            return Object.assign({}, state, {payload: action.payload});
        default:
            return state;
    }
}
