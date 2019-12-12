import { FlatObject, SerializableObject } from "interfaces";
import { GetEndpoint, GetAuthEndpoint } from "interfaces/endpoint";
import {
  withAuth,
  parse,
  useHttpCommon,
  useHttpCommonTask
} from "./http_common";
import {
  useMultiHttpCommon,
  useMultiHttpCommonTask,
  MultiEndpoint
} from "./http_multi";
import { TaskHandler } from "./task";
import { MultiTaskHandler } from "./task_multi";

const toQuery = (url: string, args: FlatObject) => {
  if (Object.keys(args).length) {
    return url;
  }
  const builder = new URLSearchParams();
  Object.entries(args).forEach(([k, v]) => builder.set(k, v.toString()));
  return `${url}?${builder.toString()}`;
};

const opts = (signal: AbortSignal, auth: boolean): RequestInit => {
  const headers = new Headers();
  if (auth) {
    withAuth(headers);
  }
  return { headers, signal };
};

/**
 * 指定されたURLにアクセスし、レスポンス値を解析する
 * @throws `***Error` `200`以外, レスポンスがJSONでない
 * @see parse
 */
export const httpGet = async <
  I extends FlatObject,
  O extends SerializableObject
>(
  ep: GetEndpoint<I, O> | GetAuthEndpoint<I, O>,
  signal: AbortSignal
) =>
  parse(
    ep.payload,
    await fetch(
      toQuery(ep.url, ep.args),
      opts(signal, "auth" in ep ? ep.auth : false)
    )
  );

export const useHttpGet = <I extends FlatObject, O extends SerializableObject>(
  endpoint: GetEndpoint<I, O> | GetAuthEndpoint<I, O>
) => {
  return useHttpCommon<GetEndpoint<I, O> | GetAuthEndpoint<I, O>, I, O>(
    endpoint,
    httpGet
  );
};

export const useMultiHttpGet = <
  I extends FlatObject,
  O extends SerializableObject
>(
  endpointList: MultiEndpoint<GetEndpoint<I, O> | GetAuthEndpoint<I, O>, I, O>
) =>
  useMultiHttpCommon<GetEndpoint<I, O> | GetAuthEndpoint<I, O>, I, O>(
    endpointList,
    httpGet
  );

/**
 * Httpリクエストを任意のタイミングで実施する。
 * ```
 * const [fire, cancel] = useHttpGetTask(endpoint, {onOK:(d)=>{ok}, onError:(e)=>{ng}, onCancel:(c)=>{cancel});
 * fire(args);
 * cancel();
 * ```
 * @see `useTask`
 */
export const useHttpGetTask = <
  I extends FlatObject,
  O extends SerializableObject
>(
  endpoint: GetEndpoint<I, O> | GetAuthEndpoint<I, O>,
  handlers: TaskHandler<I, O>
) =>
  useHttpCommonTask<GetEndpoint<I, O> | GetAuthEndpoint<I, O>, I, O>(
    endpoint,
    httpGet,
    handlers
  );

export const useMultiHttpGetTask = <
  I extends FlatObject,
  O extends SerializableObject
>(
  endpointList: GetEndpoint<I, O> | GetAuthEndpoint<I, O>,
  handlers: MultiTaskHandler<I, O>
) =>
  useMultiHttpCommonTask<GetEndpoint<I, O> | GetAuthEndpoint<I, O>, I, O>(
    endpointList,
    httpGet,
    handlers
  );
