import * as React from "react";
import { connect } from "react-redux";
import { Avatar } from "@material-ui/core";
import ToggleButton from "./Toggle";
import { MenuStatus } from "state";

const Destroyer = ToggleButton(MenuStatus.DESTROY)(() => (
  <Avatar alt="destroy" src="/assets/img/destroyer.png" />
));

export default connect(null)(Destroyer);
