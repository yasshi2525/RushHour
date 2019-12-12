import React, {
  FC,
  useContext,
  useMemo,
  useState,
  useEffect,
  createContext
} from "react";
import PixiContext from "./pixi";

type Dimension = [number, number];
const getSize = (): Dimension => [window.innerWidth, window.innerHeight];

const useWindow = () => {
  const app = useContext(PixiContext);
  const [size, setSize] = useState(getSize());

  useEffect(() => {
    const onResize = () => {
      const now = getSize();
      app.renderer.resize(...now);
      setSize(now);
    };
    window.addEventListener("resize", onResize);
    return () => {
      window.removeEventListener("resize", onResize);
    };
  }, []);

  return size;
};

const WindowContext = createContext(getSize());
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
