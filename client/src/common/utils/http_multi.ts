import { useEffect, useState, useContext } from "react";
import { OperationError, isErrors } from "interfaces/error";
import {
  Endpoint,
  RequestTask,
  OkResponse,
  ErrorResponse,
  httpCommon
} from "./http_common";
import useMultiTask, { MultiTaskHandler } from "./task_multi";
import OperationContext from "./operation";

export type MultiEndpoint<E, I, O> = E & { argsList: I[]; payload: O };

type MultiResponse<I, O> = {
  errorList: ErrorResponse<I>[];
  payloadList: OkResponse<I, O>[];
};

export type MultiFetchStatus<I, O> = false | MultiResponse<I, O>;

const multiHttpCommon = async <E, I, O>(
  endpointList: MultiEndpoint<E, I, O>,
  task: RequestTask<E, I, O>,
  signal: AbortSignal,
  maintain: (e: OperationError) => void
) => {
  const payloadList = await Promise.all(
    endpointList.argsList.map(args =>
      httpCommon({ ...endpointList, args }, task, signal, maintain)
    )
  );

  const isErrored = <I, O>(
    obj: OkResponse<I, O> | ErrorResponse<I>
  ): obj is ErrorResponse<I> => "error" in obj && isErrors(obj.error);
  const isOk = <I, O>(
    obj: OkResponse<I, O> | ErrorResponse<I>
  ): obj is OkResponse<I, O> => "payload" in obj;

  return {
    errorList: payloadList.filter(isErrored).map(e => e),
    payloadList: payloadList.filter(isOk).map(e => e)
  };
};

export const useMultiHttpCommon = <E, I, O>(
  endpointList: MultiEndpoint<E, I, O>,
  task: RequestTask<E, I, O>
) => {
  const [, maintain] = useContext(OperationContext);
  const [state, setState] = useState<MultiFetchStatus<I, O>>(false);

  useEffect(() => {
    console.info("effect useHttpMulti");
    const aborter = new AbortController();
    let setStateSafe = (c: MultiFetchStatus<I, O>) => setState(c);
    (async () => {
      setStateSafe(
        await multiHttpCommon(endpointList, task, aborter.signal, maintain)
      );
    })();
    return () => {
      setStateSafe = () => null;
      aborter.abort();
      setState(false);
    };
  }, [endpointList, task, maintain]);

  return state;
};

export const useMultiHttpCommonTask = <E, I, O>(
  ep: Endpoint<E, I, O>,
  task: RequestTask<E, I, O>,
  handlers: MultiTaskHandler<I, O>
) =>
  useMultiTask<I, O>(async (sig, args) => task({ ...ep, args }, sig), handlers);
