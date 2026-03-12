const path = require('path');
const webpack = require('webpack');
const CopyPlugin = require('copy-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const packageJson = require('./package.json');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = {
    mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
    entry: {
        options: './src/host/pages/index.tsx',
        popup: './src/host/popup/index.tsx',
        content: './src/host/content.ts',
        background: './src/host/background.ts'
    },
    devtool: process.env.NODE_ENV === 'production' ? 'source-map' : 'inline-source-map',
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: [{
                    loader: 'ts-loader',
                    options: {
                        configFile: 'tsconfig.build.json',
                        transpileOnly: true
                    }
                }],
                exclude: [/node_modules/, path.resolve(__dirname, 'src/debug')],
            },
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, 'css-loader'],
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
        modules: [path.resolve(__dirname, 'node_modules'), 'node_modules'],
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'dist/host'),
        clean: true,
    },
    watchOptions: {
        aggregateTimeout: 200,
        ignored: /node_modules/,
    },
    plugins: [
        process.env.NODE_ENV !== 'production' && new webpack.NormalModuleReplacementPlugin(
            /\.\/loglevel$/,
            (resource) => {
                if (resource.context.endsWith('shared')) {
                    resource.request = './loglevel-dev';
                }
            }
        ),
        new CopyPlugin({
            patterns: [
                {
                    from: 'src/host/manifest.json',
                    to: 'manifest.json',
                    transform(content, absoluteFrom) {
                        const manifest = JSON.parse(content.toString());
                        manifest.version = packageJson.version;
                        return JSON.stringify(manifest, null, 2);
                    }
                },
                { from: '../shared/assets/browser-extension', to: 'assets/browser-extension', noErrorOnMissing: true },
                { from: 'src/host/assets', to: 'assets', noErrorOnMissing: true, globOptions: { ignore: ['**/icons/**', '**/browser-extension/**'] } },
                { from: 'src/ndf/assets', to: 'assets', noErrorOnMissing: true },
            ],
        }),
        new HtmlWebpackPlugin({
            template: 'src/host/index.html',
            filename: 'options.html',
            chunks: ['options'],
        }),
        new HtmlWebpackPlugin({
            template: 'src/host/popup.html',
            filename: 'popup.html',
            chunks: ['popup'],
        }),
        new MiniCssExtractPlugin({
            filename: '[name].css',
        }),
    ].filter(Boolean),
    performance: {
        hints: false,
    },
};
