import {
  useState,
  useCallback,
  useEffect,
  useRef,
  useMemo,
  useContext
} from "react";
import {
  Errors,
  isErrors,
  CancelError,
  isCancelError,
  isOperationError
} from "interfaces/error";
import useSnack from "./snack";
import OperationContext from "./operation";

export type Task<I, O> = (signal: AbortSignal, args: I) => Promise<O>;

/**
 * [aborter, task, args]
 */
type TaskState<I, O> =
  | [undefined, undefined, undefined]
  | [AbortController, Task<I, O>, I];

export interface TaskHandler<I, O> {
  onOK?: (payload: O, args: I) => void;
  onError?: (e: Errors, args: I) => void;
  onCancel?: (e: CancelError, args: I) => void;
}

type Handlers<T> = [(args: T) => void, () => void];

/**
 * 非同期にタスクを実行する。タスクをキックするメソッドと、キャンセルするメソッドを返す
 * ```
 * const [fire, cancel] = useTask(
 *                          (sig, args) => http(sig, args), maintain,
 *                          { onOK:    (d) => console.info(`OK=${d}`),
 *                            onError: (e) => console.error(`NG=${e} on request(${args})`),
 *                            onCancel:(args, c) => console.warn(`CANCELED=${c} on request(${args})`)});
 * fire(args);
 * cancel(); // onCancel が指定されたとき cancel してもエラー扱いしない
 * ```
 */
const useTask = <I, O>(
  task: Task<I, O>,
  handlers?: TaskHandler<I, O>
): Handlers<I> => {
  const [, maintain] = useContext(OperationContext);
  const initState = useMemo<TaskState<I, O>>(
    () => [undefined, undefined, undefined],
    []
  );
  const [exec, setExecutor] = useState<TaskState<I, O>>(initState);
  const prevState = useRef<TaskState<I, O>>(exec);
  const snack = useSnack();

  /**
   * `fire()` でタスクを開始する
   */
  const fire = useCallback(
    (args: I) => {
      if (exec[0]) {
        exec[0].abort();
      }
      setExecutor([new AbortController(), task, args]);
    },
    [task, exec]
  );

  /**
   * `cancel()` でタスクを強制停止する
   */
  const cancel = useCallback(() => {
    if (exec[0]) {
      exec[0].abort();
      setExecutor(initState);
    }
  }, [exec]);

  /**
   * タスクが完了したら `onOK`、エラーなら `onError` をコールし、初期状態に戻す。
   * 実行中の場合、キャンセルする
   */
  useEffect(() => {
    if (prevState.current[0]) {
      console.warn("abort previsos task because new task is fired");
      prevState.current[0].abort();
    }
    (async () => {
      if (exec[0]) {
        try {
          const data = await exec[1](exec[0].signal, exec[2]);
          if (handlers?.onOK) {
            handlers.onOK(data, exec[2]);
          } else {
            // default message
            console.info("task ended");
          }
        } catch (e) {
          if (isErrors(e)) {
            if (isOperationError(e)) {
              maintain(e);
            } else if (isCancelError(e) && handlers?.onCancel) {
              handlers.onCancel(e, exec[2]);
            } else if (handlers?.onError) {
              handlers.onError(e, exec[2]);
            } else {
              // default message
              snack(e);
              console.warn("task ended error");
              console.warn(e);
            }
          } else {
            throw e;
          }
        }
        setExecutor(initState);
      }
    })();
  }, [prevState, exec, snack]);

  useEffect(() => {
    console.info(`update prevState running=${exec[0] !== undefined}`);
    prevState.current = exec;
  }, [exec]);

  return [fire, cancel];
};

export default useTask;
