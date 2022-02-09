'use strict';

const sass = require('node-sass');

// docker-compose run nodejs su node

module.exports = function (grunt) {
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-sass');

    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        sass: {
            options: {
                implementation: sass,
                sourceMap: false
            },
            dist: {
                files: {
                    'public/assets/css/style.css': 'public/assets/sass/style.scss',
                    'public/assets/css/menu.css': 'public/assets/sass/menu.scss',
                    'public/assets/css/post.css': 'public/assets/sass/post.scss',
                    'public/assets/css/select.css': 'public/assets/sass/select.scss',
                    'public/assets/css/comments.css': 'public/assets/sass/comments.scss',
                }
            }
        },
        concat: {
            css_main: {
                options: {
                    stripBanners: {
                        block: true
                    },
                    banner: '/*! <%= pkg.name %> --- <%= grunt.template.today("dd mmm yyyy HH:MM:ss") %> */\n'
                },
                src: [
                    'node_modules/normalize.css/normalize.css',
                    'node_modules/foundation-sites/css/foundation.css',
                    'public/assets/css/vendor/magnific-popup.css',
                    'public/assets/css/vendor/owl.carousel.css',
                    'public/assets/css/vendor/owl.theme.default.css',
                    'public/assets/css/menu.css',
                    'public/assets/css/select.css',
                    'public/assets/css/style.css',
                    'public/assets/css/post.css',
                    'public/assets/css/comments.css',
                ],
                dest: 'public/assets/css/<%= pkg.name %>_main.css'
            },
            js: {
                options: {
                    stripBanners: {
                        block: true
                    },
                    banner: '/*! <%= pkg.name %> --- <%= grunt.template.today("dd mmm yyyy HH:MM:ss") %> */\n'
                },
                src: [
                    'node_modules/foundation-sites/js/vendor/modernizr.js',
                    'node_modules/jquery/dist/jquery.js',
                    'node_modules/foundation-sites/js/foundation.js',
                    'public/assets/js/vendor/headroom.js',
                    'public/assets/js/vendor/jquery.headroom.js',
                    'public/assets/js/vendor/reading-time.js',
                    'public/assets/js/vendor/imagesloaded.js',
                    'public/assets/js/vendor/nice-scroll.js',
                    'public/assets/js/vendor/shuffle.js',
                    'public/assets/js/vendor/magnific-popup.js',
                    'public/assets/js/vendor/select-fx.js',
                    'public/assets/js/vendor/owl.carousel.js',
                    'public/assets/js/vendor/velocity.js',
                    'public/assets/js/anima.js',
                ],
                dest: 'public/assets/js/<%= pkg.name %>.js'
            }
        },
        cssmin: {
            options: {
                shorthandCompacting: false,
                roundingPrecision: -1,
                sourceMap: false
            },
            target: {
                files: {
                    'public/assets/css/<%= pkg.name %>.min.css': ['<%= concat.css_main.dest %>']
                }
            }
        },
        uglify: {
            options: {
                banner: '/*! <%= pkg.name %> v<%= pkg.version %> ' +
                    '--- <%= grunt.template.today("dd mmm yyyy HH:MM:ss") %> */\n'
            },
            dist: {
                files: {
                    'public/assets/js/<%= pkg.name %>.min.js': ['<%= concat.js.dest %>']
                }
            }
        }
    });

    grunt.registerTask('build', ['sass', 'concat', 'cssmin', 'uglify']);
    grunt.registerTask('style', ['sass', 'concat', 'cssmin']);
    grunt.registerTask('default', ['build']);
};
