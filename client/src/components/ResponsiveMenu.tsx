import * as React from "react";
import { useDispatch } from "react-redux";
import { Hidden, Fab } from "@material-ui/core";
import ExpandIcon from "@material-ui/icons/Add";
import MinimizeIcon from "@material-ui/icons/Remove";
import GameModel from "common/models";
import { setMenu } from "actions";
import { MenuStatus } from "state";

interface ModelProperty {
  children: JSX.Element;
  model: GameModel;
}

export default (props: ModelProperty) => {
  const dispatch = useDispatch();

  const [expands, setExpand] = React.useState(false);

  const toggle = () => {
    let newState = !expands;
    if (!newState) {
      dispatch(setMenu.request({ model: props.model, menu: MenuStatus.IDLE }));
    }
    setExpand(newState);
  };

  return (
    <>
      {/* PC向け */}
      <Hidden xsDown>{props.children}</Hidden>
      {/* スマホ向け */}
      <Hidden smUp>
        {/* メニュー表示なし */}
        {expands ? (
          <Fab hidden={!expands} onClick={toggle}>
            <MinimizeIcon fontSize="large" />
          </Fab>
        ) : (
          <Fab color="primary" onClick={toggle}>
            <ExpandIcon fontSize="large" />
          </Fab>
        )}

        {/* メニュー表示あり */}
        {expands && props.children}
      </Hidden>
    </>
  );
};
