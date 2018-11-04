const path = require('path');

module.exports = {
    watch: true,
    mode: 'development',
    entry: {
        index: ['./web/index.js', './web/game.jsx']
    },
    output: {
        path: path.join(__dirname, 'public/js'),
        filename: 'bundle.[name].js',
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader'],
            },      
            {
                test: /\.jsx$/,
                exclude: /node_modules/,
                use: ['babel-loader']
            }
        ],
    }, 
    resolve: {
        extensions: ['.js', '.jsx']
    },
};