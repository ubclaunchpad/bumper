var config = {
  entry: './index.js',
  output: {
    path: `${__dirname}/public/`,
    filename: 'bundle.js',
  },
  devServer: {
    port: 8080,
    inline: true,
    contentBase: './public'
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: '/node_modules/',
        loader: 'babel-loader',
        query: {
          presets: ['es2015', 'react']
        }
      }
    ]
  }
}

module.exports = config;