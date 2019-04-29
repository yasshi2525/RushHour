import "./style.css";
import * as React from "react";
import { connect } from "react-redux";
import * as PIXI from "pixi.js";
import GameModel from "../common/model";
import { RushHourStatus } from "../state";
import { MouseDragHandler, TouchDragHandler } from "../common/handlers/drag";
import { PinchHandler } from "../common/handlers/pinch";

const imageResources = ["residence", "company", "station", "train"];

// Pixi.js が作成する canvas を管理するコンポーネント
export class Canvas extends React.Component<RushHourStatus, RushHourStatus> {
    app: PIXI.Application;
    model: GameModel;
    ref: React.RefObject<HTMLDivElement>;
    mouse: MouseDragHandler;
    touch: TouchDragHandler;
    pinch: PinchHandler;

    constructor(props: RushHourStatus) {
        super(props);

        this.app = new PIXI.Application({
            width: window.innerWidth,
            height: window.innerHeight,
            backgroundColor: 0x333333,
            autoStart: true,
            resolution: window.devicePixelRatio
        });

        imageResources.forEach(key => this.app.loader.add(key, `public/img/${key}.png`));
        this.app.loader.load();
        this.model = new GameModel({ app: this.app , cx: 0, cy: 0, scale: 10});
        this.ref = React.createRef<HTMLDivElement>();
        this.mouse = new MouseDragHandler(this.model);
        this.touch = new TouchDragHandler(this.model);
        this.pinch = new PinchHandler(this.model);
    }

    render() {
        return (<div ref={this.ref} 
            onMouseDown={(e) => this.mouse.onStart(e)}
            onMouseMove={(e) => this.mouse.onMove(e)}
            onMouseUp={(e) => this.mouse.onEnd(e)}
            onMouseOut={(e) => this.mouse.onEnd(e)}
            onTouchStart={(e) => { this.touch.onStart(e); this.pinch.onStart(e); }}
            onTouchMove={(e) => { this.touch.onMove(e);  this.pinch.onMove(e)} }
            onTouchEnd={(e) => { this.touch.onEnd(e);  this.pinch.onEnd(e)} }>
            </div>);
    }

    componentDidMount() {
        // 一度描画して、canvas要素を子要素にする
        if (this.ref.current !== null) {
            this.ref.current.appendChild(this.app.view);
        } 
    }

    componentDidUpdate() {
        this.model.mergeAll(this.props.map);
        if (this.model.isChanged()) {
            this.model.render();
        }
    }

    componentWillUnmount() {
        this.model.unmount();
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { map: state.map };
}

export default connect(mapStateToProps)(Canvas);