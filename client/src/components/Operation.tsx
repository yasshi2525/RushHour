import React, { Suspense, lazy, useEffect, useContext } from "react";
import OperationContext, { OperationProvider } from "common/operation";
import LoadingContext, { LoadingStatus, LoadingCircle } from "common/loading";

const Application = lazy(() => import("./Application"));
const Maintenance = lazy(() => import("./Maintenance"));

const InOperation = () => {
  const { update } = useContext(LoadingContext);
  useEffect(() => {
    console.info(`effect InOperation ${LoadingStatus.CHECKED_OPERATION}`);
    update(LoadingStatus.CHECKED_OPERATION);
    return () => {
      console.info("cleanup InOperation");
    };
  }, []);
  return (
    <Suspense fallback={<LoadingCircle />}>
      <Application />
    </Suspense>
  );
};

const UnderMaintenance = () => {
  const { update } = useContext(LoadingContext);
  useEffect(() => {
    console.info(`useEffect UnderMaintenance ${LoadingStatus.END}`);
    update(LoadingStatus.END);
    return () => {
      console.info("cleanup UnderMaintenance");
    };
  }, []);
  return (
    <Suspense fallback={<LoadingCircle />}>
      <Maintenance />
    </Suspense>
  );
};

const Operation = () => {
  const { status } = useContext(OperationContext);
  return status ? <InOperation /> : <UnderMaintenance />;
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
  const { update } = useContext(LoadingContext);
  useEffect(() => {
    console.info(`useEffect Operation ${LoadingStatus.CREATED_OPERATION}`);
    update(LoadingStatus.CREATED_OPERATION);
    return () => {
      console.info("cleanup Operation");
    };
  }, []);
  return (
    <OperationProvider>
      <Operation />
    </OperationProvider>
  );
};
