import "./style.css";
import React from "react";
import * as PIXI from "pixi.js";

// Pixi.js が作成する canvas を管理するコンポーネント
export default class Canvas extends React.Component {
    constructor(props) {
        super(props);

        var _app = new PIXI.Application({
            width: window.innerWidth,
            height: window.innerHeight,
            backgroundColor: 0x333333});

        _app.view.classList.add("gameContainer");
        
        this.state = {
            app: _app,
            renderer: _app.renderer,
            stage: _app.stage,
            view: _app.view
        };

        this.divRef = React.createRef();
    }

    componentDidMount() {
        this.divRef.current.appendChild(this.state.view);
    }

    render() {
        return <div ref={this.divRef}></div>;
    }
}

