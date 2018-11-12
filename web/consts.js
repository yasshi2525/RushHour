export const metaGameMap = [
    { type: "residence", group: "residences", category: "sprite" },
    { type: "company", group: "companies", category: "sprite" }
];
export const EMPTY_MAP = metaGameMap.reduce((previous, current) => {
    previous[current.group] = [];
    return previous;
}, {});

export const mapGroupToType = groupName => {
    let match = metaGameMap.find( elm => elm.group == groupName);
    return match === undefined ? undefined : match.type;
};

export const initialState = {
    gamemap: EMPTY_MAP
};