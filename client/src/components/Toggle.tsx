import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Fab } from "@material-ui/core";
import GameModel from "common/models";
import { setMenu } from "actions";
import { MenuStatus, RushHourStatus } from "state";

export interface ToggleProperty {
  model: GameModel;
  children: JSX.Element;
  on?: MenuStatus;
  off?: MenuStatus;
}

export default (props: ToggleProperty) => {
  const onStatus: MenuStatus =
    props.on !== undefined ? props.on : MenuStatus.IDLE;
  const offStatus: MenuStatus =
    props.off !== undefined ? props.off : MenuStatus.IDLE;

  const menu = useSelector<RushHourStatus, MenuStatus>(state => state.menu);
  const [selected, setSelected] = React.useState(menu == onStatus);
  const dispatch = useDispatch();
  const toggle = () => {
    let newState = !selected;
    if (newState) {
      dispatch(setMenu.request({ model: props.model, menu: onStatus }));
    } else {
      dispatch(setMenu.request({ model: props.model, menu: offStatus }));
    }
    setSelected(newState);
  };

  return (
    <>
      {selected ? (
        <Fab color="primary" onClick={toggle}>
          {props.children}
        </Fab>
      ) : (
        <Fab onClick={toggle}>{props.children}</Fab>
      )}
    </>
  );
};
