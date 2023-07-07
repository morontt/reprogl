$(function () {
    var form = $('#comment_form');
    if (form.length) {
        var storageLocal = window.localStorage;

        initCommentator();

        var comment_url = form.attr('data-url');
        var comments_section = $('section.comments');

        function showErrors(errorsArray) {
            errorsArray.forEach(function (err) {
                var input = $('#' + err.path);
                input.addClass('error');
                input.parent().append('<small class="error">' + err.message + '</small>');
            });
        }

        function clearErrors() {
            var arr = ['name', 'mail', 'website', 'comment_text'];
            arr.forEach(function (id) {
                var input = $('#' + id);
                if (input.hasClass('error')) {
                    input.parent().find('small.error').remove();
                    input.removeClass('error');
                }
            });
        }

        comments_section.on('submit', '#comment_form', function() {
            var formData = $('#comment_form').serialize();
            $('.ajax-loader').css('display', 'flex');

            $.ajax({
                url: comment_url,
                data: formData,
                type: 'POST',
                success: function(data) {
                    saveCommentator();
                    clearErrors();
                    if (data.valid) {
                        $('#comments-wrapper').append($('#comment_add'));
                        $('#comment_text').val('');
                        $('#parentId').val(0);
                        $('#comments_thread').load(window.location.pathname + ' #comments_thread > *');
                    } else {
                        showErrors(data.errors);
                    }

                    $('.ajax-loader').hide();
                },
                error: function () {
                    clearErrors();
                    showErrors([{path: 'comment_text', message: 'Непонятная ошибка &#x1F914;'}]);
                    $('.ajax-loader').hide();
                }
            });

            return false;
        });

        comments_section.on('click', '.comment-reply span', function() {
            var parent_id = $(this).attr('data-comment-id');

            $('#form_bottom_' + parent_id).append($('#comment_add'));
            $('#parentId').val(parent_id);
        });

        comments_section.on('click', '#topic-reply span', function() {
            $('#comments-wrapper').append($('#comment_add'));
            $('#parentId').val(0);
        });

        function initCommentator() {
            if (storageLocal) {
                var nameField = $('#name');
                if (nameField.length) {
                    if (storageLocal.getItem('nickname')) {
                        nameField.val(storageLocal.getItem('nickname'));
                    }
                    if (storageLocal.getItem('email')) {
                        $('#mail').val(storageLocal.getItem('email'));
                    }
                    if (storageLocal.getItem('website')) {
                        $('#website').val(storageLocal.getItem('website'));
                    }
                }
            }
        }

        function saveCommentator() {
            if (storageLocal) {
                var nameField = $('#name');
                if (nameField.length) {
                    storageLocal.setItem('nickname', nameField.val());
                    storageLocal.setItem('email', $('#mail').val());
                    storageLocal.setItem('website', $('#website').val());
                }
            }
        }
    }
});
