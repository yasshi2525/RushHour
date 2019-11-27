import React, { useContext, useEffect } from "react";
import { Entity, GameMap } from "interfaces";
import { players, fetchMap } from "interfaces/endpoint";
import { useFetch } from "common/http";
import { useLoading, LoadingStatus } from "common/loading";
import LoginContext from "common/auth";
import ModelContext, { useModel } from "common/model";
import { LoadingCircle } from "./Loading";

//const Canvas = lazy(() => import("./Canvas"));

const TryCanvas = () => {
  const [loading, update] = useLoading();
  useEffect(() => update(LoadingStatus.LOADED_CANVAS), []);
  return (
    <></>
    // <Suspense fallback={<LoadingCircle />}>
    //   <Canvas />
    // </Suspense>
  );
};

interface FetchMapRequest {
  cx: number;
  cy: number;
  scale: number;
  delegate: number;
}

const FetchMap = () => {
  const model = useContext(ModelContext);
  const [, update] = useLoading();
  const [completed, err, data] = useFetch<FetchMapRequest, GameMap>(fetchMap, {
    ...model.coord,
    delegate: model.delegate
  });

  useEffect(() => {
    if (completed) {
      update(LoadingStatus.FETCHED_MAP);
      if (data) {
        model.gamemap.mergeAll(data);
      }
    }
  }, [completed]);
  if (!completed) {
    return <LoadingCircle />;
  } else if (err) {
    return (
      <>
        <div>マップデータの読み込みに失敗しました。</div>
        <div>画面を更新してください</div>
        {err.messages.map(msg => (
          <div>{msg}</div>
        ))}
      </>
    );
  } else {
    return <TryCanvas />;
  }
};

const FetchPlayers = () => {
  const model = useContext(ModelContext);
  const [, update] = useLoading();
  const [completed, err, data] = useFetch<{}, Entity[]>(players);
  useEffect(() => {
    if (completed) {
      update(LoadingStatus.FETCHED_PLAYERS);
      if (data) {
        model.gamemap.mergeChildren("players", data);
      }
    }
  }, [completed]);

  if (!completed) {
    return <LoadingCircle />;
  } else if (err) {
    return (
      <>
        <div>プレイヤー一覧の読み込みに失敗しました。</div>
        <div>画面を更新してください</div>
        {err.messages.map(msg => (
          <div>{msg}</div>
        ))}
      </>
    );
  } else {
    return <FetchMap />;
  }
};

export default () => {
  const [, update] = useLoading();
  const [[, my]] = useContext(LoginContext);
  const [completed, err, model] = useModel(my);
  useEffect(() => {
    console.info(`effect GameBoard ${completed}`);
    if (completed) {
      update(LoadingStatus.LOADED_RESOURCE);
    }
  }, [completed]);

  if (!completed) {
    return <LoadingCircle />;
  } else if (err) {
    return (
      <>
        <div>画像データの読み込みに失敗しました。</div>
        <div>画面を更新してください</div>
        {err?.messages.map(msg => (
          <div>{msg}</div>
        ))}
      </>
    );
  } else {
    return (
      <ModelContext.Provider value={model}>
        <FetchPlayers />
      </ModelContext.Provider>
    );
  }
};
