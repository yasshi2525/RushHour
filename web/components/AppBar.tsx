import * as React from "react";
import { AppBar as AppBarOrg, Toolbar, Typography, Button, Avatar, Link, Theme, useTheme, useMediaQuery } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import { UserInfo } from "../common/interfaces";
import { hueToRgb } from "../common/interfaces/gamemap";
import SignIn from "./SignIn";
import UserSettings from "./UserSettings";
import { RushHourStatus } from "@/state";
import { useSelector } from "react-redux";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        item: {
            marginRight: theme.spacing(1)
        },
        profile: {
            borderStyle: "solid",
            borderWidth: "2px",
            marginRight: theme.spacing(1),
            [theme.breakpoints.up("sm")]: {
                borderWidth: "4px",
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
        setting: {
            width: 50,
                height: 50,
                [theme.breakpoints.up("sm")]: {
                    borderWidth: "4px",
                    width: 50,
                    height: 50,
                }
        }
    })
);

export default function() {
    const theme = useTheme();
    const classes = useStyles(theme);
    const isTiny = useMediaQuery(theme.breakpoints.down("xs"));
    const readOnly = useSelector<RushHourStatus, boolean>(state => state.readOnly);
    const my = useSelector<RushHourStatus, UserInfo | undefined>(state => state.my, (l, r) => {
        if (l === undefined || r === undefined) {
            return l == r
        } else {
            return Object.keys(l).filter(k => l[k] != r[k]).length == 0;
        }
    });
    
    const myColor = my !== undefined ? `rgb(${hueToRgb(my.hue).join(",")})` : "inherit";

    return (
        <AppBarOrg position="sticky">
            <Toolbar>
                <Typography className={classes.item} variant={ isTiny ? "h6" : "h4" }>RushHour</Typography>
                { my !== undefined &&
                    <>
                        <Avatar style={{borderColor: myColor}} className={classes.profile} src={my.image} />
                        <div className={classes.name}>{my.name}</div>
                        <div className={classes.grow} />
                        <UserSettings />
                        <Button className={classes.item} variant="contained">
                            <Link href="/signout">サインアウト</Link>
                        </Button>
                    </> 
                }
                { readOnly &&
                    <>
                        <div className={classes.grow} />
                        <SignIn />
                    </>
                }
            </Toolbar>
        </AppBarOrg>
    );
}
