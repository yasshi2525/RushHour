import { useCallback, useMemo } from "react";
import { Locatable, HashContainer } from "interfaces";
import {
  IFetchMapResponseKeys,
  FetchMapResponseKeys
} from "interfaces/endpoint";

export type Chunk = [number, number, number, number];

interface Pluggnable extends Locatable {
  scale: number;
  memberOf: string[];
}

export type CoreMap = {
  [index in IFetchMapResponseKeys]: HashContainer<Pluggnable>;
} & {
  timestamp: number;
};

export const newInstance = (): CoreMap => ({
  residences: {},
  companies: {},
  rail_nodes: {},
  rail_edges: {},
  timestamp: 0
});

const isMapNotNull = (obj: any): obj is CoreMap => obj;

const reduce = (prev: CoreMap, current: CoreMap) => {
  FetchMapResponseKeys.forEach(key => {
    Object.values(current[key]).forEach(obj => {
      prev[key][obj.id] = {
        ...obj,
        memberOf: [...prev[key][obj.id].memberOf, ...obj.memberOf]
      };
    });
  });
  prev.timestamp = Math.max(prev.timestamp, current.timestamp);
  return prev;
};

type Handlers = CoreMap;

/**
 * ```
 * const coreMap = useCoreMap(get, keys);
 * ```
 */
const useCoreMap = (
  get: (key: Chunk) => CoreMap | undefined,
  keys: Chunk[]
): Handlers => {
  /**
   * target 以内にあるデータを一つにまとめる
   */
  const combine = useCallback(
    (init: CoreMap) =>
      keys
        .map(key => get(key))
        .filter(isMapNotNull)
        .reduce(reduce, init),
    [get, keys]
  );

  return useMemo(() => combine(newInstance()), [combine]);
};

export default useCoreMap;
