import * as React from "react";
import { connect } from "react-redux";
import { Avatar } from "@material-ui/core";
import { MenuStatus } from "../state";
import ToggleButton from "./Toggle";

const RailCreator = ToggleButton(MenuStatus.SEEK_DEPARTURE)(() => (
  <Avatar alt="rail" src="/assets/img/rail.png" />
));

export default connect(null)(RailCreator);
