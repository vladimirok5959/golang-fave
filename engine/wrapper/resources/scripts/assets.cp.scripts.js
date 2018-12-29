function ModalSysMsg(title, html) {
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
						console.log('JavaScript eval error:', e.message);
						console.log(data);
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