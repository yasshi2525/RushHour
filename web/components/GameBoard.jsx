import React from "react";
import { connect } from "react-redux";
import ProtoTypes from "prop-types";
import { EMPTY_MAP } from "../consts";
import { requestFetchMap } from "../actions";
import Canvas from "./Canvas";

// ゲーム画面のルートコンポーネント
class GameBoard extends React.Component {

    constructor(props) {
        super(props);
        this.readOnly = props.readOnly;
        this.gamemap = EMPTY_MAP;
    }

    render () {
        return (
            <div>
                <div>{ JSON.stringify(this.props.state.payload) }</div>
                <button onClick={() => this.props.dispatch(requestFetchMap(1))}>fetch</button>
                <Canvas readOnly= {this.readOnly} gamemap={this.gamemap} />
            </div>
        );
    }
}

GameBoard.propTypes = {
    readOnly: ProtoTypes.bool
};

function mapStateToProps(state) {
    return {state};
}

export default connect(mapStateToProps)(GameBoard);