'use strict'

const path                  = require('path')
const webpack               = require('webpack')
const rucksack              = require('rucksack-css')
const execSync              = require('child_process').execSync
const CopyWebpackPlugin     = require('copy-webpack-plugin')
const HtmlWebpackPlugin     = require('html-webpack-plugin')
const FaviconsWebpackPlugin = require('favicons-webpack-plugin')

let minikubeIP = '127.0.0.1', gophrWebPort = 30443

// First, grab the minikube ip if not in production.
if (process.env.NODE_ENV !== 'production') {
  try {
    console.log('Attempting to get minikube IP address...')
    minikubeIP = execSync('minikube ip', { encoding: 'utf8' }).trim()
    console.log('Got minikube IP address, now starting webpack...')
  } catch(err) {
    console.error(
      `Failed to read the minikube IP address. ` +
      `Make sure the gophr development environment is running: ${err}.`)
    process.exit(1)
  }
}

module.exports = {
  context: path.join(__dirname, './client'),
  entry: {
    jsx: './index.js',
    vendor: [
      'react',
      'react-dom',
      'react-redux',
      'react-router',
      'react-router-redux',
      'redux'
    ]
  },
  output: {
    path: path.join(__dirname, './build'),
    publicPath: '/static/',
    filename: 'bundle.js',
  },
  module: {
    loaders: [
      {
        test: /\.html$/,
        loader: 'file?name=[name].[ext]'
      },
      {
        test: /\.(svg|eot|ttf|woff|woff2)$/,
        loader: 'url-loader?limit=100000'
      },
      {
        test: /\.css$/,
        include: /client/,
        loaders: [
          'style-loader',
          'css-loader?modules&sourceMap&importLoaders=1&localIdentName=[local]___[hash:base64:5]',
          'postcss-loader'
        ]
      },
      {
        test: /\.css$/,
        exclude: /client/,
        loader: 'style!css'
      },
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        loaders: [
          'react-hot',
          'babel-loader'
        ]
      },
    ],
  },
  resolve: {
    extensions: ['', '.js', '.jsx']
  },
  postcss: [
    rucksack({
      autoprefixer: true
    })
  ],
  plugins: [
    new webpack.optimize.CommonsChunkPlugin('vendor', 'vendor.bundle.js'),
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify(process.env.NODE_ENV || 'development')
      }
    }),
    new HtmlWebpackPlugin({
      title: 'gophr - Go Package Manager',
      minify: { collapseWhitespace: true },
      template: path.join(__dirname, 'client', 'index.ejs'),
      description: 'gophr is the package manager for the Go programming ' +
        'language. With gophr, managing and vendoring dependencies ' +
        'has never been easier.'
    }),
    new FaviconsWebpackPlugin(path.join(
      __dirname,
      'client',
      'resources',
      'images',
      'favicon.png')),
    new CopyWebpackPlugin([
      {
        from: path.join(
        __dirname,
        'client',
        'resources',
        'images',
        'og-splash.png'),
      },
      {
        from: path.join(
        __dirname,
        'client',
        'resources',
        'images',
        'favicon.ico'),
      }
    ])
  ],
  devServer: {
    contentBase: './client',
    hot: true,
    proxy: {
      '/': 'http://localhost:3000/static/',
      '/api/*': {
        target: `https://${minikubeIP}:${gophrWebPort}`,
        secure: false
      }
    }
  }
}
