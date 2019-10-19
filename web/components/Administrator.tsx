import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, useTheme, Dialog, DialogTitle, DialogActions, DialogContent, Divider, Theme, Switch, useMediaQuery, Grid, Typography, Box, CircularProgress } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import { AsyncStatus } from "../common/interfaces";
import * as Actions from "../actions";
import { RushHourStatus } from "../state";

const msg = {
    start: "稼働状態にしてもよろしいですか？",
    stop: "メンテナンス状態にしてもよろしいですか？"
};

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        bg: {
            color: "#fff",
            background: theme.palette.error.main
        },
        fr: {
            color: theme.palette.error.main
        }
    }));

export default function() {
    const theme = useTheme();
    const classes = useStyles(theme);
    const [opened, setOpened] = React.useState(false);
    const [confirmMessage, setConfirmMessage] = React.useState("");
    const isFullScreen = useMediaQuery(theme.breakpoints.down("sm"));
    const inOperation = useSelector<RushHourStatus, AsyncStatus>(state => state.inOperation);
    const dispatch = useDispatch();

    const request = () => {
        dispatch(Actions.gameStatus.request({}));
    }

    const handleClose = () => {
        setOpened(false);
    }

    const confirmGameStatus = (status: boolean) => {
        if (status) {
            setConfirmMessage(msg.start);
        } else {
            setConfirmMessage(msg.stop);
        }
    }

    const handleChange = () => {
        switch(confirmMessage) {
            case msg.start:
                dispatch(Actions.inOperation.request({ key: "inOperation", value : true }));
                break;
            case msg.stop:
                dispatch(Actions.inOperation.request({ key: "inOperation", value : false }));
                break;
        }
        setConfirmMessage("");
    }

    return (
        <>
            <Button variant="contained" className={classes.bg} onClick={() => {request(); setOpened(true);}}>管理</Button>
            <Dialog
                fullScreen={isFullScreen}
                fullWidth={true}
                maxWidth="sm"
                aria-labelledby="modal-title"
                open={opened} 
                onClose={handleClose}
            >
                <DialogTitle id="modal-title">
                    管理者機能
                </DialogTitle>
                <Divider />
                <DialogContent>
                    { inOperation.waiting ? 
                        <Box display="flex" justifyContent="center" alignItems="center">
                            <CircularProgress />
                        </Box> :
                        <Grid container alignItems="center">
                            <Grid item>
                                <Typography>稼働</Typography>
                            </Grid>
                            <Grid item>
                                <Switch
                                    checked={!inOperation.value}
                                    onChange={e => confirmGameStatus(!e.target.checked)}
                                />
                            </Grid>
                            <Grid item>
                                <Typography className={classes.fr}>停止</Typography>
                            </Grid>
                        </Grid>
                    }
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => handleClose()}>戻る</Button>
                </DialogActions>
            </Dialog>
            <Dialog
                open={confirmMessage != ""} 
                onClose={() => setConfirmMessage("")}
                fullWidth={true}
                maxWidth="xs"
                aria-labelledby="modal-confirmation-title"
            >
                <DialogTitle id="modal-confirmation-title">
                    確認
                </DialogTitle>
                <DialogContent>
                    <Typography className={classes.fr}>{confirmMessage}</Typography>
                </DialogContent>
                <DialogActions>
                    <Button className={classes.bg} variant="contained" onClick={handleChange}>変更</Button>
                    <Button onClick={() => setConfirmMessage("")}>戻る</Button>
                </DialogActions>
            </Dialog>
        </>
    )
}