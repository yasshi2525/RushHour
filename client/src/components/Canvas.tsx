import React, {
  useState,
  useRef,
  useEffect,
  useContext,
  useCallback
} from "react";
import makeStyles from "@material-ui/styles/makeStyles";
import createStyles from "@material-ui/styles/createStyles";
import ModelContext from "common/model";
import { useResize } from "common/resize";
import { CanvasProperty } from "common/interfaces";
import { MouseDragHandler, TouchDragHandler } from "common/handlers/drag";
import ResizeHandler from "common/handlers/window";
import { WheelHandler } from "common/handlers/wheel";
import { PinchHandler } from "common/handlers/pinch";
import { fetchMap, players } from "actions";
import { RushHourStatus } from "state";
import { ClickCursor, TapCursor } from "common/handlers/cursor";

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
  const [ww, wh] = useResize();
  useEffect(() => {
    if (model.resize(ww, wh)) {
    }
  }, [ww, wh]);
  useEffect(() => {}, [reDelegate]);

  this.mouse = new MouseDragHandler(this.props.model, this.props.dispatch);
  this.wheel = new WheelHandler(this.props.model, this.props.dispatch);
  this.touch = new TouchDragHandler(this.props.model, this.props.dispatch);
  this.pinch = new PinchHandler(this.props.model, this.props.dispatch);
  this.clickCursor = new ClickCursor(this.props.model, this.props.dispatch);
  this.tapCursor = new TapCursor(this.props.model, this.props.dispatch);
  this.resize = new ResizeHandler(this.props.model, this.props.dispatch);

  const divRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (divRef.current !== null) {
      divRef.current.appendChild(model.app.view);
    }
  }, [divRef]);
  const onMouseDown = useCallback(e => {}, [divRef]);

  return (
    <div
      ref={divRef}
      className={classes.canvas}
      onMouseDown={e => {
        this.clickCursor.onStart(e);
        this.mouse.onStart(e);
      }}
      onMouseMove={e => {
        this.clickCursor.onMove(e);
        this.mouse.onMove(e);
      }}
      onMouseUp={e => {
        this.clickCursor.onEnd(e);
        this.mouse.onEnd(e);
      }}
      onMouseOut={e => {
        this.clickCursor.onOut(e);
        this.mouse.onEnd(e);
      }}
      onWheel={e => {
        this.wheel.onStart(e);
        this.wheel.onMove(e);
        this.wheel.onEnd(e);
      }}
      onTouchStart={e => {
        this.tapCursor.onStart(e);
        this.touch.onStart(e);
        this.pinch.onStart(e);
      }}
      onTouchMove={e => {
        this.tapCursor.onMove(e);
        this.touch.onMove(e);
        this.pinch.onMove(e);
      }}
      onTouchEnd={e => {
        this.tapCursor.onEnd(e);
        this.touch.onEnd(e);
        this.pinch.onEnd(e);
      }}
    ></div>
  );
};

// Pixi.js が作成する canvas を管理するコンポーネント
class Canvas extends React.Component<CanvasProperty, RushHourStatus> {
  ref: React.RefObject<HTMLDivElement>;
  mouse: MouseDragHandler;
  wheel: WheelHandler;
  touch: TouchDragHandler;
  pinch: PinchHandler;
  clickCursor: ClickCursor;
  tapCursor: TapCursor;
  resize: ResizeHandler;

  constructor(props: CanvasProperty) {
    super(props);

    this.ref = React.createRef<HTMLDivElement>();

    this.mouse = new MouseDragHandler(this.props.model, this.props.dispatch);
    this.wheel = new WheelHandler(this.props.model, this.props.dispatch);
    this.touch = new TouchDragHandler(this.props.model, this.props.dispatch);
    this.pinch = new PinchHandler(this.props.model, this.props.dispatch);
    this.clickCursor = new ClickCursor(this.props.model, this.props.dispatch);
    this.tapCursor = new TapCursor(this.props.model, this.props.dispatch);
    this.resize = new ResizeHandler(this.props.model, this.props.dispatch);
  }

  componentDidMount() {
    if (this.ref.current !== null) {
      if (this.props.model.app.view instanceof Node) {
        // 単体テスト時、Node非実装のため
        // 一度描画して、canvas要素を子要素にする
        this.ref.current.appendChild(this.props.model.app.view);
      }
    }
  }

  componentDidUpdate() {
    if (!this.props.isPlayerFetched) {
      this.props.dispatch(players.request({ model: this.props.model }));
    }
    if (this.props.isFetchRequired) {
      this.props.dispatch(fetchMap.request({ model: this.props.model }));
    }
  }

  componentWillUnmount() {
    this.props.model.unmount();
  }
}

function mapStateToProps(state: RushHourStatus) {
  return {
    isFetchRequired: state.isFetchRequired,
    isPlayerFetched: state.isPlayerFetched
  };
}
