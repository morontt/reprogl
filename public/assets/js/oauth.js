$(function () {
    var request_id = $('.box').data('request-id');
    var interval_id = setInterval(function () {
        $.ajax({
            url: '/oauth/check/' + request_id,
            type: 'GET',
            success: function (data) {
                if (data.status) {
                    if (data.status === 'ok' && data.redirect_url) {
                        clearInterval(interval_id);
                        window.location.assign(data.redirect_url);
                    } else if (data.status === 'error') {
                        $('.box').html('Непонятная ошибка <span>&#x1F914;</span><br><a href="/">На главную</a>');
                        clearInterval(interval_id);
                    }
                }
            },
            error: function () {
                $('.box').html('Непонятная ошибка <span>&#x1F914;</span><br><a href="/">На главную</a>');
                clearInterval(interval_id);
            }
        });
    }, 1250);
});
