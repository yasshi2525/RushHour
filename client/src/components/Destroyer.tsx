import * as React from "react";
import { connect } from "react-redux";
import { Avatar } from "@material-ui/core";
import { MenuStatus } from "state";
import destroyer from "static/destroyer.png";
import ToggleButton from "./Toggle";

const Destroyer = ToggleButton(MenuStatus.DESTROY)(() => (
  <Avatar alt="destroy" src={destroyer} />
));

export default connect(null)(Destroyer);
