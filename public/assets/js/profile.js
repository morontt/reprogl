$(function () {
    var form = $('#profile_form');
    if (form.length) {
        var submit_url = form.attr('action');
        var defaultField = 'email';
        var allFields = ['displayName', 'email', 'username'];

        form.on('submit', function() {
            var form_data = form.serialize();
            $('.ajax-loader').css('display', 'flex');

            $.ajax({
                url: submit_url,
                data: form_data,
                type: 'POST',
                success: function(data) {
                    clearErrors(allFields);
                    if (data.valid) {
                        $('#profile_info').load(window.location.pathname + ' #profile_info > *', function () {
                            $('.ajax-loader').hide();
                        });
                    } else {
                        showErrors(data.errors, defaultField);
                        $('.ajax-loader').hide();
                    }
                },
                error: function () {
                    clearErrors(allFields);
                    showErrors([{path: defaultField, message: 'Непонятная ошибка &#x1F914;'}], defaultField);
                    $('.ajax-loader').hide();
                }
            });

            return false;
        });
    }
});
