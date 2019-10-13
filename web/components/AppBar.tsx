import * as React from "react";
import { connect } from "react-redux";
import { AppBar as AppBarOrg, Toolbar, Typography, Button, Avatar, Link, Hidden, Theme } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import { AppBarProperty } from "../common/interfaces";
import { hueToRgb } from "../common/interfaces/gamemap";
import { RushHourStatus } from "../state";
import SignIn from "./SignIn";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        item: {
            marginRight: theme.spacing(1)
        },
        profile: {
            borderStyle: "solid",
            borderWidth: "5px",
            marginRight: theme.spacing(1),
            [theme.breakpoints.up("sm")]: {
                width: 50,
                height: 50,
            }
        },
        name: {
            [theme.breakpoints.down("xs")]: {
                display: "none"
            }
        },
        grow: {
            flexGrow: 1,
        },
    })
);

function AppBar(props: AppBarProperty) {
    const classes = useStyles();

    const myColor = props.hue !== undefined ? `rgb(${hueToRgb(props.hue).join(",")})` : "inherit";

    return (
        <AppBarOrg position="sticky">
            <Toolbar>
                <div className={classes.item}>
                    {/* PC向け */}
                    <Hidden xsDown>
                        <Typography variant="h4">
                            RushHour
                        </Typography>
                    </Hidden>
                    {/* スマホ向け */}
                    <Hidden smUp>
                        <Typography variant="h6">
                            RushHour
                        </Typography>
                    </Hidden>
                </div>
                { props.readOnly ?
                    <>
                        <div className={classes.grow} />
                        <SignIn />
                    </> :
                    <>
                        <Avatar style={{borderColor: myColor}} className={classes.profile} src={props.image} />
                        <div className={classes.name}>{props.displayName}</div>
                        <div className={classes.grow} />
                        <Button className={classes.item} variant="contained">
                            <Link href="/signout">サインアウト</Link>
                        </Button>
                    </> }
            </Toolbar>
        </AppBarOrg>
    );
}

function mapStateToProps(_: RushHourStatus) {
    return {};
}

export default connect(mapStateToProps)(AppBar);
