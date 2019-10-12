import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, Dialog, DialogTitle, List, ListItem, ListItemAvatar, Avatar, ListItemText, DialogActions, Divider, Grid, DialogContent, useTheme, useMediaQuery } from "@material-ui/core";
import { RushHourStatus } from "../state";
import * as Actions from "../actions";
import PasswordLogin from "./Password";

const sns = [
    { image: "twitter", msg: "Twitterでログイン" },
    { image: "google", msg: "Googleでログイン" },
    { image: "github", msg: "GitHubでログイン" },
];

export default function() {
    const [opened, setOpened] = React.useState(false);
    const theme = useTheme();
    const isFullScreen = useMediaQuery(theme.breakpoints.down("sm"));
    const isInline = useMediaQuery(theme.breakpoints.up("sm"));
    const dispatch = useDispatch();
    const failed = useSelector<RushHourStatus, boolean>(state => state.isLoginFailed);
    const handleClose = () => {
        if (failed) {
            dispatch(Actions.resetLoginError());
        }
        setOpened(false);
    }
    return (
        <>
            <Button variant="contained" onClick={() => setOpened(true)}>新規登録/ログイン</Button>
            <Dialog
                fullScreen={isFullScreen}
                fullWidth={true}
                maxWidth="sm"
                aria-labelledby="modal-title"
                open={opened} 
                onClose={handleClose}>
                <DialogTitle id="modal-title">
                    ログイン
                </DialogTitle>
                <Divider />
                <DialogContent>
                    <Grid container >
                        <Grid item xs={12} sm={6}>
                            <List>
                                {sns.map(item => (
                                    <ListItem button key={item.msg} onClick={() => { location.href=`/${item.image}` }}>
                                        <ListItemAvatar>
                                            <Avatar src={`/public/img/${item.image}.png`} />
                                        </ListItemAvatar>
                                        <ListItemText primary={item.msg} />
                                    </ListItem>
                                ))}
                            </List>
                        </Grid>
                        { !isInline && <Grid item xs={12}><Divider/></Grid> }
                        { isInline && <Grid item sm={1}><Divider orientation="vertical" /></Grid> }
                        <Grid item xs={12} sm={5}>
                            <PasswordLogin />
                        </Grid>
                    </Grid>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>戻る</Button>
                </DialogActions>
            </Dialog>
        </>
    )
}