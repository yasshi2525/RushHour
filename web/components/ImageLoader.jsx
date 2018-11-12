import React from "react";
import PropTypes from "prop-types";
import * as PIXI from "pixi.js";
import { metaGameMap } from "../consts";

export default class ImageLoader extends React.Component {
    constructor(props) {
        super(props);

        metaGameMap.filter(elm => elm.category == "sprite").map(elm => elm.type).forEach(name => {
            props.loader.add(name, `public/img/${name}.png`);
        });

        props.loader.load();
    }

    render() {
        return (null);
    }
}

ImageLoader.propTypes = {
    loader: PropTypes.instanceOf(PIXI.loaders.Loader).isRequired
};
