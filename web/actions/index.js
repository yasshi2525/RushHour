export const FETCH_MAP_REQUESTED = "FETCH_MAP_REQUESTED";
export const FETCH_MAP_SUCCEEDED = "FETCH_MAP_SUCCEEDED";
export const FETCH_MAP_FAILED = "FETCH_MAP_FAILED";

export function requestFetchMap(timestamp) {
    return {
        type: FETCH_MAP_REQUESTED,
        payload: {
            timestamp: timestamp
        }
    }
}

