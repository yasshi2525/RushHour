import * as React from "react";
import { shallow } from "enzyme";
import GameContainer from "@/common";
import { GameBoard } from "@/components/GameBoard";
import { Canvas } from "@/components/Canvas";

const game = new GameContainer(0);
test("renders <Canvas />", () => {
    const wrapper = shallow(<GameBoard readOnly={true} displayName={undefined} image={undefined} game={game} dispatch={()=>{}} isPIXILoaded={false} isPlayersFetched={false} />);
    expect(wrapper.find(Canvas)).toBeDefined();
});
