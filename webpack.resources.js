const path = require("path");

module.exports = {
    entry: {
        generate: [ "./resources/generate.ts" ]
    },
    output: {
        path: path.join(__dirname, "resources/js"),
        filename: "bundle.[name].js"
    },
    module: {
        rules: [
            { test: /\.ts$/, use: ["ts-loader"] }
        ]
    },
    resolve: {
        extensions: [".js", ".ts"]
    },
    mode: "development"
};