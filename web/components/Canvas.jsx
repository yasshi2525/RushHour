import "./style.css";
import React from "react";
import { connect } from "react-redux";
import ProtoTypes from "prop-types";
import * as PIXI from "pixi.js";
import { mapGroupToType, initialState } from "../consts";
import { filterSprite } from ".";
import ImageLoader from "./ImageLoader";
import Sprite from "./Sprite";

// Pixi.js が作成する canvas を管理するコンポーネント
class Canvas extends React.Component {
    constructor(props) {
        super(props);

        this.state = initialState;

        let app = new PIXI.Application({
            width: window.innerWidth,
            height: window.innerHeight,
            backgroundColor: 0x333333});

        this.app = app;
        this.loader = app.loader;
        this.renderer = app.renderer;
        this.stage = app.stage;
        this.view = app.view;

        this.divRef = React.createRef();
    }

    componentDidMount() {
        this.setState(this.props.gamemap);
        this.divRef.current.appendChild(this.app.view);
    }

    renderSprite(id, typeName, params) {
        return <Sprite key={id} loader={this.loader} stage={this.stage} name={typeName} {...params} />;
    }
    
    spriteForeach(gamemap) {
        Object.keys(gamemap).filter(filterSprite).map(groupName => {
            gamemap[groupName].map(elm => this.renderSprite(elm.id, mapGroupToType(groupName), elm ));
        });
    }

    componentWillReceiveProps(props) {
        this.setState( {gamemap: props.gamemap} );
    }

    render() {
        let sprites = Object.keys(this.state.gamemap).filter(filterSprite)
                        .map(groupName => ({name: mapGroupToType(groupName), instances: this.state.gamemap[groupName] }));
        let htmls = 
            <div ref={this.divRef}>
                <ImageLoader loader={this.loader} />
                { sprites.map(entry => entry.instances.map(obj => 
                    <Sprite key={obj.id} loader={this.loader} stage={this.stage} name={entry.name} {...obj} /> )) 
                }
            </div>;

        { this.renderer.render(this.stage); }

        return (htmls);
    }
}

Canvas.propTypes = {
    gamemap: ProtoTypes.object.isRequired
};

function mapStateToProps(state) {
    return {gamemap : state.gamemap};
}

export default connect(mapStateToProps)(Canvas);