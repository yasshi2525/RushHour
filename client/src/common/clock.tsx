import React, {
  FC,
  useMemo,
  useEffect,
  useContext,
  createContext
} from "react";
import { useCounter } from "./utils/tick";
import PixiContext from "./pixi";

const ROUND = 240;
const clock = () => {};

const useClock = () => {
  const pixi = useContext(PixiContext);
  return useCounter(pixi, clock, ROUND, true);
};

const ClockContext = createContext(0);
ClockContext.displayName = "ClockContext";

export const ClockProvider: FC = props => {
  const offset = useClock();
  useEffect(() => {
    console.info("after ClockProvider");
  }, []);
  return useMemo(
    () => (
      <ClockContext.Provider value={offset}>
        {props.children}
      </ClockContext.Provider>
    ),
    [offset]
  );
};

export default ClockContext;
