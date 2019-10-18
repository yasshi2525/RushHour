import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Dialog, DialogTitle, Divider, DialogContent, DialogActions, useTheme, Theme, IconButton, useMediaQuery, Typography, CircularProgress, Box, FormControlLabel, Grid, TextField, Paper, Switch, Fade, Button } from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import SettingsIcon from "@material-ui/icons/Settings";
import EditIcon from '@material-ui/icons/Edit';
import SendIcon from '@material-ui/icons/Send';
import { Entry } from "../common/interfaces";
import * as Actions from "../actions";
import { RushHourStatus, AccountSettings } from "../state";

const use_cname = "use_cname";
const custom_name = "custom_name";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        icon: {
            color: "#fff"
        },
        label: {
            display: "flex",
            height: "100%",
            justifyContent: "center",
            alignItems: "center"
        },
        value: {
            marginLeft: theme.spacing(4),
            padding: theme.spacing(1),
            marginBottom: theme.spacing(1)
        },
        disable: {
            color: theme.palette.text.disabled
        },
        loading: {
            marginLeft: theme.spacing(3)
        },
        iconLoading: {
            position: 'absolute',
            zIndex: 1,
        },
    }));

export default function() {
    const [opened, setOpened] = React.useState(false);
    const theme = useTheme();
    const classes = useStyles();
    const isFullScreen = useMediaQuery(theme.breakpoints.down("sm"));
    const settings = useSelector<RushHourStatus, AccountSettings | undefined>(state => state.settings);
    const waitingFor = useSelector<RushHourStatus, Entry | undefined>(state => state.waitingFor);
    const dispatch = useDispatch();
    const [editingName, setEditingName] = React.useState(false);


    const request = () => {
        dispatch(Actions.settings.request({}));
    }

    const handleClose = () => {
        setOpened(false);
    }

    const useCName = (use: boolean) => {
        dispatch(Actions.editSettings.request({ key: use_cname, value: use }));
    }

    const handleEditName = () => {
        if (!editingName) {
            setEditingName(true);
        } else if (settings !== undefined) {
            if (settings.custom_name == "") {
                settings.custom_name = "ゲスト";
            }
            dispatch(Actions.editSettings.request({ key: custom_name, value: settings.custom_name }));
            setEditingName(false);
        }
    }

    const CustomName = () => {
        if (settings === undefined) {
            return null;
        }
        return (
            <Box display="flex" alignItems="center">
                <TextField
                    className={classes.value}
                    name="name"
                    label="表示名"
                    fullWidth
                    disabled={!editingName}
                    defaultValue={settings.custom_name}
                    onInput={e => settings.custom_name = (e.target as HTMLInputElement).value}
                />
                <Box marginRight="auto">
                    <IconButton 
                        disabled={ waitingFor !== undefined }
                        aria-label="edit" 
                        onClick={() => handleEditName()}>
                        { editingName ? <SendIcon color="primary" />
                            : <>
                                <EditIcon color={ waitingFor === undefined ? "primary" : undefined} />
                                { waitingFor !== undefined && waitingFor.key == custom_name &&
                                <CircularProgress className={classes.iconLoading} /> }
                            </> }
                    </IconButton>
                </Box>
            </Box>
        )
    };

    return (
        <>
            <IconButton aria-label="settings" onClick={() => {request(); setOpened(true)}}> 
                <SettingsIcon className={classes.icon} fontSize="large" />
            </IconButton>
            <Dialog
                fullScreen={isFullScreen}
                fullWidth={true}
                maxWidth="sm"
                aria-labelledby="modal-title"
                open={opened} 
                onClose={handleClose}>
                <DialogTitle id="modal-title">
                    アカウント設定
                </DialogTitle>
                <Divider />
                <DialogContent>
                { settings === undefined ? 
                    <Box display="flex" justifyContent="center" alignItems="center">
                        <CircularProgress />
                    </Box>
                    : <>
                        <Grid container alignItems="stretch" spacing={1}>
                            <Grid item xs={12} sm={3}>
                                <Paper className={classes.label}><Typography>表示名(公開)</Typography></Paper>
                            </Grid>
                            <Grid item xs={10} sm={8}>
                                <Box alignItems="center">
                                    { settings.auth_type != "RushHour" ? 
                                        <>
                                            <FormControlLabel
                                                control={
                                                    <Switch
                                                        color="primary"
                                                        checked={!settings.use_cname} 
                                                        onChange={(e) => useCName(!e.target.checked)} value="checkedName" />
                                                }
                                                disabled={ waitingFor !== undefined }
                                                label={
                                                    <Box display="flex" alignItems="center">
                                                        <Typography>{ settings.auth_type } のアカウント名を使用する</Typography>
                                                        { waitingFor !== undefined && waitingFor.key == use_cname && 
                                                            <Box marginRight="auto">
                                                                <CircularProgress className={classes.loading} size={24} />
                                                            </Box>
                                                        }
                                                    </Box>
                                                }
                                            />
                                            <Paper className={ classes.value }>
                                                <Typography className={ settings.use_cname ? classes.disable : undefined}>{ settings.oauth_name }</Typography>
                                            </Paper> 
                                            <Fade in={ settings.use_cname }><CustomName /></Fade>
                                        </>
                                    : <CustomName />}
                                </Box>
                            </Grid>
                        </Grid>
                    </> }
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => handleClose()}>戻る</Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
