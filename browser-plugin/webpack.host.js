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
  devtool: 'inline-source-map',
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: [{
          loader: 'ts-loader',
          options: {
            configFile: 'tsconfig.build.json'
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
  plugins: [
    new CopyPlugin({
      patterns: [
        { from: 'src/host/manifest.json', to: 'manifest.json' },
        { from: 'src/host/assets', to: 'assets', noErrorOnMissing: true },
        // Copy the built library into the host assets for Release mode
        // We assume build:lib has run before build:host
        { 
          from: 'dist/library/library.bundle.js', 
          to: 'assets/library.bundle.js', 
          noErrorOnMissing: true 
        },
      ],
    }),
    new HtmlWebpackPlugin({
      template: 'src/host/ui/index.html',
      filename: 'options.html',
      chunks: ['options'],
    }),
  ],
};
