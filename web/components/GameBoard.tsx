import * as React from "react";
import { connect } from "react-redux";
import Canvas from "./Canvas";
import { RushHourStatus } from "../state";

// ゲーム画面のルートコンポーネント
export class GameBoard extends React.Component<RushHourStatus, RushHourStatus> {

    constructor(props: RushHourStatus) {
        super(props);
    }

    render () {
        return (
            <div>
                <Canvas readOnly = {this.props.readOnly} />
            </div>
        );
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { map: state.map };
}

export default connect(mapStateToProps)(GameBoard);