const path = require("path");
const merge = require("webpack-merge");
const common = require("./webpack.common.js");

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
  }
});
