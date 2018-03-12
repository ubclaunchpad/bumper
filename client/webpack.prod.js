/* eslint-disable */
const webpack = require('webpack');
const MinifyPlugin = require('babel-minify-webpack-plugin');

const config = {
  entry: './index.js',
  output: {
    path: `${__dirname}/public/`,
    filename: 'bundle.js',
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: '/node_modules/',
        loader: 'babel-loader',
        query: {
          presets: ['es2015', 'react'],
        },
      },
    ],
  },
  plugins: [
    new webpack.EnvironmentPlugin(['NODE_ENV']),
    new MinifyPlugin(),
  ]
};

module.exports = env => config;
