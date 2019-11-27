import React, { Suspense, lazy, useState, useEffect } from "react";
import { ComponentProperty } from "interfaces/component";
import { game } from "interfaces/endpoint";
import OperationContext from "common/operation";
import { useFetch } from "common/http";
import { LoadingStatus, useLoading } from "common/loading";
import { LoadingCircle } from "./Loading";

const Application = lazy(() => import("./Application"));
const Maintenance = lazy(() => import("./Maintenance"));

const InOperation = () => {
  const [, update] = useLoading();
  useEffect(() => {
    console.info(`effect InOperation ${LoadingStatus.IMPORTED_APPLICATION}`);
    update(LoadingStatus.IMPORTED_APPLICATION);
  }, []);
  return (
    <Suspense fallback={<LoadingCircle />}>
      <Application />
    </Suspense>
  );
};

const UnderMaintenance = () => (
  <Suspense fallback={<LoadingCircle />}>
    <Maintenance />
  </Suspense>
);

interface OperationProperty extends ComponentProperty {
  status: boolean;
}

const Operation = (props: OperationProperty) => {
  const [inOperation, setOperation] = useState(props.status);
  return (
    <OperationContext.Provider value={[inOperation, setOperation]}>
      {inOperation ? <InOperation /> : <UnderMaintenance />}
    </OperationContext.Provider>
  );
};

/**
 * メンテナンス中か `/api/v1/game` にリクエストを送る。
 * メンテナンス中か判定数する。
 * `/api/v1/game` からのレスポンスが `200` で、 `status` キーが `true` なら稼働中と判定する
 * 結果に応じて、アプリケーション画面を構築するか、メンテナンス画面を構築するか決定する
 * 描画後、別のリクエストで `503` が帰ってきた場合は、メンテナンス画面に切り替えられるよう、
 * `setter` をコンテキスト化する
 */
export default () => {
  const [, update] = useLoading();
  const [completed, err, data] = useFetch<{}, { status: boolean }>(game.status);
  const [inOperation, setOperation] = useState<boolean | undefined>();

  useEffect(() => {
    console.info(`effect Operation ${completed} ${err} ${data}`);
    if (!completed) {
      console.info(`effect Operation WAIT`);
    } else {
      const ok = data?.status === true;
      if (ok) {
        console.info(`effect Operation OK`);
        update(LoadingStatus.CHECKED_OPERATION);
        setOperation(true);
      } else {
        console.info(`effect Operation NG`);
        update(LoadingStatus.END);
        setOperation(false);
      }
    }
    return () => {
      console.info("cleanup Operation");
    };
  }, [completed]);

  if (inOperation === undefined) {
    return <LoadingCircle />;
  } else {
    return <Operation status={inOperation === true} />;
  }
};
