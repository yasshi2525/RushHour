import * as React from "react";
import { AnyAction } from "redux";
import { connect } from "react-redux";
import { Hidden, Fab } from "@material-ui/core";
import ExpandIcon from "@material-ui/icons/Add"
import MinimizeIcon from "@material-ui/icons/Remove"
import { MenuProperty } from "../common/interfaces";
import { setMenu, MenuRequest } from "../actions";
import { RushHourStatus, MenuStatus } from "../state";

interface MenuState {
    expands: boolean
}

function mapStateToProps(state: RushHourStatus) {
    return { menu: state.menu };
}

function mapDispatchToProps(dispatch: React.Dispatch<AnyAction>) {
    return { setMenu: (opts: MenuRequest) => dispatch(setMenu.request(opts)) };
};

const Menu = () =>
    (WrappedComponent: React.ComponentType<MenuProperty>) => {
    class ResponsiveMenu extends React.Component<MenuProperty, MenuState> {
        constructor(props: MenuProperty) {
            super(props);
            this.state = { expands: false };
            this.expands = this.expands.bind(this);
        }

        render() {
            return ( 
                <>
                    {/* PC向け */}
                    <Hidden xsDown>
                        <WrappedComponent {...this.props} />
                    </Hidden>
                    {/* スマホ向け */}
                    <Hidden smUp>
                        {/* メニュー表示なし */}
                        <Fab color="primary" hidden={this.state.expands} onClick={this.expands}>
                            <ExpandIcon fontSize="large" />
                        </Fab>
    
                        {/* メニュー表示あり */}
                        <Fab hidden={!this.state.expands} onClick={this.expands}>
                            <MinimizeIcon fontSize="large" />
                        </Fab>
                        {this.state.expands ?
                            <WrappedComponent {...this.props} />
                        : null }
                    </Hidden>
                </>);
        }

        protected expands() {
            let newState = !this.state.expands;
            if (!newState) {
                this.props.setMenu({ model: this.props.model, menu: MenuStatus.IDLE });
            }
            this.setState({ expands: newState });
        }

        componentDidUpdate() {
            if (!this.state.expands && this.props.menu !== MenuStatus.IDLE) {
                this.setState({ expands: true })
            }
        }
    }

    return connect(mapStateToProps, mapDispatchToProps)(ResponsiveMenu);
}

export default Menu;