import React, { FC, createContext, useMemo, useEffect } from "react";
import { ConfigResponse, game } from "interfaces/endpoint";
import { useHttpGet } from "./utils/http_get";
import DelayProvider from "./utils/provider_delay";
import { Errors } from "interfaces/error";

type ContextContents = [
  ConfigResponse,
  (v: ConfigResponse) => void,
  (e: Errors) => void
];

const ConfigContext = createContext<ContextContents>([
  { min_scale: -1, max_scale: -1 },
  () => console.warn("not initialized ConfigContext"),
  () => console.warn("not initialized ConfigContext")
]);
ConfigContext.displayName = "ConfigContext";

const converter = (d: ConfigResponse) => d;

export const ConfigProvider: FC = props => {
  const response = useHttpGet(game.const);

  useEffect(() => {
    console.info(`after ConfigProvider`);
  }, []);

  return useMemo(
    () => (
      <DelayProvider
        initialFetch={response}
        onError={
          <div>
            <p>ゲーム設定の読み込みに失敗しました</p>
            <p>画面を更新してください</p>
          </div>
        }
        convert={converter}
        ctx={ConfigContext}
      >
        {props.children}
      </DelayProvider>
    ),
    [response]
  );
};

export default ConfigContext;
