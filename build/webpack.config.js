'use strict';

const path              = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  entry: path.join(__dirname, '..', 'web', 'src', 'main.ts'),
  output: {
    path: path.join(__dirname, '..', 'web', 'dist'),
    filename: 'bundle.js'
  },
  resolve: {
    extensions: [
      '',
      '.ts',
      '.js'
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.join(__dirname, '..', 'web', 'src', 'main.html'),
      inject: 'body'
    })
  ],
  module: {
    loaders: [
      {
        test: /\.ts$/,
        loader: 'ts-loader'
      },
      {
        test: /\.css$/,
        loader: 'style-loader!css-loader'
      },
      {
        test: /\.(jpg|png|woff)$/,
        loader: 'url-loader?limit=100000'
      },
    ]
  }
};
