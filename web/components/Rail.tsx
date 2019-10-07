import * as React from "react";
import { connect } from "react-redux";
import { Fab, Avatar } from "@material-ui/core";
import GameModel from "../common/models";
import { MenuStatus } from "@/state";

interface RailProperty {
    model: GameModel
}

interface RailState {
    selected: boolean
}

export class Rail extends React.Component<RailProperty, RailState> {
    constructor(props: any) {
        super(props);
        this.state = { selected: false };
        this.toggleSelection = this.toggleSelection.bind(this);
    }

    render() {
        return (
            <>
                <Fab color="primary" hidden={!this.state.selected} onClick={this.toggleSelection}>
                    <Avatar alt="rail" src="/public/img/rail.png" />
                </Fab>
                <Fab hidden={this.state.selected} onClick={this.toggleSelection}>
                    <Avatar alt="rail" src="/public/img/rail.png" />
                </Fab>
            </>
        );
    }

    protected toggleSelection() {
        let newState = !this.state.selected;
        if (newState) {
            this.props.model.setMenuState(MenuStatus.SEEK_DEPARTURE);
        } else {
            this.props.model.setMenuState(MenuStatus.IDLE);
        }
        this.setState({ selected: newState })
    }

    componentWillUnmount() {
        this.props.model.setMenuState(MenuStatus.IDLE);
    }
}

function mapStateToProps(_: any) {
    return {};
}

export default connect(mapStateToProps)(Rail);