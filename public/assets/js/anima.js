'use strict';

var anima = (function ($) {

    var

    userAgentInit = function() {
        document.documentElement.setAttribute('data-useragent', navigator.userAgent);
    },

    // main menu open/close init
    mainMenu = function () {
        $('.js-open-main-menu').on('click', function (e) {
            e.preventDefault();
            $('.js-menu-container').addClass('show-menu');
            $('.js-menu-overlay').addClass('visible');
        });
        $('.js-close-main-menu, .js-menu-overlay').on('click', function (e) {
            e.preventDefault();
            $('.js-menu-container').removeClass('show-menu');
            $('.js-menu-overlay').removeClass('visible');
        });
    },

    // nice scroll plugin init
    niceSroll = function() {
        var $niceScrollHandler = $('html:has(body.nicescroll)');
        if ($niceScrollHandler.length) {
            $niceScrollHandler.niceScroll({
                cursorcolor:'#222',
                cursorwidth: 6,
                cursorborder: '1px solid #222',
                zindex: 99999,
                cursorborderradius: 0,
                scrollspeed: 80,
                mousescrollstep: 20,
                autohidemode: 'leave',
                railpadding: {right: 2},
                smoothscroll: true
            });
        }
    },

    niceScrollReinit = function () {
        if (!$('html').getNiceScroll().length) {
            niceSroll();
        }
        document.addEventListener('visibilitychange', function() {
            if(!document.hidden && !$('html').getNiceScroll().length) {
                niceSroll();
            }
        });
    },

    niceScrollShowEvent = function () {
        document.addEventListener('visibilitychange', function() {
            if(!document.hidden && $('html').getNiceScroll().length) {
                $('html').getNiceScroll()[0].show();
            }
        });
    },

    // sticky header
    headroom = function () {
        $('.headroom').headroom({
            offset: 200,
            classes : {
                pinned : '',
                unpinned : ''
            }
        });
    },

    // go to top button
    goToTopBtn = function() {
        var $backBtn = $('.js-back-to-top-btn');
        if($backBtn.length) {
            $(window).scroll(function () {
                if ($(this).scrollTop() > 300) {
                    $backBtn.removeClass('hidden');
                } else {
                    $backBtn.addClass('hidden');
                }
            });
            $backBtn.on('click', function (e) {
                e.preventDefault();
                if ($('body.nicescroll').length) {
                    $('html').getNiceScroll()[0].doScrollTop(0, 1000);
                } else {
                    $('body, html').stop(true, false).animate({
                        scrollTop: 0
                    }, 400);
                }
            });
        }
    },

    // based on : 'Reading Position Indicator' article
    // http://css-tricks.com/reading-position-indicator/
    positionIndicator = function () {
        if ($('.js-post-reading-time').is(':visible')) {
            imagesLoaded('.post-view-article', function () {
                var getMax = function() {
                    return $(document).height() - $(window).height();
                };
                var getValue = function() {
                    return $(window).scrollTop();
                };
                var progressBar, max, value, width;
                if('max' in document.createElement('progress')){
                    // Browser supports progress element
                    progressBar = $('progress');
                    // Set the Max attr for the first time
                    progressBar.attr({ max: getMax() });
                    $(document).on('scroll', function(){
                        // On scroll only Value attr needs to be calculated
                        progressBar.attr({ value: getValue() });
                    });
                    $(window).resize(function(){
                        // On resize, both Max/Value attr needs to be calculated
                        progressBar.attr({ max: getMax(), value: getValue() });
                    });
                }
                else {
                    progressBar = $('.progress-bar'),
                        max = getMax(),
                        value, width;
                    var getWidth = function(){
                        // Calculate width in percentage
                        value = getValue();
                        width = (value/max) * 100;
                        width = width + '%';
                        return width;
                    };
                    var setWidth = function(){
                        progressBar.css({ width: getWidth() });
                    };
                    $(document).on('scroll', setWidth);
                    $(window).on('resize', function(){
                        // Need to reset the Max attr
                        max = getMax();
                        setWidth();
                    });
                }
            });
        }
    },

    readingTime = function () {
        var $postArticleContent = $('.post-article-content');
        if ($postArticleContent.length) {
            $postArticleContent.readingTime({
                wordCountTarget: $('.js-post-reading-time').find('.js-word-count'),
                lessThanAMinuteString: 'меньше минуты',
                minShortFormString: 'мин'
            });
        }
    },

    // magnific popup init
    // http://dimsemenov.com/plugins/magnific-popup/
    imageLightbox = function () {
        if ($('.anima-image-popup').length) {
            var mfOptions = {
                type: 'image',
                removalDelay: 200,
                midClick: true,
                callbacks: {
                    beforeOpen: function() {
                        this.st.image.markup = this.st.image.markup.replace('mfp-figure', 'mfp-figure mfp-with-anim');
                        this.st.mainClass = 'mfp-zoom-in';
                        /* other css animations: http://codepen.io/dimsemenov/pen/GAIkt */
                    }
                },
                closeOnContentClick: true,
                gallery: {
                    enabled: true
                },
                image: {
                    titleSrc: function (item) {
                        if (item.img && item.img.length > 0) {
                            return item.img.attr('alt');
                        }
                    }
                }
            };
            $('.anima-image-popup').magnificPopup(mfOptions);
        }
    },

    // filter tags selector styling
    filterTagsSelector = function () {
        [].slice.call(document.querySelectorAll('select.cs-select')).forEach(function(el) {
            new SelectFx(el, {
                newTab: false,
            });
        });
    },

    flassMessages = function() {
        var successMessage = $('.success-flash');
        if(successMessage.length) {
            setTimeout(function () {
                successMessage.hide(300);
            }, 3000);
        }
    },

    // anima javascripts initialization
    init = function () {
        $(document).foundation();
        userAgentInit();
        mainMenu();
        goToTopBtn();
        readingTime();
        positionIndicator();
        headroom();
        niceSroll();
        niceScrollShowEvent();
        filterTagsSelector();
        imageLightbox();
        flassMessages();
    };

    return {
        init: init
    };

})(jQuery);

(function () {
    anima.init();
})();
