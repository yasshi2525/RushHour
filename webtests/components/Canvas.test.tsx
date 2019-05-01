import * as React from "react";
import { shallow, mount, ReactWrapper } from "enzyme";
import { Canvas } from "@/components/Canvas";
import { defaultState, RushHourStatus, GameMap } from "@/state";

const testMap: GameMap =  {
    "companies": [],
    "gates": [],
    "humans": [],
    "line_tasks": [],
    "platforms": [],
    "players": [],
    "rail_edges": [],
    "rail_nodes": [],
    "rail_lines": [],
    "residences": [{id: "1", x: 100, y: 100}],
    "stations": [],
    "trains": [],
};

test("renders canvas", () => {
    const wrapper = shallow(<Canvas {...defaultState} />);
    expect(wrapper.name()).toEqual("div");
});

describe("updates gamemodel", () => {
    let wrapper: ReactWrapper<RushHourStatus, RushHourStatus, Canvas>;
    beforeEach(() => {
        wrapper = mount(<Canvas {...defaultState} dispatch={()=>{}} />);
    });

    test("renders sprites in first time", () => {
        wrapper.setProps({ map: testMap });
    });

    test("skips rendering in second time", () => {
        wrapper.setProps({ map: testMap });
        wrapper.setProps({ map: testMap });
    });

    afterEach(() => {
        wrapper.unmount();
    });

});
