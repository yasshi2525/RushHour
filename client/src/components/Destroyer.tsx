import * as React from "react";
import { Avatar } from "@material-ui/core";
import { MenuStatus } from "state";
import GameModel from "common/models";
import destroyer from "static/destroyer.png";
import ToggleButton from "./Toggle";

interface ModelProperty {
  model: GameModel;
}
export default (props: ModelProperty) => (
  <ToggleButton {...props} on={MenuStatus.DESTROY}>
    <Avatar alt="destroy" src={destroyer} />
  </ToggleButton>
);
