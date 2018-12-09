import * as React from "react";
import { shallow } from "enzyme";
import { GameBoard } from "@/components/GameBoard";
import { defaultState } from "@/state";
import { Canvas } from "@/components/Canvas";

test("renders <Canvas />", () => {
    const wrapper = shallow(<GameBoard {...defaultState} />);
    expect(wrapper.find(Canvas)).toBeDefined();
});
