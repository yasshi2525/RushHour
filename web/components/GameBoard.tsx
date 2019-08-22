import * as React from "react";
import { connect } from "react-redux";
import { GameBoardProperty } from "../common/interfaces";
import ResizeHandler from "../common/handlers/window"
import { RushHourStatus } from "../state";
import { initPIXI } from "../actions";
import Canvas from "./Canvas";
import ToolBar from "./Toolbar";



// ゲーム画面のルートコンポーネント
export class GameBoard extends React.Component<GameBoardProperty, RushHourStatus> {
    resize: ResizeHandler;
    
    constructor(props: GameBoardProperty) {
        super(props);

        this.resize = new ResizeHandler(props.game.model, this.props.dispatch);
    }
    
    componentDidMount() {
        this.props.dispatch(initPIXI.request(this.props.game));
    }

    render () {
        return (
            <div>
                { this.props.isLoaded ? 
                    <>
                        <Canvas readOnly={this.props.readOnly} model={this.props.game.model} />
                        <ToolBar readOnly={this.props.readOnly} model={this.props.game.model} />
                    </>
                  : "ロード中" }
            </div>
        );
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { readOnly: state.readOnly, isLoaded: state.isLoaded };
}

export default connect(mapStateToProps)(GameBoard);