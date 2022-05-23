module.exports = {
    publicPath: "./",
    devServer: {
        port: 80,
        proxy: 'http://127.0.0.1:8080'
    },
    pages: {
        index: {
            // entry for the page
            entry: 'src/main.js',
            // the source template
            template: 'public/index.html',
            // output as dist/index.html
            filename: 'index.html',
            // when using title option,
            // template title tag needs to be <title><%= htmlWebpackPlugin.options.title %></title>
            title: 'Tunn WebUI',
        },
    }
}