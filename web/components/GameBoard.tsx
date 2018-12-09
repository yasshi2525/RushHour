import * as React from "react";
import { connect } from "react-redux";
import { requestFetchMap, moveSprite, destroySprite } from "../actions";
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
                <button onClick={() => this.props.dispatch(requestFetchMap())}>fetch</button>
                <button onClick={() => this.props.dispatch(moveSprite("residences", "1", 500, 500))}>move</button>
                <button onClick={() => this.props.dispatch(destroySprite("residences", "1"))}>destroy</button>
                <Canvas readOnly = {this.props.readOnly} />
            </div>
        );
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { map: state.map };
}

export default connect(mapStateToProps)(GameBoard);