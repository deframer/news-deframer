const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    mode: 'development',
    entry: './src/debug/index.tsx',
    output: {
        path: path.resolve(__dirname, 'debug-dist'),
        filename: 'bundle.js',
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: {
                    loader: 'ts-loader',
                    options: {
                        transpileOnly: true,
                    },
                },
                exclude: /node_modules/,
            },
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader'],
            },
        ],
    },
    plugins: [
        new webpack.NormalModuleReplacementPlugin(
            /\.\/mode$/,
            (resource) => {
                if (resource.context.endsWith('shared')) {
                    resource.request = './mode-dev';
                }
            }
        ),
        new HtmlWebpackPlugin({
            template: './src/debug/index.html',
        }),
    ],
    devServer: {
        static: {
            directory: path.join(__dirname, 'debug-dist'),
        },
        port: 9000,
        hot: true,
    },
};