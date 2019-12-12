import React, { ReactNode, Context, useMemo } from "react";
import { Errors } from "interfaces/error";
import LoadingCircle from "./loading";
import {
  FetchStatus,
  isErrorResponse,
  OkResponse,
  RequestTask,
  Endpoint
} from "./http_common";
import useTask from "./task";
import { DelayContextProperties, useDelayContext } from "./provider_delay";

type Contents<I, S> = [S, (args: I) => void, (e: Errors) => void];

type Handlers<I, S, C> = [
  Context<C>,
  S,
  (args: I) => void,
  Errors | undefined,
  (e: Errors) => void
];

const useAsyncContext = <E, I, O, S, C>(
  ctx: Context<C>,
  initial: OkResponse<I, O>,
  endpoint: Endpoint<E, I, O>,
  reloadTask: RequestTask<E, I, O>,
  convert: (data: O, args: I) => S
): Handlers<I, S, C> => {
  const [context, status, setStatus, error, report] = useDelayContext<
    I,
    O,
    S,
    C
  >(ctx, initial, convert);

  const handlers = useMemo(
    () => ({
      onOK: (data: O, args: I) => setStatus(convert(data, args)),
      onError: (err: Errors) => report(err)
    }),
    [report, convert]
  );
  const [fire] = useTask<I, O>(
    (signal, args) => reloadTask({ ...endpoint, args }, signal),
    handlers
  );

  return [useMemo(() => context, [context]), status, fire, error, report];
};

interface ContextProperties<E, I, O, S, C>
  extends DelayContextProperties<I, O, S, C> {
  endpoint: Endpoint<E, I, O>;
  reloadTask: RequestTask<E, I, O>;
}

const ContextProvider = <E, I, O, S>(
  props: ContextProperties<E, I, O, S, Contents<I, S>>
) => {
  const [AsyncStateContext, status, fire, error, report] = useAsyncContext(
    props.ctx,
    props.initial,
    props.endpoint,
    props.reloadTask,
    props.convert
  );

  return useMemo(() => {
    if (error) {
      return <>{props.onError}</>;
    } else {
      return (
        <AsyncStateContext.Provider value={[status, fire, report]}>
          {props.children}
        </AsyncStateContext.Provider>
      );
    }
  }, [AsyncStateContext, status, fire, error, report]);
};

type AsyncHandlers<E, I, O, S> = [
  Context<Contents<I, S>>,
  FetchStatus<I, O>,
  E,
  RequestTask<E, I, O>,
  (data: O, args: I) => S
];

const useAsyncStatus = <E, I, O, S>(
  ctx: Context<Contents<I, S>>,
  initialFetch: FetchStatus<I, O>,
  endpoint: E,
  reloadTask: RequestTask<E, I, O>,
  convert: (data: O, args: I) => S
): AsyncHandlers<E, I, O, S> => {
  return [
    useMemo(() => ctx, [ctx]),
    useMemo(() => initialFetch, [initialFetch]),
    useMemo(() => endpoint, [endpoint]),
    useMemo(() => reloadTask, [reloadTask]),
    useMemo(() => convert, [convert])
  ];
};

interface AsyncProperties<E, I, O, S> {
  children?: ReactNode;
  initialFetch: FetchStatus<I, O>;
  endpoint: Endpoint<E, I, O>;
  reloadTask: RequestTask<E, I, O>;
  convert: (data: O, args: I) => S;
  onError: ReactNode;
  ctx: Context<Contents<I, S>>;
}

/**
 * ```
 * const [state, reload, report] = useContext(Context);
 * reload("logout");
 * reload("login");
 * report(e);
 * ```
 */
const AsyncProvider = <E, I, O, S>(props: AsyncProperties<E, I, O, S>) => {
  const [context, initialFetch, endpoint, reloadTask, convert] = useAsyncStatus(
    props.ctx,
    props.initialFetch,
    props.endpoint,
    props.reloadTask,
    props.convert
  );

  return useMemo(() => {
    if (!initialFetch) {
      return <LoadingCircle />;
    } else if (isErrorResponse(initialFetch)) {
      return <>{props.onError}</>;
    } else {
      return (
        <ContextProvider
          ctx={context}
          initial={initialFetch}
          endpoint={endpoint}
          reloadTask={reloadTask}
          convert={convert}
          onError={props.onError}
        >
          {props.children}
        </ContextProvider>
      );
    }
  }, [initialFetch]);
};

export default AsyncProvider;
