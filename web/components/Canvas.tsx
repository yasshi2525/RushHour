import "./style.css";
import * as React from "react";
import { connect } from "react-redux";
import * as PIXI from "pixi.js";
import GameModel from "../common/model";
import { RushHourStatus } from "../state";

const imageResources = ["residence", "company"];

// Pixi.js が作成する canvas を管理するコンポーネント
class Canvas extends React.Component<RushHourStatus, RushHourStatus> {
    app: PIXI.Application;
    model: GameModel;
    ref: React.RefObject<HTMLDivElement>;

    constructor(props: RushHourStatus) {
        super(props);

        this.app = new PIXI.Application({
            width: window.innerWidth,
            height: window.innerHeight,
            backgroundColor: 0x333333});

        imageResources.forEach(key => this.app.loader.add(key, `public/img/${key}.png`));
        this.app.loader.load();
        this.model = new GameModel({ app: this.app });
        this.ref = React.createRef<HTMLDivElement>();
        this.state = Object.assign({}, props);
    }

    componentWillReceiveProps(props: RushHourStatus) {
        this.setState({map: props.map}, () => {
            this.model.mergeAll(this.state.map);
            if (this.model.isChanged()) {
                this.model.render();
            }
        } );
    }

    render() {
        return (<div ref={this.ref}></div>);
    }

    componentDidMount() {
        // 一度描画して、canvas要素を子要素にする
        if (this.ref.current !== null) {
            this.ref.current.appendChild(this.app.view);
        } 
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { map: state.map };
}

export default connect(mapStateToProps)(Canvas);