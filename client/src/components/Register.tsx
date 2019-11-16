import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, Theme, TextField, Box, Dialog, DialogTitle, Divider, DialogContent, DialogActions, useTheme, useMediaQuery, Grid, Typography, Slider, Container } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import { hueToRgb } from "../common/interfaces/gamemap";
import { RushHourStatus } from "../state";
import * as Actions from "../actions";

interface Attributes {
    hue: number
}

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        title: {
            marginTop: theme.spacing(1)
        },
        text: {
            display: "block",
            margin: theme.spacing(1)
        },
        error: {
            marginTop: theme.spacing(1),
            fontSize: theme.typography.overline.fontSize,
            color: theme.palette.error.main,
        },
        label: {
            marginTop: theme.spacing(1),
            marginLeft: theme.spacing(1)
        },
        lineColor: {
            marginLeft: theme.spacing(1),
            padding: "0px",
            display: "flex",
            alignItems: "flex-end"
        },
        sliderBox: {
            padding: "0px",
            margin: "0px",
            width: "180px",
        },
        sliderItem: {
            padding: "0px",
            margin: "0px",
            width: "180px"
        },
        colorImage: {
            padding: "0px",
            margin: "0px",
            width: "180px",
            height: "10px"
        },
        sampleBox: {
            marginLeft: theme.spacing(2),
            width: "40px",
            height: "30px",
            display: "flex"
        },
        sample: {
            width: "40px",
            height: "30px",
        },
        hidden: {
            display: "none"
        }
    })
);

export default function(props: Attributes) {
    const [opened, setOpened] = React.useState(false);
    const theme = useTheme();
    const classes = useStyles(theme);
    const isFullScreen = useMediaQuery(theme.breakpoints.down("sm"));
    const isInline = useMediaQuery(theme.breakpoints.up("sm"));
    const [id, setUserID] = React.useState("");
    const [name, setName] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [confirm, setConfirm] = React.useState("");
    const [idError, setIDError] = React.useState(false);
    const [pwError, setPWError] = React.useState(false);
    const [cnfError, setCnfError] = React.useState(false);
    const [hue, setHue] = React.useState(props.hue);
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
        dispatch(Actions.register.request({id, name: (name == "" ? "ゲスト" : name), password, hue}));
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

    const sampleRef = React.useRef<HTMLDivElement>(null);

    const handleChange = (_: React.ChangeEvent<{}>, newValue: number | number[]) => {
        let hue = newValue as number;
        setHue(hue);
        if (sampleRef.current !== null) {
            sampleRef.current.style.backgroundColor = `rgb(${hueToRgb(hue).join(",")})`;
        }
    }
    
    return (
        <>
            <Button variant="outlined" color="primary" onClick={() => setOpened(true)}>新規登録</Button>
            <Dialog
                fullScreen={isFullScreen}
                fullWidth={true}
                maxWidth="sm"
                aria-labelledby="modal-title"
                open={opened} 
                onClose={handleClose}>
                <DialogTitle id="modal-title">
                    新規ユーザ登録
                </DialogTitle>
                <Divider />
                <DialogContent>      
                    <form action="/" method="POST" ref={formRef} onSubmit={() => handleSubmit()}>
                        <Grid container >
                            <Grid item xs={12} sm={5}>
                                <TextField
                                    error={idError}
                                    name="id"
                                    label="メールアドレス"
                                    value={id} onInput={e => setUserID((e.target as HTMLInputElement).value)}
                                    className={classes.text}
                                    onChange={() => setIDError(false)}
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
                            </Grid>
                            { !isInline && <Grid item xs={12}><Divider/></Grid> }
                            { isInline && <Grid item sm={1}><Divider orientation="vertical" /></Grid> }
                            <Grid item xs={12} sm={6}>
                                <TextField
                                        name="name"
                                        label="表示名"
                                        placeholder="ゲスト"
                                        value={name} onInput={e => setName((e.target as HTMLInputElement).value)}
                                        className={classes.text}
                                />
                                <Typography className={classes.label}>路線カラー</Typography>
                                <Container className={classes.lineColor}>
                                    <Box className={classes.sliderBox}>
                                        <Slider
                                            className={classes.sliderItem}
                                            value={hue}
                                            min={0} max={360}
                                            onChange={handleChange} />
                                        <img className={classes.colorImage} src={`/assets/bundle/spritesheet/color@${window.devicePixelRatio}x.png`} />
                                    </Box>
                                    <Box className={classes.sampleBox}>
                                        <div 
                                            ref={sampleRef} 
                                            className={classes.sample}
                                            style={{backgroundColor: `rgb(${hueToRgb(props.hue).join(",")})`}}
                                        />
                                    </Box>
                                </Container>
                            </Grid>
                        </Grid>
                        <input className={classes.hidden} type="submit"/>
                    </form>
                </DialogContent>
                <DialogActions>
                    <Button
                        variant="contained" 
                        color="primary" 
                        onClick={handleRegister}>
                        登録
                    </Button>
                    <Button onClick={handleClose}>戻る</Button>
                </DialogActions>
            </Dialog>
        </>
    );
}