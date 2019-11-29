import * as PIXI from "pixi.js";
import { Monitorable, MonitorContainer } from "interfaces/monitor";
import {
  SpriteContainerProperty,
  AnimatedSpriteProperty,
  AnimatedSpriteContainerProperty,
  SpriteProperty
} from "interfaces/pixi";
import { ZoomablePointModel, ZoomablePointModelContainer } from "./zoom";

/**
 * x, y座標はPointModelで初期化済み
 */
const defaultValues: {
  spscale: number;
  alpha: number;
  tint: number;
} = {
  spscale: 0.5,
  alpha: 1,
  tint: 0xffffff
};

export abstract class SpriteModel extends ZoomablePointModel
  implements Monitorable {
  sprite: PIXI.Sprite;

  constructor(options: SpriteProperty) {
    super(options);
    this.sprite = new PIXI.Sprite(options.texture);
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }

  setupBeforeCallback() {
    super.setupBeforeCallback();
    this.addBeforeCallback(() => {
      this.sprite.anchor.set(0.5, 0.5);
      this.sprite.alpha = this.props.alpha;
      this.sprite.scale.set(this.props.spscale, this.props.spscale);
      this.container.addChild(this.sprite);
    });
  }

  setupUpdateCallback() {
    super.setupUpdateCallback();

    // 値の更新時、Spriteを更新するように設定
    this.addUpdateCallback("tint", (tint: number) => (this.sprite.tint = tint));
    this.addUpdateCallback("spscale", (value: number) =>
      this.sprite.scale.set(value, value)
    );
    this.addUpdateCallback(
      "alpha",
      (value: number) => (this.sprite.alpha = value)
    );
  }

  setupAfterCallback() {
    super.setupAfterCallback();
    this.addAfterCallback(() => this.container.removeChild(this.sprite));
  }

  updateDisplayInfo() {
    super.updateDisplayInfo();
    this.setDisplayPosition();
  }

  protected getPIXIObject() {
    return this.sprite;
  }
}

export abstract class SpriteContainer<
  T extends SpriteModel,
  C extends SpriteProperty
> extends ZoomablePointModelContainer<T, C> implements MonitorContainer {
  protected texture: PIXI.Texture;

  constructor(
    options: SpriteContainerProperty,
    newInstance: { new (props: C): T }
  ) {
    super(options, newInstance);
    this.texture = options.texture;
  }

  protected getBasicChildOptions(): SpriteProperty {
    return {
      ...super.getBasicChildOptions(),
      texture: this.texture
    };
  }
}

export abstract class AnimatedSpriteModel extends SpriteModel
  implements Monitorable {
  constructor(options: AnimatedSpriteProperty) {
    super({ texture: PIXI.Texture.EMPTY, ...options });

    let sprite = new PIXI.AnimatedSprite(options.animation);
    sprite.gotoAndPlay(options.offset);
    this.sprite = sprite;
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues({ spscale: 1.0 });
  }
}

const animateDefaultValues: { offset: number } = { offset: 0 };

export abstract class AnimatedSpriteContainer<
  T extends AnimatedSpriteModel,
  C extends AnimatedSpriteProperty
> extends SpriteContainer<T, C> implements MonitorContainer {
  protected animation: PIXI.Texture[];
  protected offset: number;

  constructor(
    options: AnimatedSpriteContainerProperty,
    newInstance: { new (props: C): T }
  ) {
    super({ texture: PIXI.Texture.EMPTY, ...options }, newInstance);
    this.offset = 0;
    this.animation = options.animation;
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(animateDefaultValues);
  }
  setupUpdateCallback() {
    super.setupUpdateCallback();
    this.addUpdateCallback("offset", v => (this.offset = v));
  }

  protected getBasicChildOptions(): AnimatedSpriteProperty {
    return {
      ...super.getBasicChildOptions(),
      animation: this.animation,
      offset: this.offset
    };
  }
}
