import React, {
  FC,
  useMemo,
  useState,
  useCallback,
  useEffect,
  useContext,
  createContext
} from "react";
import { useTicker } from "./utils/tick";
import PixiContext from "./pixi";

const ROUND = 240;

const useClock = () => {
  const pixi = useContext(PixiContext);
  const [offset, setOffset] = useState(0);

  const tick = useCallback(() => {
    setOffset(v => v + 1);
  }, []);

  useTicker(pixi, tick);

  useEffect(() => {
    if (offset >= ROUND) {
      setOffset(0);
    }
  }, [offset]);

  return offset;
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
