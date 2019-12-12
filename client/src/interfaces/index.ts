export type Primitive = boolean | number | string;
export type FlatObject = {
  [index: string]: Primitive | Array<Primitive>;
  [index: number]: Primitive | Array<Primitive>;
};

export type SerializableObject = {
  [index: string]: Primitive | Array<Primitive> | Object | Array<Object>;
  [index: number]: Primitive | Array<Primitive> | Object | Array<Object>;
};

type Base = {
  [index: string]: Primitive | Array<Primitive> | SerializableObject;
};

export interface Entity extends Base {
  id: number;
}

type Position = { x: number; y: number };

export interface Locatable extends Entity {
  pos: Position;
  scale: number;
}

export type Delegatable = Locatable &
  ({ mul: 1; cid: number } | { mul: number }) &
  ({} | { pids: number[] });

export const EMPTY_DLG: Delegatable = {
  id: 0,
  pos: { x: 0, y: 0 },
  mul: 1,
  scale: 0
};

export type Dimension = [number, number];
export type Coordinates = [number, number, number];

export type HashContainer<T extends Entity> = { [index: number]: T };
