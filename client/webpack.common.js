const path = require("path");
const EnvironmentPlugin = require("webpack").EnvironmentPlugin;
const HtmlWebpackPlugin = require("html-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");

module.exports = {
  entry: {
    index: "./src/index.tsx"
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    publicPath: "/assets/bundle",
    filename: "[name].js"
  },
  module: {
    rules: [
      {
        test: /\.(html)$/,
        use: {
          loader: "html-loader",
          options: {
            attrs: [":data-src"]
          }
        }
      },
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader" /* css-loader?modules */]
      },
      {
        test: /\.tsx?$/,
        use: {
          loader: "ts-loader",
          options: {
            transpileOnly: true,
            experimentalWatchApi: true
          }
        },
        exclude: /node_modules/
      },
      {
        test: /\.(png|svg|jpg|gif)$/,
        use: "file-loader"
      },
      {
        test: /\.(woff|woff2|eot|ttf)$/,
        use: {
          loader: "file-loader",
          options: {
            name: "[name].[ext]"
            /* outputPath: "../fonts" ,
                            publicPath: "assets/bundle/fonts" */
          }
        }
      }
    ]
  },
  resolve: {
    modules: [path.resolve(__dirname, "src"), "node_modules"],
    extensions: [".ts", ".tsx", ".js", ",jsx"]
  },
  plugins: [
    new CleanWebpackPlugin(),
    new EnvironmentPlugin({ baseurl: "http://localhost:8080" }),
    new HtmlWebpackPlugin({})
  ]
};
