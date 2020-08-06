/* eslint-disable */
const webpack = require('webpack');

const config = {
  mode: process.env.NODE_ENV,
  node: {
    fs: 'empty',
  },
  entry: './index.js',
  output: {
    path: `${__dirname}/public/`,
    filename: 'bundle.js',
  },
  devServer: {
    port: 8080,
    inline: true,
    contentBase: './public',
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: '/node_modules/',
        loader: 'babel-loader',
        query: {
          presets: ['@babel/env', '@babel/react'],
        },
      },
    ],
  },
  plugins: [
    new webpack.EnvironmentPlugin(['NODE_ENV'])
  ],
};

module.exports = config;
