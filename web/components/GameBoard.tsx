import * as React from "react";
import { connect } from "react-redux";
import Button from '@material-ui/core/Button';
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
                <Button variant="contained" color="primary">Hello World</Button>
            </div>
        );
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { map: state.map };
}

export default connect(mapStateToProps)(GameBoard);