import React, { FC, useCallback } from "react";
import makeStyles from "@material-ui/styles/makeStyles";
import createStyles from "@material-ui/styles/createStyles";
import Box from "@material-ui/core/Box";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import Close from "@material-ui/icons/Close";
import { VariantType, useSnackbar } from "notistack";
import { Errors, isErrors } from "interfaces/error";

interface Envelop {
  type: VariantType;
  summaries: string[];
  details: string[];
}

const isEnvelop = (obj: any): obj is Envelop =>
  obj instanceof Object &&
  "type" in obj &&
  "summaries" in obj &&
  obj.summaries instanceof Array &&
  "details" in obj &&
  obj.details instanceof Array;

const useStyles = makeStyles(() =>
  createStyles({
    root: { width: "200px" },
    tiny: { lineHeight: "0.1em" }
  })
);

const EnvelopView: FC<{ contents: Envelop }> = props => {
  const classes = useStyles();
  return (
    <Box className={classes.root}>
      {props.contents.summaries.map((msg, i) => (
        <Typography key={i} variant="body2">
          {msg}
        </Typography>
      ))}

      {props.contents.details.map((msg, i) => (
        <Box key={i} lineHeight={1} fontSize="small">
          {msg}
        </Box>
      ))}
    </Box>
  );
};

const convert = (err: Errors): Envelop => ({
  type: "error",
  summaries: err.summaries,
  details: err.messages
});

const useSnack = () => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  return useCallback(
    (_msg: string | Errors | Envelop) => {
      const msg = isEnvelop(_msg)
        ? _msg
        : isErrors(_msg)
        ? convert(_msg)
        : ({ type: "info", summaries: [_msg], details: [] } as Envelop);
      enqueueSnackbar(<EnvelopView contents={msg} />, {
        variant: msg.type,
        action: key => (
          <IconButton size="small" onClick={() => closeSnackbar(key)}>
            <Close />
          </IconButton>
        )
      });
    },
    [enqueueSnackbar, closeSnackbar]
  );
};

export default useSnack;
