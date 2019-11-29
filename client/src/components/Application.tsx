import React, { Suspense, useEffect, useState } from "react";
import { LoadingCircle } from "common/loading";

const GameBoard = React.lazy(() => import("./GameBoard"));
const AppBar = React.lazy(() => import("./AppBar"));

export default () => {
  // AppBar -> TryGameBoard の順でロードするため
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    console.info(`useEffect Application ${loaded}`);
    setLoaded(true);
    return () => {
      console.info("cleanup Application");
    };
  }, []);

  return (
    <Suspense fallback={<LoadingCircle />}>
      <AppBar />
      {loaded && (
        <Suspense fallback={<LoadingCircle />}>
          <GameBoard />
        </Suspense>
      )}
    </Suspense>
  );
};
