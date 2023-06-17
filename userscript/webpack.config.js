const path = require('path')
const { UserscriptPlugin } = require('webpack-userscript')
const ClosurePlugin = require('closure-webpack-plugin')
const dev = process.env.NODE_ENV === 'development'
const version = require('./package.json').version

module.exports = {
    mode: dev ? 'development' : 'production',
    entry: path.resolve(__dirname, 'src', 'index.ts'),
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'treediagram.user.js'
    },
    devServer: {
        static: {
            directory: path.join(__dirname, 'dist'),
        },
    },
    plugins: [
        new UserscriptPlugin({
            headers: {
                name: dev ? 'TreeDiagram-dev' : 'TreeDiagram',
                version: dev ? `${version}-build.[buildNo]` : `${version}`,
                namespace: 'https://www.sshz.org/',
                description: 'Make Eventernote Better',
                match: 'https://www.eventernote.com/*',
            }
        })
    ],
    optimization: {
        minimizer: [
            new ClosurePlugin({mode: 'STANDARD'}, {
                strict_mode_input: false,
                language_out: 'ECMASCRIPT_2019'
            })
        ]
    }
}
