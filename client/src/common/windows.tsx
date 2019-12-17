import React, {
  FC,
  useContext,
  useMemo,
  useState,
  useEffect,
  createContext
} from "react";
import PixiContext from "./pixi";

const useWindow = () => {
  const app = useContext(PixiContext);
  const [width, setWidth] = useState(window.innerWidth);
  const [height, setHeight] = useState(window.innerHeight);

  useEffect(() => {
    console.info("effect useWindow");
    const onResize = () => {
      const width = window.innerWidth;
      const height = window.innerHeight;
      app.renderer.resize(width, height);
      setWidth(width);
      setHeight(height);
    };
    window.addEventListener("resize", onResize);
    return () => {
      window.removeEventListener("resize", onResize);
    };
  }, [app]);

  return [width, height];
};

const WindowContext = createContext([window.innerWidth, window.innerHeight]);
WindowContext.displayName = "WindowContext";

export const WindowProvider: FC = props => {
  const size = useWindow();
  useEffect(() => {
    console.info("after WindowProvider");
  }, []);
  return useMemo(
    () => (
      <WindowContext.Provider value={size}>
        {props.children}
      </WindowContext.Provider>
    ),
    [size]
  );
};

export default WindowContext;
