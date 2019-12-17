import {
  useCallback,
  useState,
  useMemo,
  useEffect,
  useRef,
  useContext
} from "react";
import {
  MultiError,
  CancelError,
  isErrors,
  ErrorType,
  isCancelError,
  isOperationError
} from "interfaces/error";
import { OkResponse, ErrorResponse } from "./http_common";
import { Task } from "./task";
import useSnack from "./snack";
import OperationContext from "./operation";

type MultiTaskState<I, O> =
  | [undefined, undefined, undefined]
  | [AbortController, Task<I, O>, I[]];

type CancelResult<I> = { args: I; error: CancelError };
type ResultType<I, O> = OkResponse<I, O> | ErrorResponse<I> | CancelResult<I>;

export interface MultiTaskHandler<I, O> {
  onOK?: (payloadList: OkResponse<I, O>[]) => void;
  onError?: (payloadList: ErrorResponse<I>[]) => void;
  onCancel?: (payloadList: CancelResult<I>[]) => void;
}

const wrap = async <I, O>(
  task: Task<I, O>,
  signal: AbortSignal,
  args: I
): Promise<ResultType<I, O>> => {
  try {
    return { args, payload: await task(signal, args) };
  } catch (error) {
    if (isErrors(error)) {
      return { args, error };
    } else {
      throw error;
    }
  }
};

type Handlers<T> = [(argsList: T[]) => void, () => void];

/**
 * ```
 * const [fire, cancel] = useMultiTask((sig, args)=>http(sig, args));
 * fire(argsList);
 * ```
 */
const useMultiTask = <I, O>(
  task: Task<I, O>,
  handlers?: MultiTaskHandler<I, O>
): Handlers<I> => {
  const [, maintain] = useContext(OperationContext);
  const initState = useMemo<MultiTaskState<I, O>>(
    () => [undefined, undefined, undefined],
    []
  );
  const [exec, setExecutor] = useState<MultiTaskState<I, O>>(initState);
  const prevState = useRef<MultiTaskState<I, O>>(exec);
  const snack = useSnack();

  const isCanceled = <I>(obj: ResultType<I, O>): obj is CancelResult<I> =>
    "error" in obj && isCancelError(obj.error);
  const isErrored = <I>(obj: ResultType<I, O>): obj is ErrorResponse<I> =>
    "error" in obj && isErrors(obj.error);
  const isOk = <I, O>(obj: ResultType<I, O>): obj is OkResponse<I, O> =>
    "payload" in obj;

  const fire = useCallback(
    (argsList: I[]) => {
      if (exec[0]) {
        exec[0].abort();
      }
      setExecutor([new AbortController(), task, argsList]);
    },
    [task, exec]
  );

  const cancel = useCallback(() => {
    if (exec[0]) {
      exec[0].abort();
      setExecutor(initState);
    }
  }, [exec]);

  useEffect(() => {
    console.info("effect useTaskMulti");
    if (prevState.current[0]) {
      console.warn("abort previsos task because new task is fired");
      prevState.current[0].abort();
    }
    (async () => {
      if (exec[0]) {
        const payloadList = await Promise.all(
          exec[2].map(args => wrap(task, exec[0].signal, args))
        );
        if (handlers?.onOK) {
          handlers.onOK(payloadList.filter(isOk));
        } else {
          // default message
          console.info("multitask ended");
        }

        let errorList = payloadList.filter(isErrored);

        if (handlers?.onCancel) {
          const cancels = errorList.filter(isCanceled);
          if (cancels.length) {
            handlers.onCancel(cancels);
          }
          errorList.filter(e => e.error.type !== ErrorType.CANCEL);
        }

        if (errorList.length) {
          const opError = errorList
            .map(payload => payload.error)
            .find(isOperationError);
          if (opError) {
            maintain(opError);
          } else if (handlers?.onError) {
            handlers.onError(errorList);
          } else {
            // default message
            errorList.forEach(e => snack(e.error));
            console.warn("task ended error");
            console.warn(
              new MultiError(errorList.map(payload => payload.error))
            );
          }
        }
        setExecutor(initState);
      }
    })();
  }, [prevState, maintain, exec, snack]);

  useEffect(() => {
    console.info(`update prevState running=${exec[0] !== undefined}`);
    prevState.current = exec;
  }, [exec]);

  return [fire, cancel];
};

export default useMultiTask;
