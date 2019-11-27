export type GeneralObject = { [index: string]: NonNullable<any> };
export interface Entity extends GeneralObject {
  id: number;
}

export interface Locatable extends Entity {
  x: number;
  y: number;
}

export interface GameMap {
  companies: Locatable[];
  gates: Locatable[];
  humans: Locatable[];
  line_tasks: Locatable[];
  platforms: Locatable[];
  rail_edges: Locatable[];
  rail_lines: Locatable[];
  rail_nodes: Locatable[];
  residences: Locatable[];
  stations: Locatable[];
  trains: Locatable[];
}
