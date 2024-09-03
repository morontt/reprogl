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
                        $('.box').html('<span>Непонятная ошибка</span>><span class="r">&#x1F914;</span><a href="/">На главную</a>');
                        clearInterval(interval_id);
                    }
                }
            },
            error: function () {
                $('.box').html('<span>Непонятная ошибка</span><span class="r">&#x1F914;</span><a href="/">На главную</a>');
                clearInterval(interval_id);
            }
        });
    }, 1250);
});
