import { useCallback, useMemo } from "react";
import { FetchMap, FetchMapResponse, fetchMap } from "interfaces/endpoint";
import { useHttpGetTask, useMultiHttpGetTask } from "./http_get";
import { OkResponse } from "./http_common";
import { Chunk } from "./map_core";
import { encode } from "./map_server";

type Handlers = [() => void, () => void];

/**
 * ```
 * const [reload, bulkReload] = useServerMapReload(current, keys, put, putAll);
 * ```
 */
const useServerMapReload = (
  currentChunk: Chunk,
  keys: Chunk[],
  insert: (payload: OkResponse<FetchMap, FetchMapResponse>) => void,
  bulkInsert: (payloadList: OkResponse<FetchMap, FetchMapResponse>[]) => void
): Handlers => {
  const handler = useMemo(
    () => ({
      onOK: (payload: FetchMapResponse, args: FetchMap) =>
        insert({ payload, args })
    }),
    [insert]
  );

  const [fetchTask] = useHttpGetTask(fetchMap, handler);

  const reload = useCallback(() => {
    fetchTask(encode(currentChunk));
  }, [fetchTask, currentChunk]);

  const bulkHandler = useMemo(() => ({ onOK: bulkInsert }), [bulkInsert]);

  const [bulkFetch] = useMultiHttpGetTask(fetchMap, bulkHandler);

  const bulkReload = useCallback(() => {
    bulkFetch(keys.map(encode));
  }, [bulkFetch, keys]);

  return [reload, bulkReload];
};

export default useServerMapReload;
