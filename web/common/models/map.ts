import { Monitorable } from "../interfaces/monitor";
import { config } from "../interfaces/gamemap";
import { GameMap } from "../../state";
import GroupModel from "./group";
import { ResidenceContainer, CompanyContainer } from "./background";
import { StationContainer } from "./station";
import { RailEdge, RailNodeContainer, RailEdgeContainer, RailNode } from "./rail";

export default class extends GroupModel {

    init() {
        let textures = this.model.app.loader.resources;
        this.containers["residences"] = new ResidenceContainer({ model: this.model, app: this.model.app, texture: textures["residence"].texture});
        this.containers["companies"] = new CompanyContainer({ model: this.model, app: this.model.app, texture: textures["company"].texture});
        this.containers["stations"] = new StationContainer({ model: this.model, app: this.model.app, texture: textures["station"].texture});
        this.containers["rail_nodes"] = new RailNodeContainer({ model: this.model,app: this.model.app});
        this.containers["rail_edges"] = new RailEdgeContainer({ model: this.model, app: this.model.app});
    
        super.init();
    }

    mergeChild(type: string, props: {id: string}): undefined | Monitorable {
        if (this.containers[type] === undefined) {
            return undefined;
        } 
        return this.containers[type].mergeChild(props);
    }

    mergeAll(payload: GameMap) {
        config.zIndices.forEach(key => {
            if (this.containers[key] !== undefined) {
                this.containers[key].mergeChildren(payload[key], {coord: this.model.coord});
                if (this.containers[key].isChanged()) {
                    this.changed = true;
                }
            }
        });
        this.resolve();
        this.model.controllers.updateAnchor();
    }

    protected resolve() {
        if (this.containers["rail_nodes"] !== undefined) {
            this.containers["rail_nodes"].forEachChild((rn : RailNode) => {
                rn.resolve(this.get("rail_nodes", rn.get("pid")))
            });
        }
        if (this.containers["rail_edges"] !== undefined) {
            this.containers["rail_edges"].forEachChild((re: RailEdge) => 
                re.resolve(
                    this.get("rail_nodes", re.get("from")),
                    this.get("rail_nodes", re.get("to")),
                    this.get("rail_edges", re.get("eid"))
                )
            );
        }
    }
}