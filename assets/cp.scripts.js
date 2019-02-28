(function(window, $) {
	var fave = function(window, $) {
		// Private
		var FormDataWasChanged = false;

		function GetModalAlertTmpl(title, message, error) {
			return '<div class="alert alert-' + (!error?'success':'danger') + ' alert-dismissible fade show" role="alert"><strong>' + title + '</strong> ' + message + '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>';
		};

		function ShowSystemMsg(title, message, error) {
			var modal_alert_place = $('.modal.show .sys-messages');
			if(!modal_alert_place.length) {
				modal_alert_place = $('form.alert-here .sys-messages');
			}
			if(modal_alert_place.length) {
				modal_alert_place.html(GetModalAlertTmpl(title, message, error));
			}
		};

		function AjaxDone(data) {
			try {
				eval(data);
			} catch(e) {
				if(e instanceof SyntaxError) {
					console.log(data);
					console.log('Error: JavaScript code eval error', e.message)
				}
			}
		};

		function AjaxFail(data, status, error) {
			console.log('Error: data sending error, page will be reloaded', data, status, error);
			setTimeout(function() {
				window.location.reload(false);
			}, 1000);
		};

		function FormToAjax() {
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
					button.addClass('progress-bar-striped')
						.addClass('progress-bar-animated');

					// Another button
					if(button.attr('data-target') != '') {
						$('#' + button.attr('data-target')).addClass('progress-bar-striped')
							.addClass('progress-bar-animated');
					}

					// Clear form messages
					form.find('.sys-messages').html('');

					$.ajax({
						type: "POST",
						url: form.attr('action'),
						data: form.serialize()
					}).done(function(data) {
						FormDataWasChanged = false;
						AjaxDone(data)
					}).fail(function(xhr, status, error) {
						AjaxFail(xhr.responseText, status, error);
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

					// Prevent submit action
					e.preventDefault();
				});

				// Bind to another button
				var button = $(this).find('button[type=submit]');
				if(button.attr('data-target') != '') {
					$('#' + button.attr('data-target')).click(function() {
						button.click();
					});
				}

				// Mark body if any data in form was changed
				if($(this).hasClass('prev-data-lost')) {
					$(this).find('input, textarea, select').on('input', function() {
						if(!FormDataWasChanged) {
							FormDataWasChanged = true;
						}
					});
				}
			});
		};

		function FixFormInModal() {
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
		};

		function BindWindowBeforeUnload() {
			// Prevent page reload if data was changed
			$(window).bind('beforeunload', function(){
				if(FormDataWasChanged) {
					return 'Some data was changed and not saved. Are you sure want to leave page?';
				}
			});
		};

		function Initialize() {
			// Check if jQuery was loaded
			if(typeof $ == 'function') {
				FormToAjax();
				FixFormInModal();
				BindWindowBeforeUnload();
			} else {
				console.log('Error: jQuery is not loaded!');
			}
		};

		// Initialize
		if(window.addEventListener) {
			// W3C standard
			window.addEventListener('load', Initialize, false);
		} else if(window.attachEvent) {
			// Microsoft
			window.attachEvent('onload', Initialize);
		};

		// Public
		return {
			ShowMsgSuccess: function(title, message) {
				ShowSystemMsg(title, message, false);
			},

			ShowMsgError: function(title, message) {
				ShowSystemMsg(title, message, true);
			},

			ActionLogout: function(message) {
				if(confirm(message)) {
					$.ajax({
						type: "POST",
						url: '/cp/',
						data: {
							action: 'index-user-logout',
						}
					}).done(function(data) {
						AjaxDone(data)
					}).fail(function(xhr, status, error) {
						AjaxFail(xhr.responseText, status, error);
					});
				}
			},

			ActionDataTableDelete: function(object, action, id, message) {
				if(confirm(message)) {
					$.ajax({
						type: "POST",
						url: '/cp/',
						data: {
							action: action,
							id: id,
						}
					}).done(function(data) {
						AjaxDone(data)
					}).fail(function(xhr, status, error) {
						AjaxFail(xhr.responseText, status, error);
					});
				}
			},
		};
	}(window, $);

	// Make it public
	window.fave = fave;
}(window, jQuery));