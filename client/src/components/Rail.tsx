import * as React from "react";
import { connect } from "react-redux";
import { Avatar } from "@material-ui/core";
import rail from "static/rail.png";
import { MenuStatus } from "state";
import ToggleButton from "./Toggle";

const RailCreator = ToggleButton(MenuStatus.SEEK_DEPARTURE)(() => (
  <Avatar alt="rail" src={rail} />
));

export default connect(null)(RailCreator);
