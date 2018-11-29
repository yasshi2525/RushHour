export const isChangedXY = (my, oth) => {
    return isChangedObject(my, oth, ["x", "y"]);
};

/**
 * ディープ比較する。
 * 比較対象のプロパティが undefined のときは比較をスキップする
 * @param {Object} my 比較元
 * @param {Object} oth 比較対象
 * @param {array} attrs 検査対象プロパティ 
 */
export const deepEquals = (my, oth, attrs) => !attrs.find(key => oth[key] !== undefined && my[key] !== oth[key]);

export const positionGameToView = (x, y, centerX, centerY, scale) => {
    return { x, y };
};

export const positionViewToGame = (x, y) => {
    return { x, y };
};