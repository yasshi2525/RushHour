import React, { ReactNode, useMemo } from "react";
import { FetchStatus, isErrorResponse } from "./http_common";
import LoadingCircle from "./loading";

type DelayHandlers<I, O> = [FetchStatus<I, O>];

const useDelay = <I, O>(
  initialFetch: FetchStatus<I, O>
): DelayHandlers<I, O> => {
  return [useMemo(() => initialFetch, [initialFetch])];
};

interface DelayProperties<I, O> {
  children?: ReactNode;
  initialFetch: FetchStatus<I, O>;
  onError: ReactNode;
}

const DelayComponent = <I, O>(props: DelayProperties<I, O>) => {
  const [response] = useDelay(props.initialFetch);

  return useMemo(() => {
    if (!response) {
      return <LoadingCircle />;
    } else if (isErrorResponse(response)) {
      return <>{props.onError}</>;
    } else {
      return <>{props.children}</>;
    }
  }, [response]);
};

export default DelayComponent;
