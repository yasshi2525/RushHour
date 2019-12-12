import { useContext, useState, useMemo } from "react";
import ModelContext from "./model";

type Point = [number, number];
const ZERO: Point = [0, 0];

const useClitoSrv = () => {
  const model = useContext(ModelContext);
  const [client, setClient] = useState(ZERO);
  const server = useMemo<Point>(() => {
    const w = model.renderer.width;
    const h = model.renderer.height;
    const size = Math.max(model.renderer.width, model.renderer.height);
    return [(client[0] - w / 2) / size, (client[1] - h / 2) / size];
  }, [client]);
};

export const useCoord = () => {};
