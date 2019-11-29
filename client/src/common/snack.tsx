import React, { useCallback } from "react";
import makeStyles from "@material-ui/styles/makeStyles";
import createStyles from "@material-ui/styles/createStyles";
import Box from "@material-ui/core/Box";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import Close from "@material-ui/icons/Close";
import { VariantType, useSnackbar } from "notistack";
import { ComponentProperty } from "interfaces/component";
import { ServerErrors, isServerErrors } from "interfaces/error";

interface Envelop {
  type: VariantType;
  summary: string;
  details: string[];
}

const isEnvelop = (obj: any): obj is Envelop =>
  obj instanceof Object &&
  "type" in obj &&
  "summary" in obj &&
  "details" in obj;

interface EnvelopProperty extends ComponentProperty {
  contents: Envelop;
}

const useStyles = makeStyles(() =>
  createStyles({
    root: { width: "200px" },
    tiny: { lineHeight: "0.1em" }
  })
);

const EnvelopView = (props: EnvelopProperty) => {
  const classes = useStyles();
  return (
    <Box className={classes.root}>
      <Typography variant="body2">{props.contents.summary}</Typography>
      <Box lineHeight={1} fontSize="small">
        {props.contents.details.map(msg => msg)}
      </Box>
    </Box>
  );
};

const convert = (err: ServerErrors): Envelop => ({
  type: "error",
  summary: err.summary,
  details: err.messages
});

export const useSnack = () => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  return useCallback((_msg: string | ServerErrors | Envelop) => {
    const msg = isEnvelop(_msg)
      ? _msg
      : isServerErrors(_msg)
      ? convert(_msg)
      : ({ type: "info", summary: _msg, details: [] } as Envelop);
    enqueueSnackbar(<EnvelopView contents={msg} />, {
      variant: msg.type,
      action: key => (
        <IconButton size="small" onClick={() => closeSnackbar(key)}>
          <Close />
        </IconButton>
      )
    });
  }, []);
};
