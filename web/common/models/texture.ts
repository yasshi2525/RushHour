export function generateTextures(app: PIXI.Application): {[index: string]: PIXI.Texture} {
    // from image
    ["residence", "company", "station", "train"].forEach(key => app.loader.add(key, `public/img/${key}.png`));
    app.loader.load();

    return {
        residence: app.loader.resources["residence"].texture,
        company: app.loader.resources["company"].texture,
    }
};



