import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  Button,
  Dialog,
  DialogTitle,
  List,
  ListItem,
  ListItemAvatar,
  Avatar,
  ListItemText,
  DialogActions,
  Divider,
  Grid,
  DialogContent,
  Box,
  useTheme,
  useMediaQuery
} from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import twitter from "static/twitter.png";
import google from "static/google.png";
import github from "static/github.png";
import { RushHourStatus } from "state";
import * as Actions from "actions";
import PasswordLogin from "./Password";
import Register from "./Register";

const sns = [
  { image: twitter, msg: "Twitterでログイン", link: "twitter" },
  { image: google, msg: "Googleでログイン", link: "google" },
  { image: github, msg: "GitHubでログイン", link: "github" }
];

const useStyles = makeStyles(() =>
  createStyles({
    register: {}
  })
);

export default function() {
  const classes = useStyles();
  const [opened, setOpened] = React.useState(false);
  const theme = useTheme();
  const isFullScreen = useMediaQuery(theme.breakpoints.down("sm"));
  const isInline = useMediaQuery(theme.breakpoints.up("sm"));
  const dispatch = useDispatch();
  const failed = useSelector<RushHourStatus, boolean>(
    state => state.isLoginFailed
  );
  const handleClose = () => {
    if (failed) {
      dispatch(Actions.resetLoginError());
    }
    setOpened(false);
  };
  return (
    <>
      <Button variant="contained" onClick={() => setOpened(true)}>
        新規登録/ログイン
      </Button>
      <Dialog
        fullScreen={isFullScreen}
        fullWidth={true}
        maxWidth="sm"
        aria-labelledby="modal-title"
        open={opened}
        onClose={handleClose}
      >
        <DialogTitle id="modal-title">ログイン</DialogTitle>
        <Divider />
        <DialogContent>
          <Grid container>
            <Grid item xs={12} sm={6}>
              <List>
                {sns.map(item => (
                  <ListItem
                    button
                    key={item.msg}
                    onClick={() => {
                      location.href = `/${item.link}`;
                    }}
                  >
                    <ListItemAvatar>
                      <Avatar src={item.image} />
                    </ListItemAvatar>
                    <ListItemText primary={item.msg} />
                  </ListItem>
                ))}
              </List>
            </Grid>
            {!isInline && (
              <Grid item xs={12}>
                <Divider />
              </Grid>
            )}
            {isInline && (
              <Grid item sm={1}>
                <Divider orientation="vertical" />
              </Grid>
            )}
            <Grid item xs={12} sm={5}>
              <PasswordLogin />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Box className={classes.register}>
            <Register hue={Math.round(Math.random() * 360)} />
          </Box>
          <Button onClick={handleClose}>戻る</Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
