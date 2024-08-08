const esbuild = require('esbuild');

esbuild.context({
    entryPoints: ['./views/js/app.js'],
    bundle: true,
    outfile: './public/js/app.js',
    minify: true,
}).then(ctx => {
    ctx.watch();
}).catch(() => process.exit(1));
