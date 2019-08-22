import * as React from "react";
import { shallow } from "enzyme";
import GameContainer from "@/common";
import { GameBoard } from "@/components/GameBoard";
import { Canvas } from "@/components/Canvas";

const game = new GameContainer();
test("renders <Canvas />", () => {
    const wrapper = shallow(<GameBoard readOnly={true} game={game} dispatch={()=>{}} />);
    expect(wrapper.find(Canvas)).toBeDefined();
});
