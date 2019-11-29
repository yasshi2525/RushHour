import { useEffect, useState } from "react";
import {
  OperationError,
  AuthError,
  RequestError,
  UnhandledError,
  ServerErrors,
  isServerErrors
} from "interfaces/error";
import { GeneralObject, Entity } from "interfaces";
import { Method, Endpoint } from "interfaces/endpoint";
import { useTask } from "./task";
import { useSnack } from "./snack";

interface ErrContents {
  err: string[];
}

/**
 * サーバはエラーがあるとき、`err` キーに 配列を入れて返す。
 * @param obj 判定対象のオブジェクト
 */
const isErrContents = (obj: any): obj is ErrContents => {
  return obj.err !== undefined && obj.err instanceof Array;
};

/**
 * エラー時のレスポンスボディからエラーメッセージを取り出す。
 * エラー時は`err`キーにメッセージが配列で格納されている。
 * @throws `UnhandledServerError` JSON形式でないとき、`err`キーがないときは解析に失敗したとき。
 */
const parseError = async (res: Response) =>
  await res
    .json()
    .then(data => (!isErrContents(data) ? [`parse error: ${data}`] : data.err))
    .catch(e => {
      throw new UnhandledError(`unknown server error: ${e.name}: ${e.message}`);
    });

/**
 * レスポンスボディのJSON文字列をオブジェクトに変換する。
 * @throws RequestError `400`
 * @throws AuthError `401`
 * @throws OperationError `503` `504`
 * @throws UnhandledServerError それ以外、エラーメッセージの解釈失敗
 */
const parse = async <T>(res: Response) => {
  switch (res.status) {
    case 200:
      try {
        return <T>await res.json();
      } catch (e) {
        throw new UnhandledError(e);
      }
    case 400:
      throw new RequestError(await parseError(res));
    case 401:
      throw new AuthError(await parseError(res));
    case 503:
    case 504:
      throw new OperationError(await parseError(res));
    default:
      throw new UnhandledError(`invalid response status ${res.status}`);
  }
};

/**
 * 指定されたURL, Methodにアクセスし、レスポンス値を解析する
 * @throws `***Error` `200`以外, レスポンスがJSONでない
 * @see parseResponse
 */
const _http = async <I extends GeneralObject = {}, O = {}>(
  e: Endpoint,
  signal: AbortSignal,
  args?: I
) => {
  const method = e.method !== undefined ? e.method : Method.GET;
  const headers = new Headers();
  if (e.auth) {
    const jwt = localStorage.getItem("jwt");
    if (jwt == null) {
      throw new AuthError("no jwt");
    }
    headers.set("Authorization", `Bearer ${jwt}`);
  }
  if (method === Method.GET) {
    let url = e.url;
    if (args) {
      let params = new URLSearchParams();
      Object.entries(args).forEach(([k, v]) => params.set(k, v.toString()));
      url = `${e.url}?${params.toString()}`;
    }
    return parse<O>(await fetch(url, { headers, signal }));
  } else {
    headers.set("Content-type", "application/json");
    return parse<O>(
      await fetch(e.url, {
        method,
        headers,
        body: JSON.stringify(args),
        signal
      })
    );
  }
};

export type Http<T = undefined> = [ServerErrors, undefined] | [undefined, T];

/**
 * 指定されたURL, Methodにアクセスし、レスポンス値を解析する
 * @see _http
 */
export const http = async <
  I extends GeneralObject = {},
  O extends GeneralObject | Entity[] = {}
>(
  ep: Endpoint,
  signal: AbortSignal,
  args?: I
): Promise<Http<O>> =>
  await _http<I, O>(ep, signal, args)
    .then(d => [undefined, d] as [undefined, O])
    .catch(e => {
      if (isServerErrors(e)) {
        console.warn(e);
        return [e, undefined];
      } else {
        throw e;
      }
    });

const initErr = new UnhandledError("under fetching...");

/**
 * `endpoint` にアクセスし、結果を返す。
 * 結果がすぐほしいとき使用する。
 * 特定のイベント時にほしいときは `useHttpTask`を使用する
 * ```
 * const [completed, error, data] = useFetch(endpoint, args);
 * ```
 */
export const useFetch = <I = {}, O = {}>(ep: Endpoint, args?: I) => {
  const [completed, setCompleted] = useState(false);
  const [contents, setContents] = useState<Http<O>>([initErr, undefined]);
  const snack = useSnack();
  useEffect(() => {
    console.info(`effect useFetch ${ep.url}`);
    const aborter = new AbortController();
    let setContentsSafe = (c: Http<O>) => setContents(c);
    (async () => {
      const result = await http<I, O>(ep, aborter.signal, args);
      if (result[0]) {
        snack(result[0]);
      }
      setContentsSafe(result);
      setCompleted(true);
    })();
    return () => {
      console.info(`cleanup useFetch ${ep.url}`);
      setContentsSafe = () => null;
      aborter.abort();
      setContents([initErr, undefined]);
      setCompleted(false);
    };
  }, []);
  return [completed, ...contents] as [
    boolean,
    ServerErrors | undefined,
    O | undefined
  ];
};

/**
 * Httpリクエストを任意のタイミングで実施する。
 * ```
 * const [fire, cancel] = useHttpTask(endpoint, (d)=>{ok}, (e)=>{ng});
 * fire(args);
 * cancel();
 * ```
 * @see `useTask`
 */
export const useHttpTask = <I, O>(
  ep: Endpoint,
  ok?: (d: O) => void,
  ng?: (e: ServerErrors) => void
) => {
  return useTask((sig, args: I) => _http(ep, sig, args), ok, ng);
};
