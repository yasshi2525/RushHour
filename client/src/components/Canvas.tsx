import React, { useRef, useEffect, useContext } from "react";
import makeStyles from "@material-ui/styles/makeStyles";
import createStyles from "@material-ui/styles/createStyles";
import ModelContext from "common/model";
import { useDrag, useSwipe } from "common/scroll";
import { useWheel, usePinch } from "common/zoom";
import { useResize } from "common/resize";
import { useCursor } from "common/cursor";
import LoadingContext, { LoadingStatus } from "common/loading";

const useStyles = makeStyles(() =>
  createStyles({
    canvas: {
      position: "fixed",
      left: "0px",
      top: "0px",
      width: "100%",
      height: "100%"
    }
  })
);

export default () => {
  const classes = useStyles();
  const model = useContext(ModelContext);
  useResize();
  const drag = useDrag();
  const onWheel = useWheel();
  const swipe = useSwipe();
  const pinch = usePinch();
  const cursor = useCursor();
  const divRef = useRef<HTMLDivElement>(null);
  const { update } = useContext(LoadingContext);

  useEffect(() => {
    if (divRef.current !== null) {
      divRef.current.appendChild(model.app.view);
      update(LoadingStatus.INITED_CONTROLLER);
      update(LoadingStatus.END);
    }
  }, [divRef]);

  return (
    <div
      ref={divRef}
      className={classes.canvas}
      onMouseDown={e => {
        cursor.onMouseDown(e);
        drag.onMouseDown(e);
      }}
      onMouseMove={e => {
        cursor.onMouseMove(e);
        drag.onMouseMove(e);
      }}
      onMouseUp={e => {
        cursor.onMouseUp(e);
        drag.onMouseUp(e);
      }}
      onMouseOut={e => {
        cursor.onMouseOut(e);
        drag.onMouseOut(e);
      }}
      onWheel={e => onWheel(e)}
      onTouchStart={e => {
        cursor.onTouchStart(e);
        swipe.onTouchStart(e);
        pinch.onTouchStart(e);
      }}
      onTouchMove={e => {
        cursor.onTouchMove(e);
        swipe.onTouchMove(e);
        pinch.onTouchMove(e);
      }}
      onTouchEnd={e => {
        cursor.onTouchEnd(e);
        swipe.onTouchEnd(e);
        pinch.onTouchEnd(e);
      }}
    ></div>
  );
};
