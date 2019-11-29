import React, { useContext } from "react";
import { Theme } from "@material-ui/core/styles/createMuiTheme";
import useTheme from "@material-ui/core/styles/useTheme";
import makeStyles from "@material-ui/styles/makeStyles";
import createStyles from "@material-ui/styles/createStyles";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import Box from "@material-ui/core/Box";
import Container from "@material-ui/core/Container";
import LinearProgress from "@material-ui/core/LinearProgress";
import { ComponentProperty } from "interfaces/component";
import LoadingContext, { LoadingStatus } from "common/loading";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      position: "absolute",
      top: "100px",
      zIndex: 1000
    },
    area: {
      padding: theme.spacing(),
      margin: theme.spacing()
    },
    bar: {
      margin: theme.spacing()
    },
    desc: {
      textAlign: "center"
    }
  })
);

interface LinearProperty extends ComponentProperty {
  phase: LoadingStatus;
}

const Linear = (props: LinearProperty) => {
  const theme = useTheme();
  const classes = useStyles(theme);
  return (
    <Container maxWidth="xs" className={classes.root}>
      <Paper className={classes.area}>
        <LinearProgress
          className={classes.bar}
          aria-describedby="loading-description"
          aria-busy={true}
          variant="determinate"
          value={LoadingStatus.progress(props.phase)}
        />
        <Box className={classes.desc}>
          <Typography id="loading-description" variant="subtitle1">
            {LoadingStatus.description(props.phase)}
          </Typography>
        </Box>
      </Paper>
    </Container>
  );
};

export default () => {
  const { status } = useContext(LoadingContext);
  return status !== LoadingStatus.END ? <Linear phase={status} /> : null;
};
