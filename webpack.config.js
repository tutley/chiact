const path = require('path');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var HtmlWebpackPluginConfig = new HtmlWebpackPlugin({
  template: __dirname + '/client/src/index.html',
  filename: 'index.html',
  inject: 'body'
})

module.exports = {
  entry: path.join(__dirname, '/client/src/index.jsx'),
  output: {
    path: path.join(__dirname, '/client'),
    filename: '/client/index_bundle.js'
  },
  module: {
    loaders: [{
      test: /\.jsx?$/,
      include: path.join(__dirname, '/client/src'),
      loader: 'babel',
      query: {
        presets: ["react", "es2015"]
      }
    }, {
      test: /\.css$/,
      loader: "style-loader!css-loader"
    }]
  },
  plugins: [HtmlWebpackPluginConfig]
}
