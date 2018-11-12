import { metaGameMap, mapGroupToType } from "../consts";

export const filterSprite = groupName => 
    metaGameMap.find( elm => 
        elm.type == mapGroupToType(groupName) && elm.category == "sprite"
    ) !== undefined;