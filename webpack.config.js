const path = require('path');

module.exports = {
  entry: path.join(__dirname, '/client/src/index.jsx'),
  output: {
    path: path.join(__dirname, '/client'),
    filename: '/index_bundle.js'
  },
  module: {
    loaders: [{
      test: /\.js|\.jsx?$/,
      include: path.join(__dirname, '/client/src'),
      loader: 'babel',
      query: {
        presets: ["react", "es2015"]
      }
    }, {
      test: /\.css$/,
      loader: "style-loader!css-loader"
    }]
  }
}
