$(function () {
    var form = $('#comment_form');
    if (form.length) {
        var comment_url = form.attr('data-url');
        var comments_section = $('section.comments');

        comments_section.on('submit', '#comment_form', function() {
            var formData = $('#comment_form').serialize();
            $('.ajax-loader').css('display', 'flex');

            $.ajax({
                url: comment_url,
                data: formData,
                type: 'POST',
                success: function() {
                    $('#comments-wrapper').append($('#comment_add'));
                    $('#comment_text').val('');
                    $('#parentId').val(0);
                    $('#comments_thread').load(window.location.pathname + ' #comments_thread > *');

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
    }
});
