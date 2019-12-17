import React, {
  FC,
  useContext,
  useMemo,
  useEffect,
  createContext
} from "react";
import WindowContext from "./windows";
import CoordContext from "./coord";
import ConfigContext from "./config";

const W = 0;

const level = (width: number) => {
  if (width < 600) {
    // sm
    return 2;
  } else if (width < 960) {
    // md
    return 3;
  } else if (width < 1280) {
    // lg
    return 3;
  } else {
    // xl
    return 4;
  }
};

const useDelegate = () => {
  const [config] = useContext(ConfigContext);
  const window = useContext(WindowContext);
  const [, , scale] = useContext(CoordContext);
  const delegate = useMemo(
    () => Math.min(level(window[W]), Math.floor(scale) - config.min_scale),
    [config, window, scale]
  );
  return delegate;
};

const DelegateContext = createContext(level(window.innerWidth));
DelegateContext.displayName = "DelegateContext";

export const DelegateProvider: FC = props => {
  const delegate = useDelegate();
  useEffect(() => {
    console.info("after DelegateProvider");
  }, []);
  return useMemo(
    () => (
      <DelegateContext.Provider value={delegate}>
        {props.children}
      </DelegateContext.Provider>
    ),
    [delegate]
  );
};

export default DelegateContext;
