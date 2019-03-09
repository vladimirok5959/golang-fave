// https://github.com/jackmoore/autosize
(function (global, factory) {
	if (typeof define === "function" && define.amd) {
		define(['module', 'exports'], factory);
	} else if (typeof exports !== "undefined") {
		factory(module, exports);
	} else {
		var mod = {
			exports: {}
		};
		factory(mod, mod.exports);
		global.autosize = mod.exports;
	}
})(this, function (module, exports) {
	'use strict';

	var map = typeof Map === "function" ? new Map() : function () {
		var keys = [];
		var values = [];

		return {
			has: function has(key) {
				return keys.indexOf(key) > -1;
			},
			get: function get(key) {
				return values[keys.indexOf(key)];
			},
			set: function set(key, value) {
				if (keys.indexOf(key) === -1) {
					keys.push(key);
					values.push(value);
				}
			},
			delete: function _delete(key) {
				var index = keys.indexOf(key);
				if (index > -1) {
					keys.splice(index, 1);
					values.splice(index, 1);
				}
			}
		};
	}();

	var createEvent = function createEvent(name) {
		return new Event(name, { bubbles: true });
	};
	try {
		new Event('test');
	} catch (e) {
		// IE does not support `new Event()`
		createEvent = function createEvent(name) {
			var evt = document.createEvent('Event');
			evt.initEvent(name, true, false);
			return evt;
		};
	}

	function assign(ta) {
		if (!ta || !ta.nodeName || ta.nodeName !== 'TEXTAREA' || map.has(ta)) return;

		var heightOffset = null;
		var clientWidth = null;
		var cachedHeight = null;

		function init() {
			var style = window.getComputedStyle(ta, null);

			if (style.resize === 'vertical') {
				ta.style.resize = 'none';
			} else if (style.resize === 'both') {
				ta.style.resize = 'horizontal';
			}

			if (style.boxSizing === 'content-box') {
				heightOffset = -(parseFloat(style.paddingTop) + parseFloat(style.paddingBottom));
			} else {
				heightOffset = parseFloat(style.borderTopWidth) + parseFloat(style.borderBottomWidth);
			}
			// Fix when a textarea is not on document body and heightOffset is Not a Number
			if (isNaN(heightOffset)) {
				heightOffset = 0;
			}

			update();
		}

		function changeOverflow(value) {
			{
				// Chrome/Safari-specific fix:
				// When the textarea y-overflow is hidden, Chrome/Safari do not reflow the text to account for the space
				// made available by removing the scrollbar. The following forces the necessary text reflow.
				var width = ta.style.width;
				ta.style.width = '0px';
				// Force reflow:
				/* jshint ignore:start */
				ta.offsetWidth;
				/* jshint ignore:end */
				ta.style.width = width;
			}

			ta.style.overflowY = value;
		}

		function getParentOverflows(el) {
			var arr = [];

			while (el && el.parentNode && el.parentNode instanceof Element) {
				if (el.parentNode.scrollTop) {
					arr.push({
						node: el.parentNode,
						scrollTop: el.parentNode.scrollTop
					});
				}
				el = el.parentNode;
			}

			return arr;
		}

		function resize() {
			if (ta.scrollHeight === 0) {
				// If the scrollHeight is 0, then the element probably has display:none or is detached from the DOM.
				return;
			}

			var overflows = getParentOverflows(ta);
			var docTop = document.documentElement && document.documentElement.scrollTop; // Needed for Mobile IE (ticket #240)

			ta.style.height = '';
			ta.style.height = ta.scrollHeight + heightOffset + 'px';

			// used to check if an update is actually necessary on window.resize
			clientWidth = ta.clientWidth;

			// prevents scroll-position jumping
			overflows.forEach(function (el) {
				el.node.scrollTop = el.scrollTop;
			});

			if (docTop) {
				document.documentElement.scrollTop = docTop;
			}
		}

		function update() {
			resize();

			var styleHeight = Math.round(parseFloat(ta.style.height));
			var computed = window.getComputedStyle(ta, null);

			// Using offsetHeight as a replacement for computed.height in IE, because IE does not account use of border-box
			var actualHeight = computed.boxSizing === 'content-box' ? Math.round(parseFloat(computed.height)) : ta.offsetHeight;

			// The actual height not matching the style height (set via the resize method) indicates that 
			// the max-height has been exceeded, in which case the overflow should be allowed.
			if (actualHeight < styleHeight) {
				if (computed.overflowY === 'hidden') {
					changeOverflow('scroll');
					resize();
					actualHeight = computed.boxSizing === 'content-box' ? Math.round(parseFloat(window.getComputedStyle(ta, null).height)) : ta.offsetHeight;
				}
			} else {
				// Normally keep overflow set to hidden, to avoid flash of scrollbar as the textarea expands.
				if (computed.overflowY !== 'hidden') {
					changeOverflow('hidden');
					resize();
					actualHeight = computed.boxSizing === 'content-box' ? Math.round(parseFloat(window.getComputedStyle(ta, null).height)) : ta.offsetHeight;
				}
			}

			if (cachedHeight !== actualHeight) {
				cachedHeight = actualHeight;
				var evt = createEvent('autosize:resized');
				try {
					ta.dispatchEvent(evt);
				} catch (err) {
					// Firefox will throw an error on dispatchEvent for a detached element
					// https://bugzilla.mozilla.org/show_bug.cgi?id=889376
				}
			}
		}

		var pageResize = function pageResize() {
			if (ta.clientWidth !== clientWidth) {
				update();
			}
		};

		var destroy = function (style) {
			window.removeEventListener('resize', pageResize, false);
			ta.removeEventListener('input', update, false);
			ta.removeEventListener('keyup', update, false);
			ta.removeEventListener('autosize:destroy', destroy, false);
			ta.removeEventListener('autosize:update', update, false);

			Object.keys(style).forEach(function (key) {
				ta.style[key] = style[key];
			});

			map.delete(ta);
		}.bind(ta, {
			height: ta.style.height,
			resize: ta.style.resize,
			overflowY: ta.style.overflowY,
			overflowX: ta.style.overflowX,
			wordWrap: ta.style.wordWrap
		});

		ta.addEventListener('autosize:destroy', destroy, false);

		// IE9 does not fire onpropertychange or oninput for deletions,
		// so binding to onkeyup to catch most of those events.
		// There is no way that I know of to detect something like 'cut' in IE9.
		if ('onpropertychange' in ta && 'oninput' in ta) {
			ta.addEventListener('keyup', update, false);
		}

		window.addEventListener('resize', pageResize, false);
		ta.addEventListener('input', update, false);
		ta.addEventListener('autosize:update', update, false);
		ta.style.overflowX = 'hidden';
		ta.style.wordWrap = 'break-word';

		map.set(ta, {
			destroy: destroy,
			update: update
		});

		init();
	}

	function destroy(ta) {
		var methods = map.get(ta);
		if (methods) {
			methods.destroy();
		}
	}

	function update(ta) {
		var methods = map.get(ta);
		if (methods) {
			methods.update();
		}
	}

	var autosize = null;

	// Do nothing in Node.js environment and IE8 (or lower)
	if (typeof window === 'undefined' || typeof window.getComputedStyle !== 'function') {
		autosize = function autosize(el) {
			return el;
		};
		autosize.destroy = function (el) {
			return el;
		};
		autosize.update = function (el) {
			return el;
		};
	} else {
		autosize = function autosize(el, options) {
			if (el) {
				Array.prototype.forEach.call(el.length ? el : [el], function (x) {
					return assign(x, options);
				});
			}
			return el;
		};
		autosize.destroy = function (el) {
			if (el) {
				Array.prototype.forEach.call(el.length ? el : [el], destroy);
			}
			return el;
		};
		autosize.update = function (el) {
			if (el) {
				Array.prototype.forEach.call(el.length ? el : [el], update);
			}
			return el;
		};
	}

	exports.default = autosize;
	module.exports = exports['default'];
});

// ---

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

		function MakeTextAreasAutoSized() {
			autosize($('textarea.autosize'));
		}

		function Initialize() {
			// Check if jQuery was loaded
			if(typeof $ == 'function') {
				AllFormsToAjax();
				BindWindowBeforeUnload();
				MakeTextAreasAutoSized();
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