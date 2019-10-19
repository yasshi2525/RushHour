import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, useTheme, Dialog, DialogTitle, DialogActions, DialogContent, Divider, Theme, Switch, useMediaQuery, Grid, Typography, Box, CircularProgress } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import { AsyncStatus } from "../common/interfaces";
import * as Actions from "../actions";
import { RushHourStatus } from "../state";

const msg = {
    start: "稼働状態にしてもよろしいですか？",
    stop: "メンテナンス状態にしてもよろしいですか？",
    purge: "すべてのユーザおよびユーザデータを削除してもよろしいですか？"
};

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        item: {
            margin: theme.spacing(2)
        },
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
    const inPurge = useSelector<RushHourStatus, AsyncStatus>(state => state.inPurge);
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
            case msg.purge:
                dispatch(Actions.purgeUserData.request({ key: "inPurge", value : true }));
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
                        <Grid className={classes.item} container alignItems="center">
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
                    <Divider />
                    { inPurge.waiting ?
                        <Box display="flex" justifyContent="center" alignItems="center">
                            <CircularProgress />
                        </Box> :
                        <Button 
                            className={`${classes.bg} ${classes.item}`} 
                            variant="contained" 
                            disabled={ inOperation.value }
                            onClick={() => setConfirmMessage(msg.purge)}
                        >
                            ユーザデータ削除
                        </Button>
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