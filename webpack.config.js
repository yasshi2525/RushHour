const path = require("path");

module.exports = {
    watch: true,
    mode: "development",
    entry: {
        index: ["./web/index.js", "./web/game.tsx"]
    },
    output: {
        path: path.join(__dirname, "public/js"),
        filename: "bundle.[name].js",
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: ["style-loader", "css-loader"],
            },      
            {
                test: /\.tsx?$/,
                use: ["awesome-typescript-loader"]
            },
            { 
                enforce: "pre", test: /\.js$/, 
                loader: "source-map-loader" 
            }
        ],
    }, 
    resolve: {
        extensions: [".js", ".jsx", ".ts", ".tsx"]
    },
};