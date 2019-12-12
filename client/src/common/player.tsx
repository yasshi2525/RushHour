import React, { FC, useMemo, createContext, useContext } from "react";
import { HashContainer } from "interfaces";
import { players, Player, PlayersResponse } from "interfaces/endpoint";
import { Errors } from "interfaces/error";
import { useHttpGet, httpGet } from "./utils/http_get";
import AsyncProvider from "./utils/provider_async";
import OperationContext from "./utils/operation";
type Contents = [HashContainer<Player>, (v: {}) => void, (e: Errors) => void];
const PlayerContext = createContext<Contents>([
  {},
  () => console.warn("not initialized PlayerContext"),
  () => console.warn("not initialized PlayerContext")
]);
PlayerContext.displayName = "PlayerContext";

const converter = (payload: PlayersResponse) => payload.players;

export const PlayerProvider: FC = props => {
  const [, maintain] = useContext(OperationContext);
  const initialFetch = useHttpGet(players);
  return useMemo(
    () => (
      <AsyncProvider
        initialFetch={initialFetch}
        endpoint={players}
        reloadTask={(ep, sig) => httpGet(ep, sig)}
        convert={converter}
        onError={
          <div>
            <p>プレイヤー一覧の読み込みに失敗しました</p>
            <p>画面を更新してください</p>
          </div>
        }
        ctx={PlayerContext}
      >
        {props.children}
      </AsyncProvider>
    ),
    [initialFetch, maintain]
  );
};

/**
 * const [players, reload] = useContext(PlayerContext)
 */
export default PlayerContext;
