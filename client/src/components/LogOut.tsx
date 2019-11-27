import React, { useContext, useCallback } from "react";
import Box from "@material-ui/core/Box";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import Link from "@material-ui/core/Link";
import useTheme from "@material-ui/core/styles/useTheme";
import makeStyles from "@material-ui/styles/makeStyles";
import createStyles from "@material-ui/styles/createStyles";
import { Theme } from "@material-ui/core/styles/createMuiTheme";
import { ComponentProperty } from "interfaces/component";
import LoginContext from "common/auth";
import AdminContext from "common/admin";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    error: {
      color: theme.palette.error.main
    }
  })
);

const InputLogOut = () => {
  const [, , logout] = useContext(LoginContext);

  return <Link onClick={logout}>ログアウト</Link>;
};

interface RequiredLogOutProperty extends ComponentProperty {
  messages: string[];
}

const RequiredLogOut = (props: RequiredLogOutProperty) => {
  const theme = useTheme();
  const classes = useStyles(theme);
  return (
    <Box>
      <Button variant="contained" className={classes.error}>
        <InputLogOut />
      </Button>
      <Box>
        {props.messages.map(msg => (
          <Typography variant="caption">{msg}</Typography>
        ))}
      </Box>
    </Box>
  );
};

export default () => {
  const [auth] = useContext(LoginContext);
  const isAdminPage = useContext(AdminContext);

  if (auth[0]) {
    return (
      <RequiredLogOut
        messages={["エラーが発生しました", "一度ログアウトしてください"]}
      />
    );
  } else if (auth[1]) {
    if (isAdminPage && !auth[1].admin) {
      return (
        <RequiredLogOut
          messages={[
            "一般権限でログインしています",
            "一度ログアウトしてください"
          ]}
        />
      );
    } else {
      return (
        <Button variant="contained">
          <InputLogOut />
        </Button>
      );
    }
  } else {
    return null;
  }
};
