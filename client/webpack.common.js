const path = require("path");
const EnvironmentPlugin = require("webpack").EnvironmentPlugin;
const HtmlWebpackPlugin = require("html-webpack-plugin");
const HtmlWebpackHarddiskPlugin = require("html-webpack-harddisk-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
  entry: {
    index: "./src/index.tsx"
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    publicPath: "/assets/bundle/",
    filename: "[name].js"
  },
  module: {
    rules: [
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
          }
        }
      },
      {
        test: /\.html$/,
        loader: "html-loader"
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
    new HtmlWebpackPlugin({
      title: "RushHour",
      favicon: "src/favicon.ico",
      template: path.resolve(__dirname, "src", "index.ejs"),
      meta: {
        viewport: "width=device-width,initial-scale=1",
        "Content-Type": "text/html; charset=utf-8"
      },
      alwaysWriteToDisk: true
    }),
    new HtmlWebpackHarddiskPlugin(),
    new CopyPlugin([{ from: "src/static/import", to: "spritesheet" }])
  ]
};
