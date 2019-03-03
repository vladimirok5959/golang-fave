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

		function FormToAjax(form) {
			form.submit(function(e) {
				if(form.hasClass('loading')) {
					e.preventDefault();
					return;
				}

				// Block send button
				form.addClass('loading').addClass('alert-here');
				var button = form.find('button[type=submit]');
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
						button.removeClass('progress-bar-striped')
							.removeClass('progress-bar-animated');
						// Another button
						if(button.attr('data-target') != '') {
							$('#' + button.attr('data-target'))
								.removeClass('progress-bar-striped')
								.removeClass('progress-bar-animated');
						}
					}, 100);
				});

				// Prevent submit action
				e.preventDefault();
			});

			// Bind to another button
			var button = form.find('button[type=submit]');
			if(button.attr('data-target') != '') {
				$('#' + button.attr('data-target')).click(function() {
					button.click();
				});
			}

			// Mark body if any data in form was changed
			if(form.hasClass('prev-data-lost')) {
				form.find('input, textarea, select').on('input', function() {
					if(!FormDataWasChanged) {
						FormDataWasChanged = true;
					}
				});
			}
		};

		function AllFormsToAjax() {
			$('form').each(function() {
				FormToAjax($(this));
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
				AllFormsToAjax();
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

			ModalUserProfile: function() {
				var html = '<div class="modal fade" id="sys-modal-user-settings" tabindex="-1" role="dialog" aria-labelledby="sysModalUserSettingsLabel" aria-hidden="true"> \
					<div class="modal-dialog modal-dialog-centered" role="document"> \
						<div class="modal-content"> \
							<form class="form-user-settings" action="/cp/" method="post" autocomplete="off"> \
								<input type="hidden" name="action" value="index-user-update-profile"> \
								<div class="modal-header"> \
									<h5 class="modal-title" id="sysModalUserSettingsLabel">My profile</h5> \
									<button type="button" class="close" data-dismiss="modal" aria-label="Close"> \
										<span aria-hidden="true">&times;</span> \
									</button> \
								</div> \
								<div class="modal-body text-left"> \
									<div class="form-group"> \
										<label for="first_name">First name</label> \
										<input type="text" class="form-control" id="first_name" name="first_name" value="' + window.CurrentUserProfileData.first_name + '" placeholder="User first name" autocomplete="off"> \
									</div> \
									<div class="form-group"> \
										<label for="last_name">Last name</label> \
										<input type="text" class="form-control" id="last_name" name="last_name" value="' + window.CurrentUserProfileData.last_name + '" placeholder="User last name" autocomplete="off"> \
									</div> \
									<div class="form-group"> \
										<label for="email">Email</label> \
										<input type="email" class="form-control" id="email" name="email" value="' + window.CurrentUserProfileData.email + '" placeholder="User email" autocomplete="off" required> \
									</div> \
									<div class="form-group"> \
										<label for="password">New password</label> \
										<input type="password" class="form-control" id="password" aria-describedby="passwordHelp" name="password" value="" placeholder="User new password" autocomplete="off"> \
										<small id="passwordHelp" class="form-text text-muted">Leave this field empty if you don\'t want change your password</small> \
									</div> \
									<div class="sys-messages"></div> \
								</div> \
								<div class="modal-footer"> \
									<button type="submit" class="btn btn-primary">Save</button> \
									<button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button> \
								</div> \
							</form> \
						</div> \
					</div> \
				</div>';
				$('#sys-modal-user-settings-placeholder').html(html);
				$("#sys-modal-user-settings").modal({
					backdrop: 'static',
					keyboard: false,
					show: false,
				});
				$('#sys-modal-user-settings').on('hidden.bs.modal', function(e) {
					$('#sys-modal-user-settings-placeholder').html('');
				});
				FormToAjax($('#sys-modal-user-settings form'));
				$("#sys-modal-user-settings").modal('show');
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