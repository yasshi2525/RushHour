const path = require("path");
const merge = require("webpack-merge");
const common = require("./webpack.common.js");
const BundleAnalyzerPlugin = require("webpack-bundle-analyzer")
  .BundleAnalyzerPlugin;

const plugins = [];
if (process.env.analyze) {
  plugins.push(new BundleAnalyzerPlugin());
}

module.exports = merge(common, {
  mode: "development",
  devtool: "inline-source-map",
  devServer: {
    contentBase: path.resolve(__dirname, "dist"),
    inline: true,
    historyApiFallback: true,
    hot: true,
    port: 8090,
    proxy: {
      "/api": "http://localhost:8080"
    },
    before(app) {
      app.post("/", (req, res) => {
        res.redirect("/");
      });
    }
  },
  plugins
});
