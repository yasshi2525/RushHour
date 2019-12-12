import React, { Suspense, useEffect, useMemo, useState } from "react";
import Box from "@material-ui/core/Box";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import useTheme from "@material-ui/core/styles/useTheme";
import makeStyles from "@material-ui/styles/makeStyles";
import createStyles from "@material-ui/styles/createStyles";
import { Theme } from "@material-ui/core/styles/createMuiTheme";
import LoadingCircle from "common/utils/loading";
import { AuthProvider } from "common/auth";

const GameBoard = React.lazy(() => import("./GameBoard"));
const AppBar = React.lazy(() => import("./AppBar"));

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    error: {
      color: theme.palette.error.main
    }
  })
);

const LogOut = () => {
  const theme = useTheme();
  const classes = useStyles(theme);
  const [press, setPress] = useState(false);
  return press ? (
    <Box>
      <Button
        variant="contained"
        className={classes.error}
        onClick={() => setPress(true)}
      >
        ログアウト
      </Button>
      <Box>
        <Typography variant="caption">認証エラーが発生しました</Typography>
        <Typography variant="caption">
          ログアウトボタンを押してください
        </Typography>
      </Box>
    </Box>
  ) : (
    <Box>
      <Typography variant="caption">画面を更新してください</Typography>
    </Box>
  );
};

export default () => {
  useEffect(() => {
    console.info("after Application");
  }, []);
  return useMemo(
    () => (
      <AuthProvider onError={LogOut}>
        <Suspense fallback={<LoadingCircle />}>
          <AppBar />
          <Suspense fallback={<LoadingCircle />}>
            <GameBoard />
          </Suspense>
        </Suspense>
      </AuthProvider>
    ),
    []
  );
};
