import * as React from "react";
import { connect } from "react-redux";
import { Container } from "@material-ui/core";
import ResponsiveMenu from "./ResponsiveMenu";
import RailCreator from "./Rail";
import Destroyer from "./Destroyer";

const ActionMenu = (props: any) => 
    <Container>
        <RailCreator {...props} />
        <Destroyer {...props} />
    </Container>
        

export default connect(null)(ResponsiveMenu()(ActionMenu));
