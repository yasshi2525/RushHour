import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, useTheme, Dialog, DialogTitle, DialogActions, DialogContent, Divider, Theme, Switch, useMediaQuery, Grid, Typography } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import * as Actions from "../actions";
import { RushHourStatus } from "../state";

const msg = {
    start: "メンテナンス状態を解除してもよろしいですか？",
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
    const isFullScreen = useMediaQuery(theme.breakpoints.down("md"));
    const maintenance = useSelector<RushHourStatus, boolean>(state => state.maintenance);
    const dispatch = useDispatch();

    const handleClose = () => {
        setOpened(false);
    }

    const confirmGameStatus = (maintenance: boolean) => {
        if (maintenance) {
            setConfirmMessage(msg.stop);
        } else {
            setConfirmMessage(msg.start);
        }
    }

    const handleChange = () => {
        switch(confirmMessage) {
            case msg.start:
                dispatch(Actions.startGame.request({}));
                break;
            case msg.stop:
                dispatch(Actions.stopGame.request({}));
                break;
        }
        setConfirmMessage("");
    }

    return (
        <>
            <Button variant="contained" className={classes.bg} onClick={() => setOpened(true)}>管理</Button>
            <Dialog
                fullScreen={isFullScreen}
                fullWidth={true}
                maxWidth="md"
                aria-labelledby="modal-title"
                open={opened} 
                onClose={handleClose}
            >
                <DialogTitle id="modal-title">
                    管理者機能
                </DialogTitle>
                <Divider />
                <DialogContent>
                    <Grid container alignItems="center">
                        <Grid item>
                            <Typography>稼働</Typography>
                        </Grid>
                        <Grid item>
                            <Switch
                                checked={maintenance}
                                onChange={e => confirmGameStatus(e.target.checked)}
                            />
                        </Grid>
                        <Grid item>
                            <Typography className={classes.fr}>停止</Typography>
                        </Grid>
                    </Grid>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => handleClose()}>戻る</Button>
                </DialogActions>
            </Dialog>
            <Dialog
                open={confirmMessage != ""} 
                onClose={() => setConfirmMessage("")}
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