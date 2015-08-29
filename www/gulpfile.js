// gulpfile.js

// --- INIT
var gulp = require('gulp');  

var connect = require('gulp-connect');
var browserify = require('browserify');
var babelify = require('babelify');
var uglify = require('gulp-uglify'); // minifies JS
var gutil = require('gulp-util');
var clean = require('gulp-clean');
var rename = require('gulp-rename');
var source = require('vinyl-source-stream');
/*
var concat = require('gulp-concat');
*/

// Paths variables
var paths = {  
    'src': {
        'base': "./src",
        'js': './src/js/',
        'css': './src/css/',
        'html': './src/html/',
        'npm': './node_modules/',
        'config': './src/config/'
    },
    'build': {
        'base': "./dist",
        'css': './dist/assets/css/',
        'js': './dist/assets/js/',
        'font': './dist/assets/font/',
        'html': './dist/'
    }
};

gulp.task('default', ['serve']);

gulp.task('serve', ['build-dev'], function() {
    gulp.watch(paths.src.js + '*.js', ['js']);
    gulp.watch(paths.src.html + '*.html', ['html']);
    gulp.watch(paths.src.css + '*.css', ['css']);
    connect.server({
        root: 'dist',
        port: 8001,
        livereload: true
    });
});

gulp.task('deploy', ['build-prod']);

gulp.task('build-dev', ['js-dev','vendor','html','css']);
gulp.task('build-prod', ['js-prod','vendor','html','css']);
gulp.task('build', ['build-dev']);


var javascriptTaskWithConfig = function(configSourceFile) {
    return browserify(paths.src.js+'app.js', {
            transform: [babelify],
            debug: true
        })
        .require(paths.src.config+configSourceFile, {expose:'config'})
        //.external(['jquery','react'])
        .bundle()
        .pipe(source('app.js'))
        //.pipe(uglify())
        .pipe(gulp.dest(paths.build.js));
};

gulp.task('js-dev', function () {
    return javascriptTaskWithConfig("development.js");
});
gulp.task('js-prod', function () {
    return javascriptTaskWithConfig("production.js");
});
gulp.task('js', ['js-dev']);


/*
gulp.task('js', function () {
    return browserify([
          paths.src.js+'*.js'
        ])
        .pipe(browserify({
            exclude: ['jquery', 'react'],
            transform: [babelify],
            debug: true
        }).on('error', gutil.log))
       // .pipe(uglify())
        .pipe(gulp.dest(paths.build.js));
});
*/


gulp.task('html', function () {
    return gulp.src([
          paths.src.html+'*.html'
        ])
        .pipe(gulp.dest(paths.build.html));
});

gulp.task('css', function () {
    return gulp.src([
          paths.src.css+'*.css'
        ])
        .pipe(gulp.dest(paths.build.css));
});

gulp.task('vendor', ['vendor-js','vendor-css','vendor-font'], function() {
});

gulp.task('vendor-js', function () {
    return gulp.src([
          paths.src.npm+'jquery/dist/jquery.js',
          paths.src.npm+'jquery/dist/jquery.min.js',
          paths.src.npm+'bootstrap/dist/js/bootstrap.js',
          paths.src.npm+'bootstrap/dist/js/bootstrap.min.js',
          paths.src.npm+'marked/marked.min.js'
        ])
        .pipe(gulp.dest(paths.build.js));
});

gulp.task('vendor-css', function () {
    return gulp.src([
          paths.src.npm+'bootstrap/dist/css/bootstrap.css',
          paths.src.npm+'bootstrap/dist/css/bootstrap.min.css',
          paths.src.npm+'bootstrap/dist/css/bootstrap-theme.min.css',
          paths.src.npm+'bootstrap/dist/css/*.map'
        ])
        .pipe(gulp.dest(paths.build.css));
});

gulp.task('vendor-font', function () {
    return gulp.src([
            paths.src.npm+'bootstrap/dist/fonts/*'
        ])
        .pipe(gulp.dest(paths.build.font));
});

gulp.task('clean', function () {
    return gulp.src(paths.build.base)
        .pipe(clean());
});
