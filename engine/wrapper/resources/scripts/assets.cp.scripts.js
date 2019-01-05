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

	// Fix navbar bottstrap menu hover
	/*
	$('#navbarCollapse ul').each(function() {
		var ul = $(this);
		ul.find('li').mouseover(function() {
			if(ul.find('li.show').length > 0) {
				var li = $(this);
				ul.find('li').removeClass('show');
				ul.find('li div.dropdown-menu').removeClass('show');
				ul.find('li a[role=button]').attr('aria-expanded', false);
				li.addClass('show');
				li.find('div.dropdown-menu').addClass('show');
				li.find('a[role=button]:first').attr('aria-expanded', true);
			}
		});
	});
	*/
});