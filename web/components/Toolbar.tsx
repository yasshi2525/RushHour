import * as React from "react";
import { connect } from "react-redux";
import { Hidden, Fab, Container } from "@material-ui/core";
import ExpandIcon from "@material-ui/icons/Add"
import MinimizeIcon from "@material-ui/icons/Remove"
import { GameComponentProperty } from "../common/interfaces";
import { RushHourStatus } from "../state";
import Rail from "./Rail";

interface ToolBardState {
    expandsMenu: boolean
}

export class ToolBar extends React.Component<GameComponentProperty, ToolBardState> {
    constructor(props: GameComponentProperty) {
        super(props);
        this.state = { expandsMenu: false };
        this.toggleMenu = this.toggleMenu.bind(this);
    }

    render() {
        if (this.props.readOnly) {
            return ""
        }
        return ( 
            <Container>
                {/* PC向け */}
                <Hidden xsDown>
                    <Container>
                        <Rail model={this.props.model} />
                    </Container>
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
            </Container>);
    }

    toggleMenu() {
        this.setState({ expandsMenu: !this.state.expandsMenu });
    }
}

function mapStateToProps(_: RushHourStatus) {
    return {};
}

export default connect(mapStateToProps)(ToolBar);