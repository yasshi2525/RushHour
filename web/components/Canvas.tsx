import "./style.css";
import * as React from "react";
import { connect } from "react-redux";
import * as PIXI from "pixi.js";
import { config } from "../common/interfaces/gamemap";
import GameModel from "../common/model";
import { MouseDragHandler, TouchDragHandler } from "../common/handlers/drag";
import { WheelHandler } from "../common/handlers/wheel";
import { PinchHandler } from "../common/handlers/pinch";
import { RushHourStatus } from "../state";
import { fetchMap } from "../actions";

// Pixi.js が作成する canvas を管理するコンポーネント
export class Canvas extends React.Component<RushHourStatus, RushHourStatus> {
    app: PIXI.Application;
    model: GameModel;
    ref: React.RefObject<HTMLDivElement>;
    mouse: MouseDragHandler;
    wheel: WheelHandler;
    touch: TouchDragHandler;
    pinch: PinchHandler;

    constructor(props: RushHourStatus) {
        super(props);

        this.app = new PIXI.Application({
            width: window.innerWidth,
            height: window.innerHeight,
            backgroundColor: 0x333333,
            autoStart: true,
            antialias: true,
            resolution: window.devicePixelRatio,
            autoDensity: true
        });

        this.model = new GameModel({
            app: this.app , 
            cx: config.gamePos.default.x, 
            cy: config.gamePos.default.y, 
            scale: config.scale.default,
            zoom: 0
        });

        ["residence", "company", "station", "train"].forEach(key => this.app.loader.add(key, `public/img/${key}.png`));
        this.app.loader.load(() => {
            this.model.attach({
                residence: this.app.loader.resources["residence"].texture,
                company: this.app.loader.resources["company"].texture,
                station: this.app.loader.resources["station"].texture
            });
            this.fetchMap();
        });

        this.ref = React.createRef<HTMLDivElement>();

        window.addEventListener("resize", () => {
            let beforeW = this.model.renderer.width;
            let beforeH = this.model.renderer.height;
            let afterW = window.innerWidth;
            let afterH = window.innerHeight;

            let dW = afterW - beforeW;
            let dH = afterH - beforeH;

            let dx = dW / Math.pow(2, this.model.coord.scale - 1);
            let dy = dH / Math.pow(2, this.model.coord.scale - 1);

            this.model.setCenter(this.model.coord.cx - dx, this.model.coord.cy - dy);
            this.model.renderer.resize(afterW, afterH);
        })

        this.mouse = new MouseDragHandler(this.model, this.props.dispatch);
        this.wheel = new WheelHandler(this.model, this.props.dispatch);
        this.touch = new TouchDragHandler(this.model, this.props.dispatch);
        this.pinch = new PinchHandler(this.model, this.props.dispatch);
    }

    render() {
        return (<div ref={this.ref} 
            onMouseDown={(e) => this.mouse.onStart(e)}
            onMouseMove={(e) => this.mouse.onMove(e)}
            onMouseUp={(e) => this.mouse.onEnd(e)}
            onMouseOut={(e) => this.mouse.onEnd(e)}
            onWheel={(e) => { this.wheel.onStart(e); this.wheel.onMove(e); this.wheel.onEnd(e); }}
            onTouchStart={(e) => {this.touch.onStart(e); this.pinch.onStart(e); }}
            onTouchMove={(e) => {this.touch.onMove(e);  this.pinch.onMove(e)} }
            onTouchEnd={(e) => {this.touch.onEnd(e);  this.pinch.onEnd(e)} }>
            </div>);
    }

    componentDidMount() {
        if (this.ref.current !== null) {
            if (this.app.view instanceof Node) { // 単体テスト時、Node非実装のため
                // 一度描画して、canvas要素を子要素にする
                this.ref.current.appendChild(this.app.view);
            }
        } 
    }

    componentDidUpdate() {
        //let beforeFetch = new Date().getTime();
        this.model.timestamp = this.props.timestamp;
        this.model.mergeAll(this.props.map);
        //let beforeRender = new Date().getTime();
        if (this.model.isChanged()) {
            this.model.render();
        }
        //let afterRender = new Date().getTime();
        //this.model.debugValue = (beforeRender - beforeFetch) + ", " + (afterRender - beforeRender);
    }

    componentWillUnmount() {
        this.model.unmount();
    }

    protected fetchMap() {
        this.props.dispatch(fetchMap.request({
            cx: this.model.coord.cx, 
            cy: this.model.coord.cy, 
            scale: this.model.coord.scale
        }));
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { 
        timestamp: state.timestamp, 
        map: state.map
    };
}

export default connect(mapStateToProps)(Canvas);