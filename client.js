const { createCanvas } = require("canvas");
const gl = require("gl")(100,100, { stencil: true });

function createCanvasMock() {
    const canvas = createCanvas(100, 100);
    const __getContext = canvas.getContext.bind(canvas);

    // pixi.js Textureの読み込みで、indexOfが使われるため、モック化
    Object.defineProperty(canvas, "indexOf", { value: () => 0});

    // pixi.js Rendererの作成で、addEventListenerが使われるためモック化
    Object.defineProperty(canvas, "addEventListener", {value: ()=>{}});

    // pixi.js autoDensityの中で style属性が使われるためモック化
    Object.defineProperty(canvas, "style", { value: {}});

    // jsdomではwebglを使えないので、作成したwebglを返す
    Object.defineProperty(canvas, "getContext", {
        value: function(value, options) {
            if (value == "webgl") {
                return gl;
            } else {
                return __getContext(value, options);
            }
    }});
    return canvas;
}

const __createElement = document.createElement.bind(document);

// jsdomではcanvasを作れないので、作成したcanvasを返す
Object.defineProperty(document, "createElement", {
    value: function(value, options) {
        if (value == "canvas") {
            return createCanvasMock();
        } else {
            return __createElement(value, options);
        }
    }
});

// WebGLRenderingContext を定義しないと、pixi.jsがWebGL非サポートブラウザと判定してしまう
Object.defineProperty(window, "WebGLRenderingContext", { value: gl });
