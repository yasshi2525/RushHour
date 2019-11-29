import { useState, useCallback, useEffect } from "react";
import { UnhandledError, ServerErrors, isServerErrors } from "interfaces/error";
import { useSnack } from "./snack";

/**
 * * `[]` タスク停止中
 * * `[aborter]` タスク実行中
 */
type TaskState<T> = [AbortController, T] | [];

/**
 * 非同期にタスクを実行する。タスクをキックするメソッドと、キャンセルするメソッドを返す
 * ```
 * const [fire, cancel] = useTask(
 *                          (sig, args) => http(sig, args),
 *                          (d) => console.log(`OK=${d}`),
 *                          (e) => console.error(`NG=${e}`));
 * fire(args);
 * cancel();
 * ```
 */
export const useTask = <I, O>(
  task: (signal: AbortSignal, args: I) => Promise<O>,
  onOK?: (data: O) => void,
  onError?: (e: ServerErrors) => void
) => {
  const [exec, setExecutor] = useState<TaskState<I>>([]);
  const snack = useSnack();

  /**
   * `fire()` でタスクを開始する
   */
  const fire = useCallback(
    (args: I) => {
      console.info("callback useTask fire");
      if (exec) {
        throw new UnhandledError("cannot start running task");
      }
      setExecutor([new AbortController(), args]);
    },
    [exec]
  );

  /**
   * `cancel()` でタスクを強制停止する
   */
  const cancel = useCallback(() => {
    console.info("callback useTask cancel");
    if (!exec.length) {
      console.info("callback useTask cancel skip");
    } else {
      exec[0].abort();
      setExecutor([]);
    }
  }, [exec]);

  /**
   * タスクが完了したら `onOK`、エラーなら `onError` をコールし、初期状態に戻す。
   * 実行中の場合、キャンセルする
   */
  useEffect(() => {
    console.info(`effect useTask`);
    (async () => {
      if (exec.length) {
        await task(exec[0].signal, exec[1])
          .then(data => {
            if (onOK) {
              onOK(data);
            } else {
              console.info("effect useTask ok");
            }
          })
          .catch(e => {
            if (isServerErrors(e)) {
              if (onError) {
                onError(e);
              } else {
                snack(e);
                console.warn("effect useTask error");
                console.warn(e);
              }
            } else {
              throw e;
            }
          });
        setExecutor([]);
      } else {
        console.info(`effect useTask noop`);
      }
    })();

    return () => {
      console.info("cleanup useTask");
      if (exec.length) {
        exec[0].abort();
      }
    };
  }, [exec]);

  return [fire, cancel];
};
