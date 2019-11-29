import React, { Suspense, lazy, useContext, useEffect } from "react";
import { Entity, GameMap } from "interfaces";
import {
  FetchMap as FetchMapType,
  players,
  fetchMap
} from "interfaces/endpoint";
import { useFetch } from "common/http";
import LoadingContext, { LoadingStatus, LoadingCircle } from "common/loading";
import ModelContext, { ModelProvider } from "common/model";

const Canvas = lazy(() => import("./Canvas"));

const TryCanvas = () => {
  return (
    <Suspense fallback={<LoadingCircle />}>
      <Canvas />
    </Suspense>
  );
};

const FetchMap = () => {
  const { update } = useContext(LoadingContext);
  const model = useContext(ModelContext);
  const [completed, err, data] = useFetch<FetchMapType, GameMap>(fetchMap, {
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
  const { update } = useContext(LoadingContext);
  const model = useContext(ModelContext);
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
  const { update } = useContext(LoadingContext);
  useEffect(() => {
    console.info(`useEffect GameBoard ${LoadingStatus.LOADED_RESOURCE}`);
    update(LoadingStatus.LOADED_RESOURCE);
    return () => {
      console.info("cleanup GameBoard");
    };
  }, []);
  return (
    <ModelProvider>
      <FetchPlayers />
    </ModelProvider>
  );
};
