import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, Theme, TextField, Box } from "@material-ui/core";
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
    const classes = useStyles();
    const [id, setUserID] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [idError, setIDError] = React.useState(false);
    const [pwError, setPWError] = React.useState(false);
    const dispatch = useDispatch();
    const handleLogin = () => {
        if (id == "") {
            setIDError(true);
        }
        if (password == "") {
            setPWError(true);
        }
        if (id == "" || password == "") {
            return
        }
        dispatch(Actions.login.request({id, password}))
    };
    const failed = useSelector<RushHourStatus, boolean>(state => state.isLoginFailed);
    const succeeded = useSelector<RushHourStatus, boolean>(state => state.isLoginSucceeded);
    const formRef = React.useRef<HTMLFormElement>(null);

    if (succeeded && formRef.current !== null) {
        formRef.current.submit();
    }

    const handleSubmit = () => {
        if (!succeeded) {
            handleLogin();
        }
    }
    
    return (
        <form action="/" method="POST" ref={formRef} onSubmit={() => handleSubmit()}>
            <Box className={classes.title}>RushHourのアカウント</Box>
            <TextField
                error={idError}
                name="id"
                label="メールアドレス"
                value={id} onInput={e => setUserID((e.target as HTMLInputElement).value)}
                className={classes.text}
            />
            <TextField
                error={pwError}
                name="password"
                label="パスワード"
                type="password"
                value={password} onInput={e => setPassword((e.target as HTMLInputElement).value)}
                className={classes.text}
            />
            { idError && <Box className={classes.error}>メールアドレスを入力してください</Box> }
            { pwError && <Box className={classes.error}>パスワードを入力してください</Box> }
            { failed && <Box className={classes.error}>メールアドレスまたはパスワードが間違っています</Box> }
            <Button 
                variant="contained" 
                color="primary" 
                onClick={handleLogin}
                className={classes.button}>
                ログイン
            </Button>
            <input className={classes.hidden} type="submit"/>
        </form>
    );
}