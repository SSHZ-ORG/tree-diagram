const path = require('path')
const WebpackUserscript = require('webpack-userscript')
const dev = process.env.NODE_ENV === 'development'

module.exports = {
    mode: dev ? 'development' : 'production',
    entry: path.resolve(__dirname, 'src', 'index.js'),
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'treediagram.user.js'
    },
    devServer: {
        contentBase: path.join(__dirname, 'dist')
    },
    plugins: [
        new WebpackUserscript({
            headers: {
                name: dev ? 'TreeDiagram-dev' : 'TreeDiagram',
                version: dev ? `[version]-build.[buildNo]` : `[version]`,
                namespace: 'https://www.sshz.org/',
                description: 'Make Eventernote Better',
                match: 'https://www.eventernote.com/*',
            }
        })
    ]
}
