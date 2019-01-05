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

function AjaxDone(data) {
	try {
		eval(data);
	} catch(e) {
		if(e instanceof SyntaxError) {
			console.log(data);
			ModalShowMsg('JavaScript Eval Error', e.message)
		}
	}
}

function AjaxFail() {
	console.log('Form send fail, page will be reloaded');
	window.location.reload(false);
}

function ActionUserSettings() {
	// Reset form to remove autocomplete
	$('form.form-user-settings')[0].reset();
}

function ActionSingOut() {
	$.ajax({
		type: "POST",
		url: '/cp/',
		data: {
			action: 'singout',
		}
	}).done(function(data) {
		AjaxDone(data)
	}).fail(function() {
		AjaxFail();
	});
}

$(document).ready(function() {
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
				AjaxDone(data)
			}).fail(function() {
				AjaxFail();
			}).always(function() {
				$(form).removeClass('loading');
				$(button).removeClass('progress-bar-striped').removeClass('progress-bar-animated');
			});

			e.preventDefault();
		});
	});
});