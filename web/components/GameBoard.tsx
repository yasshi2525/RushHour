import * as React from "react";
import { connect } from "react-redux";
import Canvas from "./Canvas";
import { RushHourStatus } from "../state";
import ToolBar from "./Toolbar";

// ゲーム画面のルートコンポーネント
export class GameBoard extends React.Component<RushHourStatus, any> {

    constructor(props: RushHourStatus) {
        super(props);
    }

    render () {
        return (
            <div>
                <Canvas readOnly = {this.props.readOnly} />
                <ToolBar />
            </div>
        );
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { readOnly: state.readOnly };
}

export default connect(mapStateToProps)(GameBoard);