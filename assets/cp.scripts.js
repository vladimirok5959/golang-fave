function GetModalAlertTmpl(title, message, error) {
	return '<div class="alert alert-' + (!error?'success':'danger') + ' alert-dismissible fade show" role="alert"><strong>' + title + '</strong> ' + message + '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>';
}

function ShowSystemMsg(title, message, error) {
	var modal_alert_place = $('.modal.show .sys-messages');
	if(!modal_alert_place.length) {
		modal_alert_place = $('form.alert-here .sys-messages');
	}
	if(modal_alert_place.length) {
		modal_alert_place.html(GetModalAlertTmpl(title, message, error));
	}
}

function ShowSystemMsgSuccess(title, message) {
	ShowSystemMsg(title, message, false);
}

function ShowSystemMsgError(title, message) {
	ShowSystemMsg(title, message, true);
}

function AjaxDone(data) {
	try {
		eval(data);
	} catch(e) {
		if(e instanceof SyntaxError) {
			console.log(data);
			console.log('JavaScript Eval Error', e.message)
		}
	}
}

function AjaxFail() {
	console.log('Form send fail, page will be reloaded');
	window.location.reload(false);
}

function ActionIndexUserLogout() {
	$.ajax({
		type: "POST",
		url: '/cp/',
		data: {
			action: 'index-user-logout',
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
			if(form.hasClass('loading')) {
				e.preventDefault();
				return;
			}

			// Block send button
			form.addClass('loading').addClass('alert-here');
			var button = $(this).find('button[type=submit]');
			button.addClass('progress-bar-striped').addClass('progress-bar-animated');

			// Another button
			if(button.attr('data-target') != '') {
				$('#' + button.attr('data-target')).addClass('progress-bar-striped').addClass('progress-bar-animated');
			}

			// Clear form messages
			form.find('.sys-messages').html('');

			$.ajax({
				type: "POST",
				url: form.attr('action'),
				data: form.serialize()
			}).done(function(data) {
				AjaxDone(data)
			}).fail(function() {
				AjaxFail();
			}).always(function() {
				// Add delay for one second
				setTimeout(function() {
					form.removeClass('loading').removeClass('alert-here');
					button.removeClass('progress-bar-striped').removeClass('progress-bar-animated');
					// Another button
					if(button.attr('data-target') != '') {
						$('#' + button.attr('data-target')).removeClass('progress-bar-striped').removeClass('progress-bar-animated');
					}
				}, 100);
			});

			e.preventDefault();
		});

		// Bind to another button
		var button = $(this).find('button[type=submit]');
		if(button.attr('data-target') != '') {
			$('#' + button.attr('data-target')).click(function() {
				button.click();
			});
		}
	});

	// Remove alert from modal on close
	$('.modal.fade').on('hidden.bs.modal', function() {
		modal_alert_place = $(this).find('.sys-messages');
		if(modal_alert_place.length) {
			modal_alert_place.html('');
		}
		// Reset form at modal close
		form = $(this).find('form');
		if(form.length) {
			form[0].reset();
		}
	}).on('show.bs.modal', function() {
		// Reset form at modal open
		form = $(this).find('form');
		if(form.length) {
			form[0].reset();
		}
	});
});