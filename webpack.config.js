const path = require('path');

module.exports = {
    watch: true,
    mode: 'development',
    entry: './web/index.js',
    output: {
        path: path.join(__dirname, 'public/js'),
        filename: 'bundle.js',
    },
    module: {
        rules: [
          {
            // 対象となるファイルの拡張子(cssのみ)
            test: /\.css$/,
            // Sassファイルの読み込みとコンパイル
            use: [
              // スタイルシートをJSからlinkタグに展開する機能
              'style-loader',
              // CSSをバンドルするための機能
              'css-loader',
            ],
          },
        ],
    },
};