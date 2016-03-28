'use strict';

const path              = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  entry: path.join(__dirname, '..', 'src', 'main.ts'),
  output: {
    path: path.join(__dirname, '..', 'dist'),
    filename: 'gophr.js'
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
      template: path.join(__dirname, '..', 'src', 'main.html'),
      inject: 'body',
      minify: {
        minifyJS: true,
        minifyCSS: true,
        removeComments: true,
        collapseWhitespace: true
      }
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
        loader: 'style-loader!css-loader?modules'
      },
      {
        test: /\.(jpg|png|woff)$/,
        loader: 'url-loader?limit=100000'
      },
    ]
  }
};
