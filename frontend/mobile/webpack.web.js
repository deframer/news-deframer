const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');

module.exports = {
  mode: 'development',
  entry: path.resolve(__dirname, 'web/index.js'),
  output: {
    path: path.resolve(__dirname, 'web/dist'),
    filename: 'bundle.js',
    publicPath: '/',
    clean: true,
  },
  resolve: {
    alias: {
      'react-native$': 'react-native-web',
    },
    extensions: ['.web.js', '.web.ts', '.web.tsx', '.js', '.ts', '.tsx', '.json'],
    modules: [path.resolve(__dirname, 'node_modules'), 'node_modules'],
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx|ts|tsx)$/,
        include: [
          path.resolve(__dirname, 'App.tsx'),
          path.resolve(__dirname, 'index.js'),
          path.resolve(__dirname, 'src'),
          path.resolve(__dirname, 'web'),
          path.resolve(__dirname, 'node_modules/react-native'),
          path.resolve(__dirname, 'node_modules/react-native-safe-area-context'),
          path.resolve(__dirname, 'node_modules/@react-native-async-storage/async-storage'),
          path.resolve(__dirname, 'node_modules/lucide-react-native'),
          path.resolve(__dirname, 'node_modules/react-native-svg'),
        ],
        use: {
          loader: 'babel-loader',
        },
      },
      {
        test: /\.svg$/,
        include: [
          path.resolve(__dirname, 'src'),
          path.resolve(__dirname, '../shared/assets'),
        ],
        use: [
          {
            loader: '@svgr/webpack',
            options: {
              native: true,
            },
          },
        ],
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, 'web/index.html'),
      favicon: path.resolve(__dirname, '../shared/assets/web/favicon-32.png'),
    }),
  ],
  devServer: {
    host: '0.0.0.0',
    port: 8090,
    historyApiFallback: true,
    hot: true,
  },
};
