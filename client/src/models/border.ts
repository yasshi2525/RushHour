import { GraphicsModel, GraphicsContainer } from "./graphics";
import { Monitorable, MonitorContainer } from "interfaces/monitor";
import { BorderProperty, PIXIProperty } from "interfaces/pixi";
import {
  config,
  Coordinates,
  Chunk,
  getChunkByPos,
  getChunkByScale
} from "interfaces/gamemap";

const graphicsOpts = {
  color: 0x9e9e9e,
  lightWidth: 1,
  boldWidth: 2
};

const defaultValues: {
  index: number;
  pos: number;
  scale: number;
  delegate: number;
} = {
  index: 0,
  pos: 0,
  scale: config.scale.default,
  delegate: 0
};

interface BorderChildProperty extends PIXIProperty {
  v: boolean;
}

interface BorderContainerProperty extends BorderChildProperty {
  delegate: number;
}

export class NormalBorder extends GraphicsModel implements Monitorable {
  protected v: boolean;
  protected currentAlpha: number;
  protected destinationAlpha: number;

  constructor(props: BorderChildProperty) {
    super(props);
    this.v = props.v;
    this.currentAlpha = 1;
    this.destinationAlpha = 1;
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }

  setInitialValues(props: { [index: string]: any }) {
    super.setInitialValues(props);
    if (props.v) {
      this.props.pos = { x: props.pos, y: 0 };
    } else {
      this.props.pos = { x: 0, y: props.pos };
    }
  }

  setupBeforeCallback() {
    super.setupBeforeCallback();
    this.addBeforeCallback(() => {
      this.shape();
      this.updateDestination();
      this.moveDestination();
      this.updateDisplayInfo();
    });
  }

  setupUpdateCallback() {
    super.setupUpdateCallback();
    this.addUpdateCallback("resize", () => this.shape());
  }

  updateDestination() {
    super.updateDestination();
    this.destinationAlpha =
      1 -
      Math.abs(
        this.props.coord.scale - this.props.scale - this.props.delegate + 1
      );
    this.destinationAlpha *= 2;
  }

  moveDestination() {
    super.moveDestination();
    this.currentAlpha = this.destinationAlpha;
  }

  protected mapRatioToVariable(ratio: number) {
    super.mapRatioToVariable(ratio);
    this.currentAlpha =
      this.currentAlpha * ratio + this.destinationAlpha * (1 - ratio);
  }

  protected shape() {
    this.graphics.clear();
    let width =
      this.props.index % 2 === 0
        ? graphicsOpts.boldWidth
        : graphicsOpts.lightWidth;
    this.graphics.lineStyle(width, graphicsOpts.color);
    this.graphics.moveTo(0, 0);
    if (this.v) {
      this.graphics.lineTo(0, this.app.renderer.height);
    } else {
      this.graphics.lineTo(this.app.renderer.width, 0);
    }
  }

  updateDisplayInfo() {
    super.updateDisplayInfo();
    if (this.v) {
      this.graphics.y = 0;
    } else {
      this.graphics.x = 0;
    }
    this.graphics.alpha = this.currentAlpha;
  }
}

const borderContainerDefaultValues: {
  coord: Coordinates;
  [index: string]: any;
} = {
  coord: {
    cx: config.gamePos.default.x,
    cy: config.gamePos.default.y,
    scale: config.scale.default,
    zoom: 0
  },
  resize: false,
  delegate: 0
};

