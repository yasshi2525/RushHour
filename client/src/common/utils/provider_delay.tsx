import React, { ReactNode, Context, useMemo } from "react";
import { Errors } from "interfaces/error";
import LoadingCircle from "./loading";
import { FetchStatus, isErrorResponse, OkResponse } from "./http_common";
import { SyncHandlers, useSyncStatus } from "./provider_sync";

type ContextContents<S> = [S, (v: S) => void, (e: Errors) => void];

export const useDelayContext = <I, O, S, C>(
  ctx: Context<C>,
  initial: OkResponse<I, O>,
  convert: (data: O, args: I) => S
): SyncHandlers<S, C> => {
  const [context, status, setStatus, error, report] = useSyncStatus<S, C>(
    ctx,
    convert(initial.payload, initial.args)
  );
  return [context, status, setStatus, error, report];
};

export interface DelayContextProperties<I, O, S, C> {
  children?: ReactNode;
  initial: OkResponse<I, O>;
  convert: (data: O, args: I) => S;
  onError: ReactNode;
  ctx: Context<C>;
}

const ContextProvider = <I, O, S>(
  props: DelayContextProperties<I, O, S, ContextContents<S>>
) => {
  const [AsyncStateContext, status, setStatus, error, report] = useDelayContext(
    props.ctx,
    props.initial,
    props.convert
  );

  return useMemo(() => {
    if (error) {
      return <>{props.onError}</>;
    } else {
      return (
        <AsyncStateContext.Provider value={[status, setStatus, report]}>
          {props.children}
        </AsyncStateContext.Provider>
      );
    }
  }, [AsyncStateContext, status, error, report]);
};

type DelayContextHandlers<I, O, S> = [
  Context<ContextContents<S>>,
  FetchStatus<I, O>,
  (data: O, args: I) => S
];

const useDelay = <I, O, S>(
  ctx: Context<ContextContents<S>>,
  initialFetch: FetchStatus<I, O>,
  convert: (data: O, args: I) => S
): DelayContextHandlers<I, O, S> => {
  return [
    useMemo(() => ctx, [ctx]),
    useMemo(() => initialFetch, [initialFetch]),
    useMemo(() => convert, [convert])
  ];
};

interface DelayProperties<I, O, S> {
  children?: ReactNode;
  initialFetch: FetchStatus<I, O>;
  convert: (data: O, args: I) => S;
  onError: ReactNode;
  ctx: Context<ContextContents<S>>;
}

/**
 * ```
 * const [state, report] = useContext(Context);
 * reload("logout");
 * reload("login");
 * report(e);
 * ```
 */
const DelayProvider = <I, O, S>(props: DelayProperties<I, O, S>) => {
  const [context, initialFetch, convert] = useDelay(
    props.ctx,
    props.initialFetch,
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
          initial={initialFetch}
          ctx={context}
          convert={convert}
          onError={props.onError}
        >
          {props.children}
        </ContextProvider>
      );
    }
  }, [initialFetch]);
};

export default DelayProvider;
