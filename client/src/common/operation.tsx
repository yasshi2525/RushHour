import React, { Dispatch, createContext, useState, useMemo } from "react";
import { game } from "interfaces/endpoint";
import { ServerErrors } from "interfaces/error";
import { ComponentProperty } from "interfaces/component";
import { LoadingCircle } from "./loading";
import { useFetch } from "./http";

type OperationStatus = [
  boolean,
  ServerErrors | undefined,
  boolean,
  Dispatch<boolean>
];

const useOperation = (): OperationStatus => {
  const [completed, err, data] = useFetch<{}, { status: boolean }>(game.status);
  const [, setOperation] = useState<boolean>(false);

  return useMemo(() => {
    if (!completed) {
      console.info(`effect useOperation WAIT`);
      return [false, undefined, false, () => {}];
    } else {
      const ok = data?.status === true;
      if (ok) {
        console.info(`effect useOperation OK`);
        setOperation(true);
      } else {
        console.info(`effect useOperation NG`);
        setOperation(false);
      }
      return [true, err, ok, setOperation];
    }
  }, [completed]);
};

interface ContextStatus {
  status: boolean;
  update: Dispatch<boolean>;
}

const OperationContext = createContext<ContextStatus>({
  status: false,
  update: () => {}
});

export const OperationProvider = (props: ComponentProperty) => {
  const [completed, , status, update] = useOperation();

  if (!completed) {
    return <LoadingCircle />;
  } else {
    return (
      <OperationContext.Provider value={{ status, update }}>
        {props.children}
      </OperationContext.Provider>
    );
  }
};

export default OperationContext;
