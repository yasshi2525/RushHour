import * as React from "react";
import { Avatar } from "@material-ui/core";
import GameModel from "common/models";
import rail from "static/rail.png";
import { MenuStatus } from "state";
import ToggleButton from "./Toggle";

interface ModelProperty {
  model: GameModel;
}

export default (props: ModelProperty) => (
  <ToggleButton {...props} on={MenuStatus.SEEK_DEPARTURE}>
    <Avatar alt="rail" src={rail} />
  </ToggleButton>
);
