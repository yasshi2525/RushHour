import { useState, useCallback, useEffect } from "react";
import { UnhandledError, ServerErrors, isServerErrors } from "interfaces/error";

/**
 * * `[]` タスク停止中
 * * `[aborter]` タスク実行中
 */
type TaskState<T> = [AbortController?, T?];

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
  task: (signal: AbortSignal, args?: I) => Promise<O>,
  onOK?: (data: O) => void,
  onError?: (e: ServerErrors) => void
) => {
  const [[aborter, args], setExecutor] = useState<TaskState<I>>([]);

  /**
   * `fire()` でタスクを開始する
   */
  const fire = useCallback(
    (args?: I) => {
      console.info("callback useTask fire");
      if (aborter) {
        throw new UnhandledError("cannot start running task");
      }
      setExecutor([new AbortController(), args]);
    },
    [aborter]
  );

  /**
   * `cancel()` でタスクを強制停止する
   */
  const cancel = useCallback(() => {
    console.info("callback useTask cancel");
    if (!aborter) {
      console.info("callback useTask cancel skip");
    } else {
      aborter.abort();
      setExecutor([]);
    }
  }, [aborter]);

  /**
   * タスクが完了したら `onOK`、エラーなら `onError` をコールし、初期状態に戻す。
   * 実行中の場合、キャンセルする
   */
  useEffect(() => {
    console.info(`effect useTask`);
    (async () => {
      if (aborter) {
        await task(aborter.signal, args)
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
      if (aborter) {
        aborter.abort();
      }
    };
  }, [aborter]);

  return [fire, cancel];
};
