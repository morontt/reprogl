$(function () {
    var main = $('main.post-view-content');
    if (main.length) {
        var submit_url = main.data('trkUrl');
        if (submit_url) {
            $.ajax({
                url: submit_url,
                type: 'POST',
            });
        }
    }
});
