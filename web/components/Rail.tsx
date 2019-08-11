import RushHourTheme from ".";
import * as React from "react";
import { connect } from "react-redux";
import { Fab } from "@material-ui/core";
import { ThemeProvider } from "@material-ui/styles";
import { startDeparture, cancelEditting } from "../actions";

interface RailState {
    selected: boolean
}

export class Rail extends React.Component<any, RailState> {
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
            this.props.dispatch(startDeparture());
        } else {
            this.props.dispatch(cancelEditting());
        }
        this.setState({ selected: newState })
    }

    protected startDeparture() {
        
    }
}

function mapStateToProps(_: any) {
    return {};
}

export default connect(mapStateToProps)(Rail);