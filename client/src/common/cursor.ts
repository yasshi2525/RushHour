import React, { useState, useEffect, useCallback, useContext } from "react";
import { MenuStatus } from "interfaces/gamemap";
import ModelContext from "./model";
import { useReload } from "./map";
import { useRail } from "./rail";
import { useZoom } from "./zoom";

type Point = [number, number];
const X = 0;
const Y = 0;
const SENSITIVITY = 2;

export const useCursor = () => {
  const model = useContext(ModelContext);
  const [count, setCount] = useState(0);
  const [client, setClient] = useState<Point | undefined>();
  const [pressed, setPressed] = useState(false);
  const { depart, extend, connect, destroy } = useRail();
  const _zoom = useZoom();
  const reload = useReload();

  const zoom = useCallback((x: number, y: number) => {
    model.setCoord(..._zoom(x, y, model.coord.scale - 1));
    reload();
  }, []);

  const onMouseDown = useCallback((_: React.MouseEvent) => {
    setCount(0);
  }, []);

  const onMouseMove = useCallback(
    (e: React.MouseEvent) => {
      setClient([
        e.clientX * model.renderer.resolution,
        e.clientY * model.renderer.resolution
      ]);
      setCount(count + 1);
    },
    [count]
  );

  const onMouseOut = useCallback((_: React.MouseEvent) => {
    setClient(undefined);
    setCount(0);
  }, []);

  const onMouseUp = useCallback(
    (e: React.MouseEvent) => {
      if (count > SENSITIVITY) {
        console.warn("mouseend over move useCursor");
      } else {
        setClient([
          e.clientX * model.renderer.resolution,
          e.clientY * model.renderer.resolution
        ]);
        setPressed(true);
      }
    },
    [count]
  );

  const onTouchStart = useCallback((e: React.TouchEvent) => {
    if (e.targetTouches.length !== 1) {
      console.info("touchstart skip multi-touch useCursor");
    } else {
      setCount(0);
    }
  }, []);

  const onTouchMove = useCallback(
    (e: React.TouchEvent) => {
      if (e.targetTouches.length >= 2) {
        console.info("touchmove on multi-touch useCursor");
      } else if (count > SENSITIVITY) {
        console.info("touchmove overmove useTouch");
      } else {
        setCount(count + 1);
      }
    },
    [count]
  );

  const onTouchEnd = useCallback(
    (e: React.TouchEvent) => {
      if (count > SENSITIVITY) {
        console.info("touchmove ignore overmove useTouch");
        setCount(0);
      } else {
        setClient([
          e.changedTouches[0].clientX * model.renderer.resolution,
          e.changedTouches[0].clientY * model.renderer.resolution
        ]);
        setCount(0);
        setPressed(true);
      }
    },
    [count]
  );

  useEffect(() => {
    const pos = client
      ? {
          x: client[X],
          y: client[Y]
        }
      : undefined;
    model.controllers.getCursor().merge("client", pos);
  }, [client]);

  useEffect(() => {
    if (!client) {
      console.warn("no client useCursor");
    } else if (pressed) {
      const server = model.controllers.getCursor().get("pos");
      if (!server) {
        console.warn("no cursor.pos useCursor");
      } else {
        switch (model.menu) {
          case MenuStatus.SEEK_DEPARTURE:
            if (!model.controllers.getCursor().selected) {
              depart(server.x, server.y);
            } else if (model.controllers.getCursor().get("mul") === 1) {
              model.setMenuState(MenuStatus.EXTEND_RAIL);
              model.controllers
                .getAnchor()
                .merge(
                  "anchor",
                  model.controllers.getCursor().genAnchorStatus()
                );
            } else {
              zoom(...client);
            }
            break;
          case MenuStatus.EXTEND_RAIL:
            const anchor = model.controllers.getAnchor().object;
            if (!anchor) {
              console.warn("no anchor useCursor");
            } else if (!model.controllers.getCursor().get("activation")) {
              console.warn("inactive anchor useCursor");
            } else {
              const cursor = model.controllers.getCursor().selected;
              if (!cursor) {
                extend(server.x, server.y, anchor.get("cid"));
              } else {
                if (model.controllers.getCursor().get("mul") === 1) {
                  connect(anchor.get("cid"), cursor.get("cid"));
                } else {
                  zoom(...client);
                }
              }
            }
            break;
          case MenuStatus.DESTROY:
            const destroyer = model.controllers.getCursor().destroyer.selected;
            if (!destroyer) {
              console.warn("no destroyer useCursor");
            } else if (destroyer.get("mul") === 1) {
              destroy(destroyer.get("id"), destroyer.get("cid"));
            } else {
              zoom(...client);
            }
        }
      }
      setPressed(false);
    }
  }, [pressed]);

  return {
    onMouseDown,
    onMouseMove,
    onMouseOut,
    onMouseUp,
    onTouchStart,
    onTouchMove,
    onTouchEnd
  };
};
