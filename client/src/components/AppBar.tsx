import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  AppBar as AppBarOrg,
  Toolbar,
  Typography,
  Button,
  Avatar,
  Link,
  Theme,
  useTheme,
  useMediaQuery
} from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import player from "static/player.png";
import { UserInfo } from "common/interfaces";
import { hueToRgb } from "common/interfaces/gamemap";
import { RushHourStatus } from "state";
import * as Actions from "actions";
import SignIn from "./SignIn";
import UserSettings from "./UserSettings";
import Administrator from "./Administrator";

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
        height: 50
      }
    },
    name: {
      [theme.breakpoints.down("xs")]: {
        display: "none"
      }
    },
    grow: {
      flexGrow: 1
    },
    setting: {
      width: 50,
      height: 50,
      [theme.breakpoints.up("sm")]: {
        borderWidth: "4px",
        width: 50,
        height: 50
      }
    },
    inOperation: {
      color: theme.palette.error.main
    }
  })
);

export default function() {
  const theme = useTheme();
  const classes = useStyles(theme);
  const isTiny = useMediaQuery(theme.breakpoints.down("xs"));
  const readOnly = useSelector<RushHourStatus, boolean>(
    state => state.readOnly
  );
  const my = useSelector<RushHourStatus, UserInfo | undefined>(
    state => state.my,
    (l, r) => {
      if (l === undefined || r === undefined) {
        return l == r;
      } else {
        return Object.keys(l).filter(k => l[k] != r[k]).length == 0;
      }
    }
  );
  const isAdmin = useSelector<RushHourStatus, boolean>(state => state.isAdmin);
  const inOperation = useSelector<RushHourStatus, boolean>(
    state => state.inOperation.value
  );
  const dispatch = useDispatch();
  dispatch(Actions.gameStatus.request({}));

  const myColor =
    my !== undefined ? `rgb(${hueToRgb(my.hue).join(",")})` : "inherit";

  const signOut = () => {
    dispatch(Actions.signout.request({}));
  };

  return (
    <AppBarOrg position="sticky">
      <Toolbar>
        <Typography className={classes.item} variant={isTiny ? "h6" : "h4"}>
          RushHour
        </Typography>
        {!inOperation && (
          <Typography
            className={classes.inOperation}
            variant={isTiny ? "h6" : "h6"}
          >
            メンテナンス中です
          </Typography>
        )}
        {my !== undefined && (
          <>
            <Avatar
              style={{ borderColor: myColor }}
              className={classes.profile}
              src={my.image != "" ? my.image : player}
            />
            <div className={classes.name}>{my.name}</div>
            <div className={classes.grow} />
            {isAdmin && <Administrator />}
            <UserSettings />
            <Button className={classes.item} variant="contained">
              <Link onClick={() => signOut()}>サインアウト</Link>
            </Button>
          </>
        )}
        {readOnly && (
          <>
            <div className={classes.grow} />
            <SignIn />
          </>
        )}
      </Toolbar>
    </AppBarOrg>
  );
}
