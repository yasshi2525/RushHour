import { useState, useEffect, useContext } from "react";
import ModelContext from "./model";
import { useReload } from "./map";

const getSize = () => [window.innerWidth, window.innerHeight];

const _useResize = () => {
  const [windowsSize, setSize] = useState(getSize());
  useEffect(() => {
    const onResize = () => {
      setSize(getSize());
    };
    window.addEventListener("resize", onResize);
    return () => {
      window.removeEventListener("resize", onResize);
    };
  });

  return windowsSize;
};

export const useResize = () => {
  const model = useContext(ModelContext);
  const reloader = useReload();
  const [w, h] = _useResize();

  useEffect(() => {
    if (model.resize(w, h)) {
      reloader();
    }
  }, [w, h]);
};
