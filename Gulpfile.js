const gulp = require('gulp'),
      del  = require('del'),
      fs   = require('fs'),
      seq  = require('run-sequence'),
      sync = require('gulp-sync')(gulp).sync,
      util = require('gulp-util'),
      zip  = require('gulp-zip');

const dad = require('./package.json');

// Client Tasks
const webpack          = require('webpack'),
      webpackConfig    = require('./webpack.config'),
      WebpackDevServer = require('webpack-dev-server');

gulp.task('client:webpack-dev-server', callback => {
  const port = webpackConfig.devServer.port;
  webpackConfig.entry.unshift(
    'react-hot-loader/patch',
    `webpack-dev-server/client?http://localhost:${port}/`,
    'webpack/hot/dev-server'
  );
  const compiler = webpack(webpackConfig);
  new WebpackDevServer(compiler, {
    hot: true,
    quiet: true,
    noInfo: true
  }).listen(port, 'localhost', err => {
        if (err) {
          throw new util.PluginError("webpack-dev-server", err);
        }
        callback();
    });
});

gulp.task('client:webpack', callback => {
  webpack(webpackConfig, err => {
    if (err) {
      throw new util.PluginError('webpack', err);
    }
    callback();
  });
});

gulp.task('client:dist', callback => {
  seq(
    'client:webpack',
    callback
  )
});

// Server Tasks
const child      = require('child_process'),
      git        = require('git-rev'),
      dateFormat = require('dateformat');

const now = () => dateFormat(new Date(), 'isoDateTime');

gulp.task('server:build', () => {
  const flags = `-X github.com/soprasteria/dad/cmd.Version=${dad.version}`;
  const build = child.spawnSync('go', ['install', '-ldflags', flags]);
  if (build.stderr.length) {
    util.log(util.colors.red(build.stderr.toString()));
  }
  return build;
});

gulp.task('server:watch', () => {
  const watcher = gulp.watch(['server/**/*.go', 'cmd/**/*.go'], sync([
    'server:build',
    'server:spawn'
  ], 'server:watch'));

  watcher.on('change', e => {
    util.log(util.colors.yellow(`File ${e.path} was ${e.type}`));
  });
});

let server;
gulp.task('server:spawn', () => {
  try {
    fs.mkdirSync('dist/');
  } catch (e) {}

  if (server) {
    server.kill();
  }

  server = child.spawn('dad', ['serve', '--level', 'debug'], {
    cwd: 'dist',
    stdio: 'inherit'
  });
});

gulp.task('server:dist', callback =>
  git.long(gitHash => {
    const flags = `
      -X github.com/soprasteria/dad/cmd.Version=${dad.version}
      -X github.com/soprasteria/dad/cmd.BuildDate=${now()}
      -X github.com/soprasteria/dad/cmd.GitHash=${gitHash}
    `;
    const build = child.spawnSync('go', ['build', '-o', 'dist/dad', '-ldflags', flags]);
    if (build.stderr.length) {
      util.log(util.colors.red(build.stderr.toString()));
    } else {
      callback();
    }
  })
);

// High-level tasks
gulp.task('clean', () => del(['./dist']));

gulp.task('archive', () =>
  gulp.src(['dist/**/*', '!dist/dad-*.zip'])
    .pipe(zip(`dad-${dad.version}.zip`))
    .pipe(gulp.dest('dist'))
);

gulp.task('dev', callback => {
  seq(
    'clean',
    'client:webpack-dev-server',
    'server:build',
    'server:watch',
    'server:spawn',
    callback
  );
});

gulp.task('dist', callback => {
  seq(
    'clean',
    'client:dist',
    'server:dist',
    'archive',
    callback
  );
});

const defaultTask = process.env.NODE_ENV === 'production' ? 'dist' : 'dev'
gulp.task('default', [defaultTask]);
