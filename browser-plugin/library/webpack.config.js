const path = require('path');

module.exports = {
    entry: './src/index.ts',
    // mode: 'production', // 'production' enables minification and other optimizations
    mode: 'development', // 'development' keeps the output readable
    devtool: 'inline-source-map', // include original sources for debugging
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
        ],
    },
    resolve: {
        extensions: ['.ts', '.js'],
    },
    output: {
        filename: 'news-deframer-lib.js',
        path: path.resolve(__dirname, 'dist'),
        library: {
            name: '_lib_test',
            type: 'window',
            export: 'default',
        },
    },
};