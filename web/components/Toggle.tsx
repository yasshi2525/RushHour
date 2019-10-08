import * as React from "react";
import { AnyAction } from "redux";
import { connect } from "react-redux";
import { Fab } from "@material-ui/core";
import { MenuProperty } from "../common/interfaces";
import { setMenu, MenuRequest } from "../actions";
import { MenuStatus, RushHourStatus } from "../state";

interface ToggleState {
    selected: boolean
}

function mapStateToProps(state: RushHourStatus) {
    return { menu: state.menu };
}

function mapDispatchToProps(dispatch: React.Dispatch<AnyAction>) {
    return { setMenu: (opts: MenuRequest) => dispatch(setMenu.request(opts)) };
};

const Toggle = (on = MenuStatus.IDLE, off = MenuStatus.IDLE) => 
    (WrappedComponent: React.ComponentType<MenuProperty>) => {

    class ToggleComponent extends React.Component<MenuProperty, ToggleState> { 
        
        constructor(props: MenuProperty) {
            super(props);
            this.state = { selected: props.menu === on };
            this.toggle = this.toggle.bind(this);
        }

        render() {
            return (
                <>
                    <Fab color="primary" hidden={!this.state.selected} onClick={this.toggle}>
                        <WrappedComponent {...this.props} />
                    </Fab>
                    <Fab hidden={this.state.selected} onClick={this.toggle}>
                        <WrappedComponent {...this.props} />
                    </Fab>
                </>
            );
        }

        protected toggle() {
            let newState = !this.state.selected;
            if (newState) {
                this.props.setMenu({ model: this.props.model, menu: on });
            } else {
                this.props.setMenu({ model: this.props.model, menu: off });
            }
            this.setState({ selected: newState })
        }

        componentDidUpdate() {
            if (this.state.selected && this.props.menu !== on) {
                this.setState({ selected: false })
            }
            if (!this.state.selected && this.props.menu === on) {
                this.setState({ selected: true })
            }
        }
    }

    return connect(mapStateToProps, mapDispatchToProps)(ToggleComponent);
}

export default Toggle;
