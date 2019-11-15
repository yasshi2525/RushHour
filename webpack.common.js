const path = require("path");
const webpack = require("webpack");

module.exports = {
    watch: false,
    watchOptions: {
        ignored: "^((?!web).)*$"
    },
    entry: {
        index: ["./web/index.tsx"]
    },
    output: {
        path: path.join(__dirname, "assets/bundle/js"),
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
                test: /\.(woff|woff2|eot|ttf|svg)$/,
                use: [
                    {
                        loader: "file-loader",
                        options: {
                            name: "[name].[ext]",
                            outputPath: "../fonts",
                            publicPath: "assets/bundle/fonts"
                        }
                    }
                ]
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
    plugins: [
        new webpack.EnvironmentPlugin({
            baseurl: "http://localhost:8080"
        })
    ]
};