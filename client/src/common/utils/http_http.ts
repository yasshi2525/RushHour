import { SerializableObject } from "interfaces";
import {
  HttpMethod,
  Method,
  HttpEndpoint,
  HttpAuthEndpoint
} from "interfaces/endpoint";
import {
  withAuth,
  parse,
  useHttpCommon,
  useHttpCommonTask
} from "./http_common";
import {
  useMultiHttpCommon,
  MultiEndpoint,
  useMultiHttpCommonTask
} from "./http_multi";
import { TaskHandler } from "./task";
import { MultiTaskHandler } from "./task_multi";

const opts = (
  signal: AbortSignal,
  method: HttpMethod | undefined,
  auth: boolean,
  args: SerializableObject
): RequestInit => {
  const headers = new Headers();
  headers.set("Content-type", "application/json");
  if (auth) {
    withAuth(headers);
  }
  const body = Object.keys(args).length ? JSON.stringify(args) : undefined;
  return { headers, signal, method: method ? method : Method.POST, body };
};

/**
 * 指定されたURLににアクセスし、レスポンス値を解析する
 * @throws `***Error` `200`以外, レスポンスがJSONでない
 * @see parse
 */
export const httpHttp = async <
  I extends SerializableObject,
  O extends SerializableObject
>(
  ep: HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>,
  signal: AbortSignal
) =>
  parse(
    ep.payload,
    await fetch(
      ep.url,
      opts(signal, ep.method, "auth" in ep ? ep.auth : false, ep.args)
    )
  );

export const useHttpHttp = <
  I extends SerializableObject,
  O extends SerializableObject
>(
  endpoint: HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>
) => {
  return useHttpCommon<HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>, I, O>(
    endpoint,
    httpHttp
  );
};

export const useMultiHttpHttp = <
  I extends SerializableObject,
  O extends SerializableObject
>(
  eps: MultiEndpoint<HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>, I, O>
) => {
  return useMultiHttpCommon<HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>, I, O>(
    eps,
    httpHttp
  );
};

/**
 * Httpリクエストを任意のタイミングで実施する。
 * ```
 * const [fire, cancel] = useHttpGetTask(endpoint, {onOK:(d)=>{ok}, onError:(e)=>{ng}, onCancel:(c)=>{cancel});
 * fire(args);
 * cancel();
 * ```
 * @see `useTask`
 */
export const useHttpHttpTask = <
  I extends SerializableObject,
  O extends SerializableObject
>(
  ep: HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>,
  handlers: TaskHandler<I, O>
) =>
  useHttpCommonTask<HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>, I, O>(
    ep,
    httpHttp,
    handlers
  );

export const useMultiHttpHttpTask = <
  I extends SerializableObject,
  O extends SerializableObject
>(
  endpointList: HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>,
  handlers: MultiTaskHandler<I, O>
) =>
  useMultiHttpCommonTask<HttpEndpoint<I, O> | HttpAuthEndpoint<I, O>, I, O>(
    endpointList,
    httpHttp,
    handlers
  );
