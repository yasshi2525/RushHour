import React, {
  Context,
  ReactNode,
  useMemo,
  useState,
  useCallback,
  useContext
} from "react";
import { Errors, isOperationError } from "interfaces/error";
import OperationContext from "./operation";

type ContextContents<S> = [S, (data: S) => void, (error: Errors) => void];

export type SyncHandlers<S, C> = [
  Context<C>,
  S,
  (v: S) => void,
  Errors | undefined,
  (e: Errors) => void
];

export const useSyncStatus = <S, C>(
  context: Context<C>,
  initial: S
): SyncHandlers<S, C> => {
  const [, maintain] = useContext(OperationContext);
  const [status, setStatus] = useState<S>(initial);
  const [error, setError] = useState<Errors>();
  const report = useCallback(
    (e: Errors) => {
      if (isOperationError(e)) {
        maintain(e);
      }
      setError(e);
    },
    [maintain]
  );
  return [useMemo(() => context, [context]), status, setStatus, error, report];
};

interface SyncProperties<S, C> {
  children?: ReactNode;
  initial: S;
  onError: ReactNode;
  ctx: Context<C>;
}

function SyncProvider<S>(props: SyncProperties<S, ContextContents<S>>) {
  const [SyncContext, status, setStatus, error, report] = useSyncStatus(
    props.ctx,
    props.initial
  );

  return useMemo(() => {
    if (error) {
      return <>{props.onError}</>;
    } else {
      return (
        <SyncContext.Provider value={[status, setStatus, report]}>
          {props.children}
        </SyncContext.Provider>
      );
    }
  }, [SyncContext, status, setStatus, error, report]);
}

export default SyncProvider;
