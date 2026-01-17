const path = require('path');

module.exports = {
  mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
  entry: './src/library/index.ts',
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
    ],
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
  },
  output: {
    filename: 'library.bundle.js',
    path: path.resolve(__dirname, 'dist/library'),
    library: {
      name: 'NewsDeframer',
      type: 'window',
      export: 'default',
    },
    clean: true,
  },
  devServer: {
    static: './dist/library',
    port: 8080,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
      "Access-Control-Allow-Headers": "X-Requested-With, content-type, Authorization"
    },
    allowedHosts: "all",
    hot: false, // Hot replacement doesn't work well with injected scripts usually, standard reload is safer
    liveReload: false, // We rely on the host to re-fetch
  },
};
