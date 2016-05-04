'use strict';

const path              = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

var webpack = require('webpack');

module.exports = {
  entry: [
    path.join(__dirname, '..', 'src', 'index'),
    'webpack-dev-server/client?http://localhost:8080',
    'webpack/hot/only-dev-server',
    './src/index'
  ],
  output: {
    path: path.join(__dirname, '..', 'dist'),
    publicPath: '/',
    filename: 'gophr.js'
  },
  resolve: {
    extensions: [
      '',
      '.ts',
      '.js',
      '.jsx'
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.join(__dirname, '..', 'src', 'index.html'),
      inject: 'body',
      minify: {
        minifyJS: true,
        minifyCSS: true,
        removeComments: true,
        collapseWhitespace: true
      }
    }),
    new webpack.HotModuleReplacementPlugin()
  ],
  devServer: {
    contentBase: './dist',
    hot: true
  },
  module: {
    loaders: [
      {
        test: /\.ts$/,
        loader: 'ts-loader'
      },
      {
        test: /\.css$/,
        loader: 'style-loader!css-loader!postcss-loader?modules'
      },
      {
        test: /\.(jpg|png|woff)$/,
        loader: 'url-loader?limit=100000'
      },
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: 'react-hot!babel'
      }
    ]
  },
  postcss: function () {
    return [require('autoprefixer'), require('precss')];
  }
};

//!postcss-loader
