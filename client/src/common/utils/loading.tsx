import React, { useMemo } from "react";

import CircularProgress from "@material-ui/core/CircularProgress";

const LoadingCircle = () =>
  useMemo(
    () => (
      <CircularProgress aria-describedby="loading-description" aria-busy={true}>
        <div id="loading-description">読み込み中</div>
      </CircularProgress>
    ),
    []
  );

export default LoadingCircle;
