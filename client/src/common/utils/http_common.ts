import { useEffect, useState, useContext } from "react";
import {
  OperationError,
  AuthError,
  RequestError,
  ServerError,
  Errors,
  isErrors,
  CancelError,
  isOperationError
} from "interfaces/error";
import useTask, { TaskHandler } from "./task";
import OperationContext from "./operation";

interface ErrContents {
  err: string[];
}

/**
 * サーバはエラーがあるとき、`err` キーに 配列を入れて返す。
 * @param obj 判定対象のオブジェクト
 */
const isErrContents = (obj: any): obj is ErrContents => {
  return obj?.err instanceof Array;
};

/**
 * エラー時のレスポンスボディからエラーメッセージを取り出す。
 * エラー時は`err`キーにメッセージが配列で格納されている。
 * @throws `ServerError` JSON形式でないとき、`err`キーがないときは解析に失敗したとき。
 */
const parseError = async (res: Response) => {
  try {
    const data = await res.json();
    if (!isErrContents(data)) {
      throw new ServerError(`parse error: ${data}`);
    } else {
      return data.err;
    }
  } catch (e) {
    if (!isErrors(e)) {
      throw new ServerError(`unknown server error: ${e}`);
    } else {
      throw e;
    }
  }
};

/**
 * レスポンスボディのJSON文字列をオブジェクトに変換する。
 * @throws RequestError `400`
 * @throws AuthError `401`
 * @throws OperationError `503` `504`
 * @throws ServerError それ以外、エラーメッセージの解釈失敗
 */
export const parse = async <T>(_: T, res: Response) => {
  switch (res.status) {
    case 200:
      try {
        return <T>await res.json();
      } catch (e) {
        throw new ServerError(e);
      }
    case 400:
      throw new RequestError(await parseError(res));
    case 401:
      throw new AuthError(await parseError(res));
    case 503:
    case 504:
      throw new OperationError(await parseError(res));
    default:
      throw new ServerError(`invalid response status ${res.status}`);
  }
};

export const withAuth = (headers: Headers) => {
  const jwt = localStorage.getItem("jwt");
  if (jwt == null) {
    throw new AuthError("no jwt");
  }
  headers.set("Authorization", `Bearer ${jwt}`);
};

export type OkResponse<I, O> = { args: I; payload: O };
export type ErrorResponse<I> = { args: I; error: Errors };

export type ResponseType<I, O> = OkResponse<I, O> | ErrorResponse<I>;

export const isErrorResponse = <I, O>(
  obj: ResponseType<I, O>
): obj is ErrorResponse<I> => "error" in obj && isErrors(obj.error);

export type Endpoint<E, I, O> = E & { args: I; payload: O };
export type RequestTask<E, I, O> = (
  endpoint: Endpoint<E, I, O>,
  signal: AbortSignal
) => Promise<O>;

export const httpCommon = async <E, I, O>(
  endpoint: Endpoint<E, I, O>,
  task: RequestTask<E, I, O>,
  signal: AbortSignal,
  maintain: (e: OperationError) => void
) => {
  try {
    return {
      args: endpoint.args,
      payload: await task(endpoint, signal)
    } as OkResponse<I, O>;
  } catch (error) {
    if (error.name == "AbortError") {
      error = new CancelError("HTTPリクエストがキャンセルされました");
    }
    if (isErrors(error)) {
      if (isOperationError(error)) {
        maintain(error);
      }
      return { args: endpoint.args, error } as ErrorResponse<I>;
    } else {
      throw error;
    }
  }
};

/**
 * `false`: wait response
 * {error: `Errors`}: error
 * {payload: `O`}: response (payloadInCancel in canceled && allowsCancel)
 */
export type FetchStatus<I, O> = false | ResponseType<I, O>;

/**
 * ```
 * const state = useHttpCommon(...);
 * ```
 */
export const useHttpCommon = <E, I, O>(
  endpoint: Endpoint<E, I, O>,
  task: RequestTask<E, I, O>
) => {
  const [, maintain] = useContext(OperationContext);
  const [state, setState] = useState<FetchStatus<I, O>>(false);

  useEffect(() => {
    const aborter = new AbortController();
    let setStateSafe = (c: FetchStatus<I, O>) => setState(c);
    (async () => {
      setStateSafe(await httpCommon(endpoint, task, aborter.signal, maintain));
    })();
    return () => {
      setStateSafe = () => null;
      aborter.abort();
      setState(false);
    };
  }, [endpoint, task, maintain]);
  return state;
};

export const useHttpCommonTask = <E, I, O>(
  endpoint: Endpoint<E, I, O>,
  task: RequestTask<E, I, O>,
  handlers: TaskHandler<I, O>
) =>
  useTask<I, O>(
    async (sig, args) => task({ ...endpoint, args }, sig),
    handlers
  );
