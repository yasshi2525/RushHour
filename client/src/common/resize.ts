import { useState, useEffect } from "react";

const getSize = () => [window.innerWidth, window.innerHeight];

export const useResize = () => {
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
