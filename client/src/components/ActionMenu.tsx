import * as React from "react";
import { Container } from "@material-ui/core";
import GameModel from "common/models";
import ResponsiveMenu from "./ResponsiveMenu";
import RailCreator from "./Rail";
import Destroyer from "./Destroyer";

interface ModelProperty {
  children?: JSX.Element;
  model: GameModel;
}

export default (props: ModelProperty) => (
  <ResponsiveMenu model={props.model}>
    <Container>
      <RailCreator {...props} />
      <Destroyer {...props} />
    </Container>
  </ResponsiveMenu>
);
