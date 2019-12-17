import React, { useContext, useMemo, useCallback, useRef } from "react";
import Box from "@material-ui/core/Box";
import {
  Paper,
  Container,
  Theme,
  useTheme,
  Button,
  IconButton
} from "@material-ui/core";
import { makeStyles, createStyles } from "@material-ui/styles";
import ZoomInIcon from "@material-ui/icons/ZoomIn";
import ZoomOutIcon from "@material-ui/icons/ZoomOut";
import ConfigContext from "common/config";
import CoordContext from "common/coord";
import DelegateContext from "common/delegate";
import useCoreMapStorage, { hash } from "common/utils/map_storage";
import useCoreMap from "common/utils/map_core";
import useServerMap from "common/utils/map_server";
import useCoreMapChunk, { Chunk } from "common/utils/map_chunk";
import { FlashStatus, GraceHandler } from "common/utils/flash";

const X = 0;
const Y = 1;
const S = 2;
const D = 3;

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    table: {
      width: "100%",
      textAlign: "center",
      margin: "auto"
    },
    row: {
      justifyContent: "space-around"
    },
    cell: {
      margin: theme.spacing(1)
    }
  })
);

const CacheMap = () => {
  const theme = useTheme();
  const classes = useStyles(theme);
  const [{ min_scale, max_scale }] = useContext(ConfigContext);
  const [coordX, coordY, coordS, , , update] = useContext(CoordContext);
  const delegate = useContext(DelegateContext);

  const [core, add, sub] = useCoreMap();

  const [current, cube] = useCoreMapChunk(
    coordX,
    coordY,
    coordS,
    delegate,
    min_scale,
    max_scale
  );

  const graceHandler = useRef<GraceHandler>({
    prepared: false,
    send: () => console.warn("not initialized grace handler")
  });

  const [get, put, status, key, keyAll] = useCoreMapStorage(
    current,
    cube,
    min_scale,
    graceHandler.current,
    add,
    sub
  );

  const [error, reload, bulkReload] = useServerMap(
    key,
    keyAll,
    put,
    graceHandler.current
  );

  const length = useMemo(() => 1 << (current[S] - min_scale), [
    current,
    min_scale
  ]);

  const allMap = useMemo<Chunk[][]>(() => {
    const result: Chunk[][] = [];
    for (var y = 0; y < length; y++) {
      const row: Chunk[] = [];
      for (var x = 0; x < length; x++) {
        row.push([x, y, current[S], current[D]]);
      }
      result.push(row);
    }
    return result;
  }, [length, current]);

  const info = useCallback(
    (ch: Chunk) => {
      const raw = status(hash(ch, min_scale));
      if (raw !== undefined) {
        switch (raw) {
          case FlashStatus.PRIMARY_ACTIVE:
            return "A";
          case FlashStatus.PRIMARY_GRACE:
            return "F";
          case FlashStatus.SECONDARY_ACTIVE:
            return "C";
          case FlashStatus.SECONDARY_GRACE:
            return "D";
        }
      }
      return "-";
    },
    [status, min_scale]
  );

  return useMemo(
    () => (
      <Container className={classes.table}>
        <Button
          variant="outlined"
          color="primary"
          onClick={() => {
            console.info("update!");
            const x = Math.random() * length * (1 << min_scale);
            const y = Math.random() * length * (1 << min_scale);
            console.info(`move to ${x} ${y} ${coordS}`);
            update(x, y, coordS);
          }}
        >
          移動
        </Button>
        <IconButton
          onClick={() => {
            console.info(`scale=${coordS} => ${coordS - 1}`);
            update(coordX, coordY, coordS - 1);
          }}
        >
          <ZoomInIcon />
        </IconButton>
        <IconButton
          onClick={() => {
            console.info(`scale=${coordS} => ${coordS + 1}`);
            update(coordX, coordY, coordS + 1);
          }}
        >
          <ZoomOutIcon />
        </IconButton>
        <p>updated at {new Date().toString()}</p>
        <p>coord={`(${coordX}, ${coordY}, ${coordS})`}</p>
        {allMap.map((row, idx) => (
          <Box className={classes.row} key={idx} display="flex">
            {row.map(cell => (
              <Paper className={classes.cell} key={hash(cell, min_scale)}>
                <div>{info(cell)}</div>
                <div>
                  {hash(current, min_scale) == hash(cell, min_scale) && "HERE"}
                </div>
              </Paper>
            ))}
          </Box>
        ))}
      </Container>
    ),
    [coordX, coordY, coordS, current, length, allMap, update, min_scale, core]
  );
};

export default CacheMap;
