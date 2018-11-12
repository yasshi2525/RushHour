export const ActionType = {
    FETCH_MAP_REQUESTED: "FETCH_MAP_REQUESTED",
    FETCH_MAP_SUCCEEDED: "FETCH_MAP_SUCCEEDED",
    FETCH_MAP_FAILED: "FETCH_MAP_FAILED"
};

export const requestFetchMap = () => ({ type: ActionType.FETCH_MAP_REQUESTED });
