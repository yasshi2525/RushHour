const path = require("path");

module.exports = {
    watch: false,
    watchOptions: {
        ignored: ["./node_modules", "./app", "./public"]
    },
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
                use: ["style-loader", "css-loader?modules"],
            },
            {
                test: /\.tsx?$/,
                use: ["ts-loader"]
            },
            { 
                enforce: "pre", test: /\.js$/, 
                loader: "source-map-loader" 
            }
        ],
    }, 
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "web")
        },
        extensions: [".js", ".jsx", ".ts", ".tsx"]
    },
};