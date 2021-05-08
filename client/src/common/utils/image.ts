import { useContext, useState, useEffect } from "react";
import { Errors, ServerError } from "interfaces/error";
import PixiContext from "../pixi";

const SHEETS = [
  "cursor",
  "anchor",
  "residence",
  "company",
  "rail_node",
  "rail_edge",
  "destroy"
];

const load = (app: PIXI.Application) =>
  new Promise<void>((resolve, reject) => {
    SHEETS.forEach(key =>
      app.loader.add(
        key,
        `assets/bundle/spritesheet/${key}@${Math.floor(
          app.renderer.resolution
        )}x.json`
      )
    );
    app.loader.onError.add((e: Error) => reject(e));
    app.loader.load(() => resolve());
  });

const useLoader = () => {
  const [completed, setCompleted] = useState(false);
  const [err, setError] = useState<Errors>();
  const app = useContext(PixiContext);
  useEffect(() => {
    console.info(`effect useLoader`);
    (async () => {
      await load(app).catch(e => setError(new ServerError(e)));
      app.loader.resources;
      setCompleted(true);
    })();
    return () => {
      app.loader.reset();
    };
  }, [app]);
  return [completed, err] as [boolean, Errors | undefined];
};

export default useLoader;