abstract class NormalBorderContainer
  extends GraphicsContainer<NormalBorder, BorderChildProperty>
  implements MonitorContainer {
  protected chunk: Chunk;
  protected v: boolean;
  protected chScale: { [index: number]: boolean };

  constructor(props: BorderContainerProperty) {
    super(props, NormalBorder);
    this.chScale = {};
    this.v = props.v;
    this.props.delegate = props.delegate;

    this.chunk = getChunkByPos(
      config.gamePos.default,
      config.scale.default - props.delegate + 1
    );
    this.genChildren(getChunkByScale(this.chunk, -1));
    this.genChildren(this.chunk);
    this.genChildren(getChunkByScale(this.chunk, +1));
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(borderContainerDefaultValues);
  }

  setInitialValues(props: { [index: string]: any }) {
    super.setInitialValues(props);
    this.props.delegate = this.model.delegate;
  }

  setupUpdateCallback() {
    super.setupUpdateCallback();
    this.addUpdateCallback("delegate", () => this.refreshChildren());
    this.addUpdateCallback("coord", (v: Coordinates) => {
      let nowChunk = getChunkByPos(
        { x: v.cx, y: v.cy },
        v.scale - this.props.delegate + 1
      );

      if (this.chunk.scale !== nowChunk.scale) {
        this.removeOutRangeChildren();
        let zoom = nowChunk.scale < this.chunk.scale;
        if (zoom) {
          for (var s = -1; s < this.chunk.scale - nowChunk.scale; s++) {
            if (!this.chScale[nowChunk.scale + s]) {
              this.genChildren(getChunkByScale(nowChunk, s));
            }
          }
        } else {
          for (var s = +1; s > this.chunk.scale - nowChunk.scale; s--) {
            if (!this.chScale[nowChunk.scale + s]) {
              this.genChildren(getChunkByScale(nowChunk, s));
            }
          }
        }
        this.chunk = nowChunk;
        Object.keys(this.chScale).forEach(scale => {
          let ch = getChunkByPos({ x: v.cx, y: v.cy }, parseInt(scale));
          this.genChildren(ch);
        });
      } else {
        for (var idx = -1; idx <= 1; idx++) {
          let beforeChunk = getChunkByScale(this.chunk, idx);
          let afterChunk = getChunkByScale(nowChunk, idx);
          let beforeOffset = this.getOffset(beforeChunk);
          let afterOffset = this.getOffset(afterChunk);

          let num = Math.pow(2, this.props.delegate - idx + 1);
          // 左(上)側を作成、右(下)側を削除
          for (var i = 0; i < beforeOffset - afterOffset; i++) {
            var offset = i - Math.floor(num / 2) - 1;
            this.genChild(beforeChunk, offset);

            var offset = i + Math.floor(num / 2) + beforeOffset;
            this.removeChild(this.getId(offset, beforeChunk.scale));
          }
          // 右(下)側を作成、左(上)側を削除
          for (var i = 0; i < afterOffset - beforeOffset; i++) {
            var offset = i + Math.floor(num / 2) + 1;
            this.genChild(beforeChunk, offset);

            var offset = i - Math.floor(num / 2) + beforeOffset;
            this.removeChild(this.getId(offset, beforeChunk.scale));
          }
        }
      }

      this.chunk = nowChunk;
    });
  }

  protected getBasicChildOptions(): BorderChildProperty {
    return {
      ...super.getBasicChildOptions(),
      v: this.v
    };
  }

  protected getChildOptions() {
    return this.getBasicChildOptions();
  }

  protected getInterval(chunk: Chunk) {
    return Math.pow(2, chunk.scale);
  }

  protected getId(offset: number, scale: number) {
    return offset + scale * config.scale.max;
  }

  protected abstract getOffset(chunk: Chunk): number;

  protected abstract isAreaIn(offset: number): boolean;

  protected genChildOpts(
    id: number,
    index: number,
    pos: number,
    scale: number
  ) {
    return {
      v: this.v,
      id: id,
      index: index,
      pos: pos,
      scale: scale,
      coord: this.props.coord,
      delegate: this.props.delegate
    };
  }

  protected genChild(chunk: Chunk, offset: number) {
    let index = this.getOffset(chunk) + offset;
    let id = this.getId(index, chunk.scale);
    let pos = index * this.getInterval(chunk);

    if (!this.isAreaIn(pos)) {
      return;
    }

    if (this.existsChild(id)) {
      return;
    }

    this.addChild(this.genChildOpts(id, index, pos, chunk.scale));
  }

  protected genChildren(chunk: Chunk) {
    this.chScale[chunk.scale] = true;
    let num = Math.pow(
      2,
      this.props.delegate + this.chunk.scale - chunk.scale + 1
    );
    for (var pos = -Math.floor(num / 2); pos < Math.floor(num / 2) + 1; pos++) {
      this.genChild(chunk, pos);
    }
  }

  protected refreshChildren() {
    this.forEachChild(c => this.removeChild(c.get("id")));
    this.chunk = getChunkByPos(
      { x: this.props.coord.cx, y: this.props.coord.cy },
      this.props.coord.scale - this.props.delegate + 1
    );

    this.genChildren(getChunkByScale(this.chunk, -1));
    this.genChildren(this.chunk);
    this.genChildren(getChunkByScale(this.chunk, +1));
  }

  protected removeOutRangeChildren() {
    this.forEachChild(c => {
      if (
        Math.abs(
          this.props.coord.scale - c.get("scale") - this.props.delegate + 1
        ) > 1
      ) {
        delete this.chScale[c.get("scale")];
        this.removeChild(c.get("id"));
      }
    });
  }
}

export class XBorderContainer extends NormalBorderContainer
  implements MonitorContainer {
  constructor(props: BorderProperty) {
    super({ ...props, v: false });
  }

  protected getOffset(chunk: Chunk) {
    return chunk.y;
  }

  protected isAreaIn(offset: number) {
    return offset > config.gamePos.min.x && offset < config.gamePos.max.x;
  }
}

export class YBorderContainer extends NormalBorderContainer
  implements MonitorContainer {
  constructor(props: BorderProperty) {
    super({ ...props, v: true });
  }

  protected getOffset(chunk: Chunk) {
    return chunk.x;
  }

  protected isAreaIn(offset: number) {
    return offset > config.gamePos.min.y && offset < config.gamePos.max.y;
  }
}
