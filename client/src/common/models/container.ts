import BaseModel, { BaseProperty } from "./base";
import { MonitorContainer, Monitorable } from "common/interfaces/monitor";
import { Chunk, ResolveError } from "common/interfaces/gamemap";
import { PointModel } from "./point";

export default abstract class<T extends Monitorable, C extends BaseProperty>
  extends BaseModel
  implements MonitorContainer {
  Child: { new (props: C): T };
  children: { [index: string]: T };

  constructor(opts: BaseProperty, newInstance: { new (props: C): T }) {
    super({ model: opts.model });
    this.Child = newInstance;
    this.children = {};
  }

  existsChild(id: string) {
    return this.children[id] !== undefined;
  }

  getChild(id: string) {
    return this.children[id];
  }

  getChildOnChunk(chunk: Chunk, oid: number): PointModel | undefined {
    let result = Object.keys(this.children)
      .map(id => this.children[id])
      .find(c => c.get("oid") === oid && c.standOnChunk(chunk));
    return result instanceof PointModel ? result : undefined;
  }

  protected getBasicChildOptions(): BaseProperty {
    return { model: this.model };
  }

  protected abstract getChildOptions(): C;

  addChild(props: { id: string; [propName: string]: any }): T {
    let child = new this.Child(this.getChildOptions());
    child.setupDefaultValues();
    child.setupUpdateCallback();
    child.setupBeforeCallback();
    child.setupAfterCallback();
    child.setInitialValues({ ...props, coord: this.model.coord });
    child.begin();

    this.children[props.id] = child;
    this.change();

    return child;
  }

  updateChild(props: { id: string; [propName: string]: any }): T {
    let target = this.children[props.id];

    target.mergeAll(props);

    if (target.isChanged()) {
      this.change();
    }
    return target;
  }

  endChildren() {
    this.forEachChild(ch => {
      if (ch.shouldEnd()) {
        this.removeChild(ch.get("id"));
      }
    });
  }

  removeChild(id: string) {
    if (this.existsChild(id)) {
      this.children[id].end();
      delete this.children[id];
      this.change();
    }
  }

  mergeChild(props: { id: string; [propName: string]: any }): T {
    if (this.existsChild(props.id)) {
      return this.updateChild(props);
    } else {
      return this.addChild(props);
    }
  }

  mergeChildren(
    payload: { id: string; [propName: string]: any }[],
    opts: { [index: string]: any }
  ) {
    if (payload === undefined) {
      return;
    }
    payload.forEach(props => {
      Object.assign(props, opts, { outMap: false });
      this.mergeChild(props);
    });

    // payloadに含まれない child に outMap をつける
    let aliveIds = payload.map(props => props.id);
    Object.keys(this.children)
      .filter(myId => !aliveIds.find(id => myId == id))
      .filter(id => !this.children[id].get("deamon"))
      .forEach(id => this.getChild(id).merge("outMap", true));
  }

  /**
   * すべてのchildのkeyにvalueを設定します
   * @param key 設定するプロパティ名
   * @param value プロパティに設定する値
   */
  merge(key: string, value: any) {
    super.merge(key, value);
    this.forEachChild(c => {
      c.merge(key, value);
      if (c.isChanged()) {
        this.change();
      }
    });
  }

  reset() {
    this.forEachChild(c => c.reset());
    super.reset();
  }

  end() {
    this.forEachChild(c => c.end());
    super.end();
  }

  resolve(error: ResolveError) {
    this.forEachChild(v => v.resolve(error));
    return error;
  }

  forEachChild(func: (c: T) => any) {
    Object.keys(this.children).forEach(id => func(this.children[id]));
  }
}
