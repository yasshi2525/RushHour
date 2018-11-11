import React from "react";
import Canvas from "./Canvas"
import { requestFetchMap } from "../actions"
import { connect } from "react-redux";

// ゲーム画面のルートコンポーネント
class GameBoard extends React.Component {

    constructor(props) {
        super(props);
        this.readOnly = props.readOnly;
    }

    render () {
        return (
            <div>
                <div>{ JSON.stringify(this.props.state.payload) }</div>
                <button onClick={() => this.props.dispatch(requestFetchMap(1))}>fetch</button>
                <Canvas readOnly= {this.readOnly} />
            </div>
        );
    }
}

function mapStateToProps(state) {
    return {state};
}

export default connect(mapStateToProps)(GameBoard);