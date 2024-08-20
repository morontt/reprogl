'use strict';

const sass = require('node-sass');

// docker-compose run nodejs su node

module.exports = function (grunt) {
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-sass');

    grunt.initConfig({
        banner_format: '/* <%= pkg.name %> v<%= pkg.version %> --- <%= grunt.template.today("dd mmm yyyy HH:MM:ss o") %> */\n',
        pkg: grunt.file.readJSON('package.json'),
        sass: {
            options: {
                implementation: sass,
                sourceMap: false
            },
            dist: {
                files: {
                    'public/assets/css/foundation.css': 'public/assets/sass/foundation.scss',
                    'public/assets/css/style.css': 'public/assets/sass/style.scss',
                    'public/assets/css/menu.css': 'public/assets/sass/menu.scss',
                    'public/assets/css/post.css': 'public/assets/sass/post.scss',
                    'public/assets/css/select.css': 'public/assets/sass/select.scss',
                    'public/assets/css/comments.css': 'public/assets/sass/comments.scss',
                    'public/assets/css/profile.css': 'public/assets/sass/profile.scss',
                    'public/assets/css/pygments.css': 'public/assets/sass/pygments.scss',
                    'public/assets/css/errors.css': 'public/assets/sass/errors.scss',
                    'public/assets/css/glyphicons.css': 'public/assets/sass/glyphicons.scss',
                    'public/assets/css/login.css': 'public/assets/sass/login.scss',
                }
            }
        },
        concat: {
            css_main: {
                options: {
                    stripBanners: {
                        block: true
                    },
                    banner: '<%= banner_format %>'
                },
                src: [
                    'node_modules/normalize.css/normalize.css',
                    'public/assets/css/foundation.css',
                    'node_modules/magnific-popup/dist/magnific-popup.css',
                    'public/assets/css/glyphicons.css',
                    'public/assets/css/style.css',
                    'public/assets/css/pygments.css',
                    'public/assets/css/menu.css',
                    'public/assets/css/select.css',
                    'public/assets/css/post.css',
                    'public/assets/css/comments.css',
                    'public/assets/css/profile.css',
                ],
                dest: 'public/assets/css/<%= pkg.name %>.css'
            },
            foundation_js: {
                src: [
                    'node_modules/foundation-sites/js/foundation/foundation.js',
                ],
                dest: 'public/assets/js/foundation.js'
            },
            js: {
                options: {
                    stripBanners: {
                        block: true
                    },
                    banner: '<%= banner_format %>'
                },
                src: [
                    'node_modules/foundation-sites/js/vendor/modernizr.js',
                    'node_modules/jquery/dist/jquery.js',
                    '<%= concat.foundation_js.dest %>',
                    'public/assets/js/vendors/headroom.js',
                    'public/assets/js/vendors/jquery.headroom.js',
                    'public/assets/js/vendors/reading-time.js',
                    'public/assets/js/vendors/imagesloaded.js',
                    'public/assets/js/vendors/nice-scroll.js',
                    'public/assets/js/vendors/shuffle.js',
                    'node_modules/magnific-popup/dist/jquery.magnific-popup.js',
                    'public/assets/js/vendors/select-fx.js',
                    'public/assets/js/vendors/velocity.js',
                    'public/assets/js/anima.js',
                    'public/assets/js/comments.js',
                ],
                dest: 'public/assets/js/<%= pkg.name %>.js'
            }
        },
        cssmin: {
            options: {
                shorthandCompacting: false,
                roundingPrecision: -1,
                format: 'keep-breaks',
                sourceMap: false
            },
            target: {
                files: {
                    'public/assets/css/<%= pkg.name %>_temp.min.css': ['<%= concat.css_main.dest %>']
                }
            }
        },
        uglify: {
            options: {
                output: {
                    ascii_only: true,
                    max_line_len: 160
                }
            },
            dist: {
                files: {
                    'public/assets/js/<%= pkg.name %>_temp.min.js': ['<%= concat.js.dest %>']
                }
            }
        }
    });

    grunt.registerTask('build', ['sass', 'concat', 'cssmin', 'uglify']);
    grunt.registerTask('style', ['sass', 'concat', 'cssmin']);
    grunt.registerTask('default', ['build']);
};
