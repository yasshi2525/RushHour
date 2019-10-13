import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, Theme, TextField, Box, Dialog, DialogTitle, Divider, DialogContent, DialogActions, useTheme, useMediaQuery } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import { RushHourStatus } from "../state";
import * as Actions from "../actions";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        title: {
            marginTop: theme.spacing(1)
        },
        text: {
            display: "block",
            margin: theme.spacing(1)
        },
        button: {
            marginTop: theme.spacing(2)
        },
        error: {
            marginTop: theme.spacing(1),
            fontSize: theme.typography.overline.fontSize,
            color: theme.palette.error.main,
        },
        hidden: {
            display: "none"
        }
    })
);

export default function() {
    const [opened, setOpened] = React.useState(false);
    const theme = useTheme();
    const classes = useStyles(theme);
    const isFullScreen = useMediaQuery(theme.breakpoints.down("xs"));
    const [id, setUserID] = React.useState("");
    const [name, setName] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [confirm, setConfirm] = React.useState("");
    const [idError, setIDError] = React.useState(false);
    const [pwError, setPWError] = React.useState(false);
    const [cnfError, setCnfError] = React.useState(false);
    const dispatch = useDispatch();
    const handleRegister = () => {
        if (id == "") {
            setIDError(true);
        }
        if (password == "") {
            setPWError(true);
        }
        if (id == "" || password == "") {
            return
        }
        if (password != confirm) {
            setCnfError(true);
            return
        }
        dispatch(Actions.register.request({id, name: (name == "" ? "ゲスト" : name), password}));
    };
    
    const failed = useSelector<RushHourStatus, boolean>(state => state.isRegisterFailed);
    const succeeded = useSelector<RushHourStatus, boolean>(state => state.isRegisterSucceeded);
    const formRef = React.useRef<HTMLFormElement>(null);
    
    if (succeeded && formRef.current !== null) {
        formRef.current.submit();
    }

    const handleClose = () => {
        setIDError(false);
        setPWError(false);
        setCnfError(false);
        setOpened(false);
    }

    const handleSubmit = () => {
        if (!succeeded) {
            handleRegister();
        }
    }
    
    return (
        <>
            <Button variant="outlined" color="primary" onClick={() => setOpened(true)}>新規登録</Button>
            <Dialog
                fullScreen={isFullScreen}
                fullWidth={true}
                maxWidth="xs"
                aria-labelledby="modal-title"
                open={opened} 
                onClose={handleClose}>
                <DialogTitle id="modal-title">
                    新規ユーザ登録
                </DialogTitle>
                <Divider />
                <DialogContent>      
                    <form action="/" method="POST" ref={formRef} onSubmit={() => handleSubmit()}>
                        <TextField
                            error={idError}
                            name="id"
                            label="メールアドレス"
                            value={id} onInput={e => setUserID((e.target as HTMLInputElement).value)}
                            className={classes.text}
                            onChange={() => setIDError(false)}
                        />
                        <TextField
                            name="name"
                            label="表示名"
                            placeholder="ゲスト"
                            value={name} onInput={e => setName((e.target as HTMLInputElement).value)}
                            className={classes.text}
                        />
                        <TextField
                            error={pwError}
                            name="password"
                            label="パスワード"
                            type="password"
                            value={password} onInput={e => setPassword((e.target as HTMLInputElement).value)}
                            className={classes.text}
                            onChange={() => {setPWError(false); setCnfError(false)}}
                        />
                        <TextField
                            error={cnfError}
                            name="confirm"
                            label="パスワード(確認)"
                            type="password"
                            value={confirm} onInput={e => setConfirm((e.target as HTMLInputElement).value)}
                            className={classes.text}
                            onChange={() => {setPWError(false); setCnfError(false)}}
                        />
                        { idError && <Box className={classes.error}>メールアドレスを入力してください</Box> }
                        { pwError && <Box className={classes.error}>パスワードを入力してください</Box> }
                        { cnfError && <Box className={classes.error}>パスワード(確認)が一致しません</Box> }
                        { failed && <Box className={classes.error}>入力したメールアドレスはすでに使われています</Box> }
                        <Button
                            variant="contained" 
                            color="primary" 
                            onClick={handleRegister}
                            className={classes.button}>
                            登録
                        </Button>
                        <input className={classes.hidden} type="submit"/>
                    </form>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>戻る</Button>
                </DialogActions>
            </Dialog>
        </>
    );
}