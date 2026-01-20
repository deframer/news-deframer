const path = require('path');
const CopyPlugin = require('copy-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
    entry: {
        options: './src/host/ui/options.tsx',
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
                exclude: /node_modules/,
            },
            {
                test: /\.css$/i,
                use: ['style-loader', 'css-loader'],
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
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
        new CopyPlugin({
            patterns: [
                { from: 'src/host/manifest.json', to: 'manifest.json' },
                { from: 'src/host/assets', to: 'assets', noErrorOnMissing: true },
            ],
        }),
        new HtmlWebpackPlugin({
            template: 'src/host/ui/index.html',
            filename: 'options.html',
            chunks: ['options'],
        }),
    ],
};