import * as React from "react";
import { connect } from "react-redux";
import { AppBar as AppBarOrg, Toolbar, Typography, Button, Fade, Dialog, Link } from "@material-ui/core";
import * as styles from "./style.css";
import { GameBarProperty } from "../common/interfaces";
import { RushHourStatus } from "../state";


interface GameBarState {
    openModal: boolean
}

class AppBar extends React.Component<GameBarProperty, GameBarState> {
    constructor(props: any) {
        super(props);
        this.state = { openModal: false };
        this.handleOpen = this.handleOpen.bind(this);
        this.handleClose = this.handleClose.bind(this);
    }
    render () {
        return (
            <AppBarOrg position="sticky">
                <Toolbar>
                    <Typography variant="h4">
                        RushHour
                    </Typography>
                    { this.props.readOnly ?
                        <>
                            <Button variant="contained" onClick={this.handleOpen}>登録/サインイン</Button>
                            <Dialog
                                aria-labelledby="modal-title"
                                aria-describedby="modal-description"
                                open={this.state.openModal} 
                                onClose={this.handleClose}>
                                <Fade in={this.state.openModal}>
                                    <div>
                                        <div id="modal-title">
                                            登録/サインイン
                                        </div>
                                        <div id="modal-description">
                                            <Button >
                                                <Link href="/twitter">Twitter</Link>
                                            </Button>
                                            <Button>
                                                <Link href="/google">Google</Link>
                                            </Button>
                                        </div>
                                    </div>
                                </Fade>
                            </Dialog>
                        </> :
                        <>
                            <img className={styles.profile} src={this.props.image} />
                            <div>{this.props.displayName}</div>
                            <Button variant="contained">
                                <Link href="/signout">サインアウト</Link>
                            </Button>
                        </> }
                </Toolbar>
            </AppBarOrg>
        );
    }

    handleOpen() {
        this.setState({ openModal: true });
    }

    handleClose() {
        this.setState({ openModal: false });
    }
}

function mapStateToProps(_: RushHourStatus) {
    return {};
}

export default connect(mapStateToProps)(AppBar);
