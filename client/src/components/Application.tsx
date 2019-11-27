import React, { Suspense, useEffect, useState } from "react";
import { LoadingStatus, useLoading } from "common/loading";
import { LoadingCircle } from "./Loading";

const GameBoard = React.lazy(() => import("./GameBoard"));
const AppBar = React.lazy(() => import("./AppBar"));

const TryGameBoard = () => {
  const [, update] = useLoading();
  useEffect(() => {
    console.info(
      `effect Application.GameBoard ${LoadingStatus.IMPORTED_BOARD}`
    );
    update(LoadingStatus.IMPORTED_BOARD);
  }, []);
  return (
    <Suspense fallback={<LoadingCircle />}>
      <GameBoard />
    </Suspense>
  );
};

export default () => {
  const [, update] = useLoading();
  // AppBar -> TryGameBoard の順でロードするため
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    console.info(`useEffect Application ${loaded}`);
    if (!loaded) {
      console.info(`effect AppBar ${LoadingStatus.IMPORTED_MENU}`);
      update(LoadingStatus.IMPORTED_MENU);
    }
    setLoaded(true);
    return () => {
      console.info("cleanup Application");
    };
  }, []);

  return (
    <>
      <Suspense fallback={<LoadingCircle />}>
        <AppBar />
        {loaded && <TryGameBoard />}
      </Suspense>
    </>
  );
};
