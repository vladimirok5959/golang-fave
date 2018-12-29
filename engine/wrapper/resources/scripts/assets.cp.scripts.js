(function($) {
	$.fn.hasScrollBar = function() {
		return this.get(0).scrollHeight > this.get(0).clientHeight;
	}
})(jQuery);

function DetectBodyScroll() {
	var body = $('body');
	if($(body).hasScrollBar()) {
		$(body).removeClass('no-scroll');
	} else {
		$(body).addClass('no-scroll');
	}
}

function ModalSysMsg(title, html) {
	DetectBodyScroll();
	var dialog = $('#sys-modal-msg');
	$('#sysModalMsgLabel').text(title);
	$('#sysModalMsgBody').html(html);
	return dialog;
}

function ModalShowMsg(title, message) {
	var dialog = ModalSysMsg(title, message);
	dialog.modal('show');
}

$(document).ready(function() {
	// Fix body scroll
	$(window).resize(function() {
		DetectBodyScroll();
	});
	DetectBodyScroll();

	// Ajax forms
	$('form').each(function() {
		$(this).submit(function(e) {
			var form = $(this);
			if($(form).hasClass('loading')) {
				e.preventDefault();
				return;
			}

			$(form).addClass('loading');
			var button = $(this).find('button[type=submit]');
			$(button).addClass('progress-bar-striped').addClass('progress-bar-animated');

			$.ajax({
				type: "POST",
				url: form.attr('action'),
				data: form.serialize()
			}).done(function(data) {
				try {
					eval(data);
				} catch(e) {
					if(e instanceof SyntaxError) {
						console.log(data);
						ModalShowMsg('JavaScript Eval Error', e.message)
					}
				}
			}).fail(function() {
				console.log('Form send fail, page will be reloaded');
				window.location.reload(false);
			}).always(function() {
				$(form).removeClass('loading');
				$(button).removeClass('progress-bar-striped').removeClass('progress-bar-animated');
			});

			e.preventDefault();
		});
	});
});