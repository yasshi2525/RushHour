import RushHourTheme from ".";
import * as React from "react";
import { connect } from "react-redux";
import { Fab } from "@material-ui/core";
import { ThemeProvider } from "@material-ui/styles";
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
            <ThemeProvider theme={RushHourTheme}>
                <Fab color="primary" hidden={!this.state.selected} onClick={this.toggleSelection}>
                        Rail
                </Fab>
                <Fab hidden={this.state.selected} onClick={this.toggleSelection}>
                        Rail
                </Fab>
            </ThemeProvider>
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