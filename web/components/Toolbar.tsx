import * as React from "react";
import { connect } from "react-redux";
import { Hidden, Fab, Container } from "@material-ui/core";
import ExpandIcon from "@material-ui/icons/Add"
import MinimizeIcon from "@material-ui/icons/Remove"
import { ToolBarProperty } from "../common/interfaces";
import { RushHourStatus } from "../state";
import Rail from "./Rail";

interface ToolBardState {
    expandsMenu: boolean
}

class ToolBar extends React.Component<ToolBarProperty, ToolBardState> {
    constructor(props: ToolBarProperty) {
        super(props);
        this.state = { expandsMenu: false };
        this.toggleMenu = this.toggleMenu.bind(this);
    }

    render() {
        return ( 
            <>
                {/* PC向け */}
                <Hidden xsDown>
                    <Rail model={this.props.model} />
                </Hidden>
                {/* スマホ向け */}
                <Hidden smUp>
                    {/* メニュー表示なし */}
                    <Fab color="primary" hidden={this.state.expandsMenu} onClick={this.toggleMenu}>
                        <ExpandIcon fontSize="large" />
                    </Fab>

                    {/* メニュー表示あり */}
                    <Fab hidden={!this.state.expandsMenu} onClick={this.toggleMenu}>
                        <MinimizeIcon fontSize="large" />
                    </Fab>
                    {this.state.expandsMenu ?
                        <Container hidden={!this.state.expandsMenu}>
                            <Rail model={this.props.model} />
                        </Container>
                    : null }
                </Hidden>
            </>);
    }

    toggleMenu() {
        this.setState({ expandsMenu: !this.state.expandsMenu });
    }
}

function mapStateToProps(_: RushHourStatus) {
    return {};
}

export default connect(mapStateToProps)(ToolBar);