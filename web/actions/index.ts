import { createAction } from "typesafe-actions";

export enum ActionType {
    FETCH_MAP_REQUESTED = "FETCH_MAP_REQUESTED",
    FETCH_MAP_SUCCEEDED = "FETCH_MAP_SUCCEEDED",
    FETCH_MAP_FAILED = "FETCH_MAP_FAILED",
    MOVE_SPRITE = "MOVE_SPRITE",
    DESTROY_SPRITE = "DESTROY_SPRITE",
};

export const requestFetchMap = createAction(ActionType.FETCH_MAP_REQUESTED);
export const moveSprite = createAction(ActionType.MOVE_SPRITE, resolve => {
    return (key: string, id: string, x: number, y: number) => resolve({key, id, x, y});
})
export const destroySprite =createAction(ActionType.DESTROY_SPRITE, resolve => {
    return (key: string, id: string) => resolve({key, id});
})
