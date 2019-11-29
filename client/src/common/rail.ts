import { useContext, useCallback } from "react";
import { Entity } from "interfaces";
import { Depart, rail, Extend, Connect } from "interfaces/endpoint";
import { MenuStatus, ResolveError } from "interfaces/gamemap";
import { useTask } from "./task";
import { useHttpTask, http } from "./http";
import ModelContext from "./model";
import { usePlayers } from "./map";

export const useRail = () => {
  const model = useContext(ModelContext);
  const players = usePlayers();

  const resolve = useCallback((err: ResolveError) => {
    if (err.hasUnresolvedOwner) {
      players();
    }
  }, []);

  const [_depart] = useHttpTask<Depart, { rn: Entity }>(rail.depart, d => {
    const anchor = model.gamemap.mergeChild("rail_nodes", d.rn);
    if (!anchor) {
      console.warn("no anchor useRail.depart");
    } else {
      resolve(anchor.resolve({}));

      model.setMenuState(MenuStatus.EXTEND_RAIL);
      model.controllers.getAnchor().merge("anchor", {
        type: "rail_nodes",
        pos: anchor.get("pos"),
        cid: anchor.get("cid")
      });
    }
  });

  const [_extend] = useHttpTask<Extend, { rn: Entity; e1: Entity; e2: Entity }>(
    rail.extend,
    d => {
      const anchor = model.gamemap.mergeChild("rail_nodes", d.rn);
      const e1 = model.gamemap.mergeChild("rail_nodes", d.e1);
      const e2 = model.gamemap.mergeChild("rail_nodes", d.e2);
      if (!anchor || !e1 || !e2) {
        console.warn(
          `rn=${anchor}, e1=${e1}, e2=${e2} required Entity useRail.extend`
        );
      } else {
        resolve(anchor.resolve({}));
        resolve(e1.resolve({}));
        resolve(e2.resolve({}));
        model.controllers.getAnchor().merge("anchor", {
          type: "rail_nodes",
          pos: anchor.get("pos"),
          cid: anchor.get("cid")
        });
      }
    }
  );

  const [_connect] = useHttpTask<Connect, { e1: Entity; e2: Entity }>(
    rail.connect,
    d => {
      const anchor = model.gamemap.mergeChild("rail_nodes", d.e1.to);
      const e1 = model.gamemap.mergeChild("rail_nodes", d.e1);
      const e2 = model.gamemap.mergeChild("rail_nodes", d.e2);
      if (!anchor || !e1 || !e2) {
        console.warn(
          `rn=${anchor}, e1=${e1}, e2=${e2} required Entity useRail.connect`
        );
      } else {
        resolve(e1.resolve({}));
        resolve(e2.resolve({}));
        model.controllers.getAnchor().merge("anchor", {
          type: "rail_nodes",
          pos: anchor.get("pos"),
          cid: anchor.get("cid")
        });
      }
    }
  );

  const [_destroy] = useTask(
    (sig, args: { id: number; cid: number }) =>
      http(rail.destroy, sig, { id: args.id }).then(() => args.cid),
    id => {
      model.gamemap.removeChild("rail_nodes", id);
    }
  );
  const destroy = useCallback((id: number, cid: number) => {
    _destroy({ id, cid });
  }, []);

  return {
    depart: useCallback(
      (x: number, y: number) =>
        _depart({
          x,
          y,
          scale: Math.floor(model.coord.scale - model.delegate + 1)
        }),
      []
    ),
    extend: useCallback(
      (x: number, y: number, rnid: number) =>
        _extend({
          x,
          y,
          rnid,
          scale: Math.floor(model.coord.scale - model.delegate + 1)
        }),
      []
    ),
    connect: useCallback(
      (from: number, to: number) =>
        _connect({
          from,
          to,
          scale: Math.floor(model.coord.scale - model.delegate + 1)
        }),
      []
    ),
    destroy: useCallback((id: number, cid: number) => destroy(id, cid), [])
  };
};
