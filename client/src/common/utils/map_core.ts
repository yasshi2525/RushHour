import { useCallback, useRef, useState, useMemo } from "react";
import { Locatable, HashContainer } from "interfaces";
import {
  IFetchMapResponseKeys,
  FetchMapResponseKeys
} from "interfaces/endpoint";

export interface Pluggnable extends Locatable {
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

const add = (body: CoreMap, key: string, value: CoreMap) => {
  FetchMapResponseKeys.forEach(i => {
    Object.values(value[i]).forEach(obj => {
      if (obj.id in body[i]) {
        body[i][obj.id].memberOf.push(key);
      } else {
        obj.memberOf.push(key);
        body[i][obj.id] = obj;
      }
    });
  });
  body.timestamp = Math.max(body.timestamp, value.timestamp);
  return body;
};

const sub = (prev: CoreMap, hash: string) => {
  FetchMapResponseKeys.forEach(key => {
    const garbage: number[] = [];
    Object.values(prev[key])
      .filter(obj => obj.memberOf.includes(hash))
      .forEach(obj => {
        obj.memberOf = obj.memberOf.filter(key => key != hash);
        if (!obj.memberOf.length) {
          garbage.push(obj.id);
        }
      });
    garbage.forEach(id => delete prev[key][id]);
  });
};

type Handlers = [
  CoreMap,
  (k: string, v: CoreMap) => void,
  (key: string) => void
];

/**
 * const [data, add, sub, rev] = useCoreMapReduce();
 */
const useCoreMap = (): Handlers => {
  const storage = useRef(newInstance());
  const [rev, setRev] = useState(0);
  const notify = useCallback(() => {
    setRev(v => v + 1);
  }, []);
  const include = useCallback(
    (k: string, v: CoreMap) => {
      add(storage.current, k, v);
      notify();
    },
    [storage]
  );
  const remove = useCallback(
    (key: string) => {
      sub(storage.current, key);
      notify();
    },
    [storage]
  );
  const core = useMemo(() => Object.assign({}, storage.current), [
    storage,
    rev
  ]);
  return [core, include, remove];
};

export default useCoreMap;
