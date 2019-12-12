import React, { FC, useContext, useEffect, useMemo, Fragment } from "react";
import Button from "@material-ui/core/Button";
import { PixiProvider } from "common/pixi";
import { PlayerProvider } from "common/player";
import LoadingCircle from "common/utils/loading";
import useLoader from "common/utils/image";
import CoordContext, { CoordProvider } from "common/coord";
import { ClockProvider } from "common/clock";
import { WindowProvider } from "common/windows";
import { DelegateProvider } from "common/delegate";
import { ConfigProvider } from "common/config";
import GameMapContext, { GameMapProvider } from "common/map";

// const Canvas = lazy(() => import("./Canvas"));

// const TryCanvas = () => {
//   return (
//     <Suspense fallback={<LoadingCircle />}>
//       <Canvas />
//     </Suspense>
//   );
// };

// const FetchMap = () => {
//   const [, update] = useContext(LoadingContext);
//   const model = useContext(ModelContext);
//   const [completed, err, data] = useFetch<FetchMapType, GameMap>(fetchMap, {
//     ...model.coord,
//     delegate: model.delegate
//   });

//   useEffect(() => {
//     if (completed) {
//       update(LoadingStatus.FETCHED_MAP);
//       if (data) {
//         model.gamemap.mergeAll(data);
//       }
//     }
//   }, [completed]);
//   if (!completed) {
//     return <LoadingCircle />;
//   } else if (err) {
//     return (
//       <>
//         <div>マップデータの読み込みに失敗しました。</div>
//         <div>画面を更新してください</div>
//         {err.messages.map(msg => (
//           <div>{msg}</div>
//         ))}
//       </>
//     );
//   } else {
//     return <TryCanvas />;
//   }
// };

const Test = () => {
  const [data] = useContext(GameMapContext);
  const [, , update] = useContext(CoordContext);
  useEffect(() => {
    console.info("after Test");
  }, []);
  return useMemo(
    () => (
      <Fragment>
        <Button
          onClick={() => update(Math.random(), Math.random(), Math.random())}
        >
          {data.timestamp}
        </Button>
      </Fragment>
    ),
    [data, update]
  );
};

const LoadResource: FC = props => {
  const [loaded, err] = useLoader();
  useEffect(() => {
    console.info("after LoadResource");
  }, []);
  return useMemo(() => {
    if (loaded) {
      if (err) {
        return (
          <>
            <p>画像データの読み込みに失敗しました</p>
            <p>画面を更新してください</p>
          </>
        );
      } else {
        return <>{props.children}</>;
      }
    } else {
      return <LoadingCircle />;
    }
  }, [loaded, err]);
};

export default () => {
  useEffect(() => {
    console.info("after GameBoard");
  }, []);
  return useMemo(
    () => (
      <ConfigProvider>
        <PixiProvider>
          <WindowProvider>
            <CoordProvider>
              <DelegateProvider>
                <ClockProvider>
                  <LoadResource>
                    <PlayerProvider>
                      <GameMapProvider>
                        <Test />
                      </GameMapProvider>
                    </PlayerProvider>
                  </LoadResource>
                </ClockProvider>
              </DelegateProvider>
            </CoordProvider>
          </WindowProvider>
        </PixiProvider>
      </ConfigProvider>
    ),
    []
  );
};
