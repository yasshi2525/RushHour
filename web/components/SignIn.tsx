import * as React from "react";
import { connect } from "react-redux";
import { Button, Fade, Dialog, DialogTitle, List, ListItem, ListItemAvatar, Avatar, ListItemText, DialogActions, Divider } from "@material-ui/core";
import { RushHourStatus } from "../state";

const sns = [
    { image: "twitter", msg: "Twitterでログイン" },
    { image: "google", msg: "Googleでログイン" },
    { image: "github", msg: "GitHubでログイン" },
];

function SignIn() {
    const [opened, setOpened] = React.useState(false);
    return (
        <>
            <Button variant="contained" onClick={() => setOpened(true)}>新規登録/ログイン</Button>
            <Dialog
                aria-labelledby="modal-title"
                open={opened} 
                onClose={() => setOpened(false)}>
                <Fade in={opened}>
                    <>
                        <DialogTitle id="modal-title">
                            新規登録/ログイン
                        </DialogTitle>
                        <Divider />
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
                        <DialogActions>
                            <Button onClick={() => setOpened(false)}>戻る</Button>
                        </DialogActions>
                    </>
                </Fade>
            </Dialog>
        </>
    )
}

function mapStateToProps(_: RushHourStatus) {
    return {};
}

export default connect(mapStateToProps)(SignIn);