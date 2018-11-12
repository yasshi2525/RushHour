import React from "react";
import * as PIXI from "pixi.js";
import PropTypes from "prop-types";

export default class Sprite extends React.Component {
    constructor(props) {
        super(props);

        let sprite = new PIXI.Sprite(props.loader.resources[props.name].texture);

        sprite.anchor.set(0.5, 0.5);
        sprite.alpha = props.alpha || 1.0;
        sprite.scale.x = props.scale || 0.5;
        sprite.scale.y = props.scale || 0.5;
        sprite.x = props.x;
        sprite.y = props.y;

        props.stage.addChild(sprite);

        this.sprite = sprite;
    }

    render() {
        return (null);
    }

}

Sprite.propTypes = {
    //stage: PropTypes.objectOf(PIXI.Container).isRequired,
    //loader: PropTypes.objectOf(PIXI.loaders.Loader).isRequired,
    // 画像種別
    name: PropTypes.string.isRequired,
    x: PropTypes.number,
    y: PropTypes.number,
    alpha: PropTypes.number,
    scale: PropTypes.number
};