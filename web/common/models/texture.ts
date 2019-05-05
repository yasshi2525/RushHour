const opts = {
    rail_node: {
        width: 2,
        color: 0x4169e1,
        radius: 10
    }
}

function generateRailNode(app: PIXI.Application): PIXI.Texture {
    let graphics = new PIXI.Graphics();
    graphics.lineStyle(opts.rail_node.width, opts.rail_node.color);
    graphics.arc(0, 0, opts.rail_node.radius, 0, Math.PI * 2);
    return app.renderer.generateTexture(graphics, PIXI.SCALE_MODES.LINEAR, app.renderer.resolution);
}

export function generateTextures(app: PIXI.Application): {[index: string]: PIXI.Texture} {
    // from image
    ["residence", "company", "station", "train"].forEach(key => app.loader.add(key, `public/img/${key}.png`));
    app.loader.load();

    return {
        residence: app.loader.resources["residence"].texture,
        company: app.loader.resources["company"].texture,
        rail_node: generateRailNode(app)
    }
};



