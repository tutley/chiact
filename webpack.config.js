var HtmlWebpackPlugin = require('html-webpack-plugin');
var HtmlWebpackPluginConfig = new HtmlWebpackPlugin({
  template: __dirname + '/client/src/index.html',
  filename: 'index.html',
  inject: 'body'
})

module.exports = {
  entry: [
    './client/src/index.js'
  ],
  output: {
    path: __dirname + '/client',
    filename: 'index_bundle.js'
  },
  module: {
    loaders: [
      {test: /\.js$/, exclude: /node_modules/, loader: "babel-loader"},
      { test: /\.css$/, loader: "style-loader!css-loader" }
    ]
  },
  plugins: [HtmlWebpackPluginConfig]
}
