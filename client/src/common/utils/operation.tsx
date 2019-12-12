import React, {
  createContext,
  useCallback,
  useState,
  useMemo,
  FC,
  ReactNode
} from "react";
import { game } from "interfaces/endpoint";
import { OperationError, Errors } from "interfaces/error";
import { useHttpGet } from "./http_get";
import DelayComponent from "./delay";

const initialReporter = (e: OperationError) => console.error(e);

type OperationStatus = [boolean, (e: Errors) => void];

const OperationContext = createContext<OperationStatus>([
  false,
  initialReporter
]);
OperationContext.displayName = "OperationContext";

interface OperatorProperties {
  maintenance: ReactNode;
}

/**
 * メンテナンス中か `/api/v1/game` にリクエストを送る。
 * メンテナンス中か判定数する。
 * `/api/v1/game` からのレスポンスが `200` で、 `status` キーが `true` なら稼働中と判定する
 * 結果に応じて、アプリケーション画面を構築するか、メンテナンス画面を構築するか決定する
 * 描画後、別のリクエストで `503` が帰ってきた場合は、メンテナンス画面に切り替えられるよう、
 * `setter` をコンテキスト化する
 */
export const OperationProvider: FC<OperatorProperties> = props => {
  const response = useHttpGet(game.status);
  const [isOperation, setOperation] = useState(true);
  const maintain = useCallback((e: OperationError) => {
    initialReporter(e);
    setOperation(false);
  }, []);
  return useMemo(
    () => (
      <DelayComponent initialFetch={response} onError={props.maintenance}>
        <OperationContext.Provider value={[isOperation, maintain]}>
          {isOperation ? props.children : props.maintenance}
        </OperationContext.Provider>
      </DelayComponent>
    ),
    [response]
  );
};

/**
 * [status, update] = useContext(OperationContext);
 */
export default OperationContext;
