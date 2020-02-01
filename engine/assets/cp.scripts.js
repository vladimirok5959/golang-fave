/*!
 * Bootstrap-select v1.13.9 (https://developer.snapappointments.com/bootstrap-select)
 *
 * Copyright 2012-2019 SnapAppointments, LLC
 * Licensed under MIT (https://github.com/snapappointments/bootstrap-select/blob/master/LICENSE)
 */

(function (root, factory) {
	if (root === undefined && window !== undefined) root = window;
	if (typeof define === 'function' && define.amd) {
		// AMD. Register as an anonymous module unless amdModuleId is set
		define(["jquery"], function (a0) {
			return (factory(a0));
		});
	} else if (typeof module === 'object' && module.exports) {
		// Node. Does not work with strict CommonJS, but
		// only CommonJS-like environments that support module.exports,
		// like Node.
		module.exports = factory(require("jquery"));
	} else {
		factory(root["jQuery"]);
	}
}(this, function (jQuery) {

(function ($) {
	'use strict';

	var DISALLOWED_ATTRIBUTES = ['sanitize', 'whiteList', 'sanitizeFn'];

	var uriAttrs = [
		'background',
		'cite',
		'href',
		'itemtype',
		'longdesc',
		'poster',
		'src',
		'xlink:href'
	];

	var ARIA_ATTRIBUTE_PATTERN = /^aria-[\w-]*$/i;

	var DefaultWhitelist = {
		// Global attributes allowed on any supplied element below.
		'*': ['class', 'dir', 'id', 'lang', 'role', 'tabindex', 'style', ARIA_ATTRIBUTE_PATTERN],
		a: ['target', 'href', 'title', 'rel'],
		area: [],
		b: [],
		br: [],
		col: [],
		code: [],
		div: [],
		em: [],
		hr: [],
		h1: [],
		h2: [],
		h3: [],
		h4: [],
		h5: [],
		h6: [],
		i: [],
		img: ['src', 'alt', 'title', 'width', 'height'],
		li: [],
		ol: [],
		p: [],
		pre: [],
		s: [],
		small: [],
		span: [],
		sub: [],
		sup: [],
		strong: [],
		u: [],
		ul: []
	}

	/**
	 * A pattern that recognizes a commonly useful subset of URLs that are safe.
	 *
	 * Shoutout to Angular 7 https://github.com/angular/angular/blob/7.2.4/packages/core/src/sanitization/url_sanitizer.ts
	 */
	var SAFE_URL_PATTERN = /^(?:(?:https?|mailto|ftp|tel|file):|[^&:/?#]*(?:[/?#]|$))/gi;

	/**
	 * A pattern that matches safe data URLs. Only matches image, video and audio types.
	 *
	 * Shoutout to Angular 7 https://github.com/angular/angular/blob/7.2.4/packages/core/src/sanitization/url_sanitizer.ts
	 */
	var DATA_URL_PATTERN = /^data:(?:image\/(?:bmp|gif|jpeg|jpg|png|tiff|webp)|video\/(?:mpeg|mp4|ogg|webm)|audio\/(?:mp3|oga|ogg|opus));base64,[a-z0-9+/]+=*$/i;

	function allowedAttribute (attr, allowedAttributeList) {
		var attrName = attr.nodeName.toLowerCase()

		if ($.inArray(attrName, allowedAttributeList) !== -1) {
			if ($.inArray(attrName, uriAttrs) !== -1) {
				return Boolean(attr.nodeValue.match(SAFE_URL_PATTERN) || attr.nodeValue.match(DATA_URL_PATTERN))
			}

			return true
		}

		var regExp = $(allowedAttributeList).filter(function (index, value) {
			return value instanceof RegExp
		})

		// Check if a regular expression validates the attribute.
		for (var i = 0, l = regExp.length; i < l; i++) {
			if (attrName.match(regExp[i])) {
				return true
			}
		}

		return false
	}

	function sanitizeHtml (unsafeElements, whiteList, sanitizeFn) {
		if (sanitizeFn && typeof sanitizeFn === 'function') {
			return sanitizeFn(unsafeElements);
		}

		var whitelistKeys = Object.keys(whiteList);

		for (var i = 0, len = unsafeElements.length; i < len; i++) {
			var elements = unsafeElements[i].querySelectorAll('*');

			for (var j = 0, len2 = elements.length; j < len2; j++) {
				var el = elements[j];
				var elName = el.nodeName.toLowerCase();

				if (whitelistKeys.indexOf(elName) === -1) {
					el.parentNode.removeChild(el);

					continue;
				}

				var attributeList = [].slice.call(el.attributes);
				var whitelistedAttributes = [].concat(whiteList['*'] || [], whiteList[elName] || []);

				for (var k = 0, len3 = attributeList.length; k < len3; k++) {
					var attr = attributeList[k];

					if (!allowedAttribute(attr, whitelistedAttributes)) {
						el.removeAttribute(attr.nodeName);
					}
				}
			}
		}
	}

	// Polyfill for browsers with no classList support
	// Remove in v2
	if (!('classList' in document.createElement('_'))) {
		(function (view) {
			if (!('Element' in view)) return;

			var classListProp = 'classList',
					protoProp = 'prototype',
					elemCtrProto = view.Element[protoProp],
					objCtr = Object,
					classListGetter = function () {
						var $elem = $(this);

						return {
							add: function (classes) {
								classes = Array.prototype.slice.call(arguments).join(' ');
								return $elem.addClass(classes);
							},
							remove: function (classes) {
								classes = Array.prototype.slice.call(arguments).join(' ');
								return $elem.removeClass(classes);
							},
							toggle: function (classes, force) {
								return $elem.toggleClass(classes, force);
							},
							contains: function (classes) {
								return $elem.hasClass(classes);
							}
						}
					};

			if (objCtr.defineProperty) {
				var classListPropDesc = {
					get: classListGetter,
					enumerable: true,
					configurable: true
				};
				try {
					objCtr.defineProperty(elemCtrProto, classListProp, classListPropDesc);
				} catch (ex) { // IE 8 doesn't support enumerable:true
					// adding undefined to fight this issue https://github.com/eligrey/classList.js/issues/36
					// modernie IE8-MSW7 machine has IE8 8.0.6001.18702 and is affected
					if (ex.number === undefined || ex.number === -0x7FF5EC54) {
						classListPropDesc.enumerable = false;
						objCtr.defineProperty(elemCtrProto, classListProp, classListPropDesc);
					}
				}
			} else if (objCtr[protoProp].__defineGetter__) {
				elemCtrProto.__defineGetter__(classListProp, classListGetter);
			}
		}(window));
	}

	var testElement = document.createElement('_');

	testElement.classList.add('c1', 'c2');

	if (!testElement.classList.contains('c2')) {
		var _add = DOMTokenList.prototype.add,
				_remove = DOMTokenList.prototype.remove;

		DOMTokenList.prototype.add = function () {
			Array.prototype.forEach.call(arguments, _add.bind(this));
		}

		DOMTokenList.prototype.remove = function () {
			Array.prototype.forEach.call(arguments, _remove.bind(this));
		}
	}

	testElement.classList.toggle('c3', false);

	// Polyfill for IE 10 and Firefox <24, where classList.toggle does not
	// support the second argument.
	if (testElement.classList.contains('c3')) {
		var _toggle = DOMTokenList.prototype.toggle;

		DOMTokenList.prototype.toggle = function (token, force) {
			if (1 in arguments && !this.contains(token) === !force) {
				return force;
			} else {
				return _toggle.call(this, token);
			}
		};
	}

	testElement = null;

	// shallow array comparison
	function isEqual (array1, array2) {
		return array1.length === array2.length && array1.every(function (element, index) {
			return element === array2[index];
		});
	};

	// <editor-fold desc="Shims">
	if (!String.prototype.startsWith) {
		(function () {
			'use strict';
			var defineProperty = (function () {
				try {
					var object = {};
					var $defineProperty = Object.defineProperty;
					var result = $defineProperty(object, object, object) && $defineProperty;
				} catch (error) {
				}
				return result;
			}());
			var toString = {}.toString;
			var startsWith = function (search) {
				if (this == null) {
					throw new TypeError();
				}
				var string = String(this);
				if (search && toString.call(search) == '[object RegExp]') {
					throw new TypeError();
				}
				var stringLength = string.length;
				var searchString = String(search);
				var searchLength = searchString.length;
				var position = arguments.length > 1 ? arguments[1] : undefined;
				var pos = position ? Number(position) : 0;
				if (pos != pos) {
					pos = 0;
				}
				var start = Math.min(Math.max(pos, 0), stringLength);
				if (searchLength + start > stringLength) {
					return false;
				}
				var index = -1;
				while (++index < searchLength) {
					if (string.charCodeAt(start + index) != searchString.charCodeAt(index)) {
						return false;
					}
				}
				return true;
			};
			if (defineProperty) {
				defineProperty(String.prototype, 'startsWith', {
					'value': startsWith,
					'configurable': true,
					'writable': true
				});
			} else {
				String.prototype.startsWith = startsWith;
			}
		}());
	}

	if (!Object.keys) {
		Object.keys = function (
			o, // object
			k, // key
			r  // result array
		) {
			// initialize object and result
			r = [];
			// iterate over object keys
			for (k in o) {
				// fill result array with non-prototypical keys
				r.hasOwnProperty.call(o, k) && r.push(k);
			}
			// return result
			return r;
		};
	}

	if (HTMLSelectElement && !HTMLSelectElement.prototype.hasOwnProperty('selectedOptions')) {
		Object.defineProperty(HTMLSelectElement.prototype, 'selectedOptions', {
			get: function () {
				return this.querySelectorAll(':checked');
			}
		});
	}

	// much faster than $.val()
	function getSelectValues (select) {
		var result = [];
		var options = select.selectedOptions;
		var opt;

		if (select.multiple) {
			for (var i = 0, len = options.length; i < len; i++) {
				opt = options[i];

				result.push(opt.value || opt.text);
			}
		} else {
			result = select.value;
		}

		return result;
	}

	// set data-selected on select element if the value has been programmatically selected
	// prior to initialization of bootstrap-select
	// * consider removing or replacing an alternative method *
	var valHooks = {
		useDefault: false,
		_set: $.valHooks.select.set
	};

	$.valHooks.select.set = function (elem, value) {
		if (value && !valHooks.useDefault) $(elem).data('selected', true);

		return valHooks._set.apply(this, arguments);
	};

	var changedArguments = null;

	var EventIsSupported = (function () {
		try {
			new Event('change');
			return true;
		} catch (e) {
			return false;
		}
	})();

	$.fn.triggerNative = function (eventName) {
		var el = this[0],
				event;

		if (el.dispatchEvent) { // for modern browsers & IE9+
			if (EventIsSupported) {
				// For modern browsers
				event = new Event(eventName, {
					bubbles: true
				});
			} else {
				// For IE since it doesn't support Event constructor
				event = document.createEvent('Event');
				event.initEvent(eventName, true, false);
			}

			el.dispatchEvent(event);
		} else if (el.fireEvent) { // for IE8
			event = document.createEventObject();
			event.eventType = eventName;
			el.fireEvent('on' + eventName, event);
		} else {
			// fall back to jQuery.trigger
			this.trigger(eventName);
		}
	};
	// </editor-fold>

	function stringSearch (li, searchString, method, normalize) {
		var stringTypes = [
					'display',
					'subtext',
					'tokens'
				],
				searchSuccess = false;

		for (var i = 0; i < stringTypes.length; i++) {
			var stringType = stringTypes[i],
					string = li[stringType];

			if (string) {
				string = string.toString();

				// Strip HTML tags. This isn't perfect, but it's much faster than any other method
				if (stringType === 'display') {
					string = string.replace(/<[^>]+>/g, '');
				}

				if (normalize) string = normalizeToBase(string);
				string = string.toUpperCase();

				if (method === 'contains') {
					searchSuccess = string.indexOf(searchString) >= 0;
				} else {
					searchSuccess = string.startsWith(searchString);
				}

				if (searchSuccess) break;
			}
		}

		return searchSuccess;
	}

	function toInteger (value) {
		return parseInt(value, 10) || 0;
	}

	// Borrowed from Lodash (_.deburr)
	/** Used to map Latin Unicode letters to basic Latin letters. */
	var deburredLetters = {
		// Latin-1 Supplement block.
		'\xc0': 'A',  '\xc1': 'A', '\xc2': 'A', '\xc3': 'A', '\xc4': 'A', '\xc5': 'A',
		'\xe0': 'a',  '\xe1': 'a', '\xe2': 'a', '\xe3': 'a', '\xe4': 'a', '\xe5': 'a',
		'\xc7': 'C',  '\xe7': 'c',
		'\xd0': 'D',  '\xf0': 'd',
		'\xc8': 'E',  '\xc9': 'E', '\xca': 'E', '\xcb': 'E',
		'\xe8': 'e',  '\xe9': 'e', '\xea': 'e', '\xeb': 'e',
		'\xcc': 'I',  '\xcd': 'I', '\xce': 'I', '\xcf': 'I',
		'\xec': 'i',  '\xed': 'i', '\xee': 'i', '\xef': 'i',
		'\xd1': 'N',  '\xf1': 'n',
		'\xd2': 'O',  '\xd3': 'O', '\xd4': 'O', '\xd5': 'O', '\xd6': 'O', '\xd8': 'O',
		'\xf2': 'o',  '\xf3': 'o', '\xf4': 'o', '\xf5': 'o', '\xf6': 'o', '\xf8': 'o',
		'\xd9': 'U',  '\xda': 'U', '\xdb': 'U', '\xdc': 'U',
		'\xf9': 'u',  '\xfa': 'u', '\xfb': 'u', '\xfc': 'u',
		'\xdd': 'Y',  '\xfd': 'y', '\xff': 'y',
		'\xc6': 'Ae', '\xe6': 'ae',
		'\xde': 'Th', '\xfe': 'th',
		'\xdf': 'ss',
		// Latin Extended-A block.
		'\u0100': 'A',  '\u0102': 'A', '\u0104': 'A',
		'\u0101': 'a',  '\u0103': 'a', '\u0105': 'a',
		'\u0106': 'C',  '\u0108': 'C', '\u010a': 'C', '\u010c': 'C',
		'\u0107': 'c',  '\u0109': 'c', '\u010b': 'c', '\u010d': 'c',
		'\u010e': 'D',  '\u0110': 'D', '\u010f': 'd', '\u0111': 'd',
		'\u0112': 'E',  '\u0114': 'E', '\u0116': 'E', '\u0118': 'E', '\u011a': 'E',
		'\u0113': 'e',  '\u0115': 'e', '\u0117': 'e', '\u0119': 'e', '\u011b': 'e',
		'\u011c': 'G',  '\u011e': 'G', '\u0120': 'G', '\u0122': 'G',
		'\u011d': 'g',  '\u011f': 'g', '\u0121': 'g', '\u0123': 'g',
		'\u0124': 'H',  '\u0126': 'H', '\u0125': 'h', '\u0127': 'h',
		'\u0128': 'I',  '\u012a': 'I', '\u012c': 'I', '\u012e': 'I', '\u0130': 'I',
		'\u0129': 'i',  '\u012b': 'i', '\u012d': 'i', '\u012f': 'i', '\u0131': 'i',
		'\u0134': 'J',  '\u0135': 'j',
		'\u0136': 'K',  '\u0137': 'k', '\u0138': 'k',
		'\u0139': 'L',  '\u013b': 'L', '\u013d': 'L', '\u013f': 'L', '\u0141': 'L',
		'\u013a': 'l',  '\u013c': 'l', '\u013e': 'l', '\u0140': 'l', '\u0142': 'l',
		'\u0143': 'N',  '\u0145': 'N', '\u0147': 'N', '\u014a': 'N',
		'\u0144': 'n',  '\u0146': 'n', '\u0148': 'n', '\u014b': 'n',
		'\u014c': 'O',  '\u014e': 'O', '\u0150': 'O',
		'\u014d': 'o',  '\u014f': 'o', '\u0151': 'o',
		'\u0154': 'R',  '\u0156': 'R', '\u0158': 'R',
		'\u0155': 'r',  '\u0157': 'r', '\u0159': 'r',
		'\u015a': 'S',  '\u015c': 'S', '\u015e': 'S', '\u0160': 'S',
		'\u015b': 's',  '\u015d': 's', '\u015f': 's', '\u0161': 's',
		'\u0162': 'T',  '\u0164': 'T', '\u0166': 'T',
		'\u0163': 't',  '\u0165': 't', '\u0167': 't',
		'\u0168': 'U',  '\u016a': 'U', '\u016c': 'U', '\u016e': 'U', '\u0170': 'U', '\u0172': 'U',
		'\u0169': 'u',  '\u016b': 'u', '\u016d': 'u', '\u016f': 'u', '\u0171': 'u', '\u0173': 'u',
		'\u0174': 'W',  '\u0175': 'w',
		'\u0176': 'Y',  '\u0177': 'y', '\u0178': 'Y',
		'\u0179': 'Z',  '\u017b': 'Z', '\u017d': 'Z',
		'\u017a': 'z',  '\u017c': 'z', '\u017e': 'z',
		'\u0132': 'IJ', '\u0133': 'ij',
		'\u0152': 'Oe', '\u0153': 'oe',
		'\u0149': "'n", '\u017f': 's'
	};

	/** Used to match Latin Unicode letters (excluding mathematical operators). */
	var reLatin = /[\xc0-\xd6\xd8-\xf6\xf8-\xff\u0100-\u017f]/g;

	/** Used to compose unicode character classes. */
	var rsComboMarksRange = '\\u0300-\\u036f',
			reComboHalfMarksRange = '\\ufe20-\\ufe2f',
			rsComboSymbolsRange = '\\u20d0-\\u20ff',
			rsComboMarksExtendedRange = '\\u1ab0-\\u1aff',
			rsComboMarksSupplementRange = '\\u1dc0-\\u1dff',
			rsComboRange = rsComboMarksRange + reComboHalfMarksRange + rsComboSymbolsRange + rsComboMarksExtendedRange + rsComboMarksSupplementRange;

	/** Used to compose unicode capture groups. */
	var rsCombo = '[' + rsComboRange + ']';

	/**
	 * Used to match [combining diacritical marks](https://en.wikipedia.org/wiki/Combining_Diacritical_Marks) and
	 * [combining diacritical marks for symbols](https://en.wikipedia.org/wiki/Combining_Diacritical_Marks_for_Symbols).
	 */
	var reComboMark = RegExp(rsCombo, 'g');

	function deburrLetter (key) {
		return deburredLetters[key];
	};

	function normalizeToBase (string) {
		string = string.toString();
		return string && string.replace(reLatin, deburrLetter).replace(reComboMark, '');
	}

	// List of HTML entities for escaping.
	var escapeMap = {
		'&': '&amp;',
		'<': '&lt;',
		'>': '&gt;',
		'"': '&quot;',
		"'": '&#x27;',
	};
	escapeMap[String.fromCharCode(96)] = '&#x60;';

	// Functions for escaping and unescaping strings to/from HTML interpolation.
	var createEscaper = function (map) {
		var escaper = function (match) {
			return map[match];
		};
		// Regexes for identifying a key that needs to be escaped.
		var source = '(?:' + Object.keys(map).join('|') + ')';
		var testRegexp = RegExp(source);
		var replaceRegexp = RegExp(source, 'g');
		return function (string) {
			string = string == null ? '' : '' + string;
			return testRegexp.test(string) ? string.replace(replaceRegexp, escaper) : string;
		};
	};

	var htmlEscape = createEscaper(escapeMap);

	/**
	 * ------------------------------------------------------------------------
	 * Constants
	 * ------------------------------------------------------------------------
	 */

	var keyCodeMap = {
		32: ' ',
		48: '0',
		49: '1',
		50: '2',
		51: '3',
		52: '4',
		53: '5',
		54: '6',
		55: '7',
		56: '8',
		57: '9',
		59: ';',
		65: 'A',
		66: 'B',
		67: 'C',
		68: 'D',
		69: 'E',
		70: 'F',
		71: 'G',
		72: 'H',
		73: 'I',
		74: 'J',
		75: 'K',
		76: 'L',
		77: 'M',
		78: 'N',
		79: 'O',
		80: 'P',
		81: 'Q',
		82: 'R',
		83: 'S',
		84: 'T',
		85: 'U',
		86: 'V',
		87: 'W',
		88: 'X',
		89: 'Y',
		90: 'Z',
		96: '0',
		97: '1',
		98: '2',
		99: '3',
		100: '4',
		101: '5',
		102: '6',
		103: '7',
		104: '8',
		105: '9'
	};

	var keyCodes = {
		ESCAPE: 27, // KeyboardEvent.which value for Escape (Esc) key
		ENTER: 13, // KeyboardEvent.which value for Enter key
		SPACE: 32, // KeyboardEvent.which value for space key
		TAB: 9, // KeyboardEvent.which value for tab key
		ARROW_UP: 38, // KeyboardEvent.which value for up arrow key
		ARROW_DOWN: 40 // KeyboardEvent.which value for down arrow key
	}

	var version = {
		success: false,
		major: '3'
	};

	try {
		version.full = ($.fn.dropdown.Constructor.VERSION || '').split(' ')[0].split('.');
		version.major = version.full[0];
		version.success = true;
	} catch (err) {
		// do nothing
	}

	var selectId = 0;

	var EVENT_KEY = '.bs.select';

	var classNames = {
		DISABLED: 'disabled',
		DIVIDER: 'divider',
		SHOW: 'open',
		DROPUP: 'dropup',
		MENU: 'dropdown-menu',
		MENURIGHT: 'dropdown-menu-right',
		MENULEFT: 'dropdown-menu-left',
		// to-do: replace with more advanced template/customization options
		BUTTONCLASS: 'btn-default',
		POPOVERHEADER: 'popover-title',
		ICONBASE: 'glyphicon',
		TICKICON: 'glyphicon-ok'
	}

	var Selector = {
		MENU: '.' + classNames.MENU
	}

	var elementTemplates = {
		span: document.createElement('span'),
		i: document.createElement('i'),
		subtext: document.createElement('small'),
		a: document.createElement('a'),
		li: document.createElement('li'),
		whitespace: document.createTextNode('\u00A0'),
		fragment: document.createDocumentFragment()
	}

	elementTemplates.a.setAttribute('role', 'option');
	elementTemplates.subtext.className = 'text-muted';

	elementTemplates.text = elementTemplates.span.cloneNode(false);
	elementTemplates.text.className = 'text';

	elementTemplates.checkMark = elementTemplates.span.cloneNode(false);

	var REGEXP_ARROW = new RegExp(keyCodes.ARROW_UP + '|' + keyCodes.ARROW_DOWN);
	var REGEXP_TAB_OR_ESCAPE = new RegExp('^' + keyCodes.TAB + '$|' + keyCodes.ESCAPE);

	var generateOption = {
		li: function (content, classes, optgroup) {
			var li = elementTemplates.li.cloneNode(false);

			if (content) {
				if (content.nodeType === 1 || content.nodeType === 11) {
					li.appendChild(content);
				} else {
					li.innerHTML = content;
				}
			}

			if (typeof classes !== 'undefined' && classes !== '') li.className = classes;
			if (typeof optgroup !== 'undefined' && optgroup !== null) li.classList.add('optgroup-' + optgroup);

			return li;
		},

		a: function (text, classes, inline) {
			var a = elementTemplates.a.cloneNode(true);

			if (text) {
				if (text.nodeType === 11) {
					a.appendChild(text);
				} else {
					a.insertAdjacentHTML('beforeend', text);
				}
			}

			if (typeof classes !== 'undefined' && classes !== '') a.className = classes;
			if (version.major === '4') a.classList.add('dropdown-item');
			if (inline) a.setAttribute('style', inline);

			return a;
		},

		text: function (options, useFragment) {
			var textElement = elementTemplates.text.cloneNode(false),
					subtextElement,
					iconElement;

			if (options.content) {
				textElement.innerHTML = options.content;
			} else {
				textElement.textContent = options.text;

				if (options.icon) {
					var whitespace = elementTemplates.whitespace.cloneNode(false);

					// need to use <i> for icons in the button to prevent a breaking change
					// note: switch to span in next major release
					iconElement = (useFragment === true ? elementTemplates.i : elementTemplates.span).cloneNode(false);
					iconElement.className = options.iconBase + ' ' + options.icon;

					elementTemplates.fragment.appendChild(iconElement);
					elementTemplates.fragment.appendChild(whitespace);
				}

				if (options.subtext) {
					subtextElement = elementTemplates.subtext.cloneNode(false);
					subtextElement.textContent = options.subtext;
					textElement.appendChild(subtextElement);
				}
			}

			if (useFragment === true) {
				while (textElement.childNodes.length > 0) {
					elementTemplates.fragment.appendChild(textElement.childNodes[0]);
				}
			} else {
				elementTemplates.fragment.appendChild(textElement);
			}

			return elementTemplates.fragment;
		},

		label: function (options) {
			var textElement = elementTemplates.text.cloneNode(false),
					subtextElement,
					iconElement;

			textElement.innerHTML = options.label;

			if (options.icon) {
				var whitespace = elementTemplates.whitespace.cloneNode(false);

				iconElement = elementTemplates.span.cloneNode(false);
				iconElement.className = options.iconBase + ' ' + options.icon;

				elementTemplates.fragment.appendChild(iconElement);
				elementTemplates.fragment.appendChild(whitespace);
			}

			if (options.subtext) {
				subtextElement = elementTemplates.subtext.cloneNode(false);
				subtextElement.textContent = options.subtext;
				textElement.appendChild(subtextElement);
			}

			elementTemplates.fragment.appendChild(textElement);

			return elementTemplates.fragment;
		}
	}

	var Selectpicker = function (element, options) {
		var that = this;

		// bootstrap-select has been initialized - revert valHooks.select.set back to its original function
		if (!valHooks.useDefault) {
			$.valHooks.select.set = valHooks._set;
			valHooks.useDefault = true;
		}

		this.$element = $(element);
		this.$newElement = null;
		this.$button = null;
		this.$menu = null;
		this.options = options;
		this.selectpicker = {
			main: {},
			current: {}, // current changes if a search is in progress
			search: {},
			view: {},
			keydown: {
				keyHistory: '',
				resetKeyHistory: {
					start: function () {
						return setTimeout(function () {
							that.selectpicker.keydown.keyHistory = '';
						}, 800);
					}
				}
			}
		};
		// If we have no title yet, try to pull it from the html title attribute (jQuery doesnt' pick it up as it's not a
		// data-attribute)
		if (this.options.title === null) {
			this.options.title = this.$element.attr('title');
		}

		// Format window padding
		var winPad = this.options.windowPadding;
		if (typeof winPad === 'number') {
			this.options.windowPadding = [winPad, winPad, winPad, winPad];
		}

		// Expose public methods
		this.val = Selectpicker.prototype.val;
		this.render = Selectpicker.prototype.render;
		this.refresh = Selectpicker.prototype.refresh;
		this.setStyle = Selectpicker.prototype.setStyle;
		this.selectAll = Selectpicker.prototype.selectAll;
		this.deselectAll = Selectpicker.prototype.deselectAll;
		this.destroy = Selectpicker.prototype.destroy;
		this.remove = Selectpicker.prototype.remove;
		this.show = Selectpicker.prototype.show;
		this.hide = Selectpicker.prototype.hide;

		this.init();
	};

	Selectpicker.VERSION = '1.13.9';

	// part of this is duplicated in i18n/defaults-en_US.js. Make sure to update both.
	Selectpicker.DEFAULTS = {
		noneSelectedText: 'Nothing selected',
		noneResultsText: 'No results matched {0}',
		countSelectedText: function (numSelected, numTotal) {
			return (numSelected == 1) ? '{0} item selected' : '{0} items selected';
		},
		maxOptionsText: function (numAll, numGroup) {
			return [
				(numAll == 1) ? 'Limit reached ({n} item max)' : 'Limit reached ({n} items max)',
				(numGroup == 1) ? 'Group limit reached ({n} item max)' : 'Group limit reached ({n} items max)'
			];
		},
		selectAllText: 'Select All',
		deselectAllText: 'Deselect All',
		doneButton: false,
		doneButtonText: 'Close',
		multipleSeparator: ', ',
		styleBase: 'btn',
		style: classNames.BUTTONCLASS,
		size: 'auto',
		title: null,
		selectedTextFormat: 'values',
		width: false,
		container: false,
		hideDisabled: false,
		showSubtext: false,
		showIcon: true,
		showContent: true,
		dropupAuto: true,
		header: false,
		liveSearch: false,
		liveSearchPlaceholder: null,
		liveSearchNormalize: false,
		liveSearchStyle: 'contains',
		actionsBox: false,
		iconBase: classNames.ICONBASE,
		tickIcon: classNames.TICKICON,
		showTick: false,
		template: {
			caret: '<span class="caret"></span>'
		},
		maxOptions: false,
		mobile: false,
		selectOnTab: false,
		dropdownAlignRight: false,
		windowPadding: 0,
		virtualScroll: 600,
		display: false,
		sanitize: true,
		sanitizeFn: null,
		whiteList: DefaultWhitelist
	};

	Selectpicker.prototype = {

		constructor: Selectpicker,

		init: function () {
			var that = this,
					id = this.$element.attr('id');

			this.selectId = selectId++;

			this.$element[0].classList.add('bs-select-hidden');

			this.multiple = this.$element.prop('multiple');
			this.autofocus = this.$element.prop('autofocus');
			this.options.showTick = this.$element[0].classList.contains('show-tick');

			this.$newElement = this.createDropdown();
			this.$element
				.after(this.$newElement)
				.prependTo(this.$newElement);

			this.$button = this.$newElement.children('button');
			this.$menu = this.$newElement.children(Selector.MENU);
			this.$menuInner = this.$menu.children('.inner');
			this.$searchbox = this.$menu.find('input');

			this.$element[0].classList.remove('bs-select-hidden');

			if (this.options.dropdownAlignRight === true) this.$menu[0].classList.add(classNames.MENURIGHT);

			if (typeof id !== 'undefined') {
				this.$button.attr('data-id', id);
			}

			this.checkDisabled();
			this.clickListener();
			if (this.options.liveSearch) this.liveSearchListener();
			this.setStyle();
			this.render();
			this.setWidth();
			if (this.options.container) {
				this.selectPosition();
			} else {
				this.$element.on('hide' + EVENT_KEY, function () {
					if (that.isVirtual()) {
						// empty menu on close
						var menuInner = that.$menuInner[0],
								emptyMenu = menuInner.firstChild.cloneNode(false);

						// replace the existing UL with an empty one - this is faster than $.empty() or innerHTML = ''
						menuInner.replaceChild(emptyMenu, menuInner.firstChild);
						menuInner.scrollTop = 0;
					}
				});
			}
			this.$menu.data('this', this);
			this.$newElement.data('this', this);
			if (this.options.mobile) this.mobile();

			this.$newElement.on({
				'hide.bs.dropdown': function (e) {
					that.$menuInner.attr('aria-expanded', false);
					that.$element.trigger('hide' + EVENT_KEY, e);
				},
				'hidden.bs.dropdown': function (e) {
					that.$element.trigger('hidden' + EVENT_KEY, e);
				},
				'show.bs.dropdown': function (e) {
					that.$menuInner.attr('aria-expanded', true);
					that.$element.trigger('show' + EVENT_KEY, e);
				},
				'shown.bs.dropdown': function (e) {
					that.$element.trigger('shown' + EVENT_KEY, e);
				}
			});

			if (that.$element[0].hasAttribute('required')) {
				this.$element.on('invalid' + EVENT_KEY, function () {
					that.$button[0].classList.add('bs-invalid');

					that.$element
						.on('shown' + EVENT_KEY + '.invalid', function () {
							that.$element
								.val(that.$element.val()) // set the value to hide the validation message in Chrome when menu is opened
								.off('shown' + EVENT_KEY + '.invalid');
						})
						.on('rendered' + EVENT_KEY, function () {
							// if select is no longer invalid, remove the bs-invalid class
							if (this.validity.valid) that.$button[0].classList.remove('bs-invalid');
							that.$element.off('rendered' + EVENT_KEY);
						});

					that.$button.on('blur' + EVENT_KEY, function () {
						that.$element.trigger('focus').trigger('blur');
						that.$button.off('blur' + EVENT_KEY);
					});
				});
			}

			setTimeout(function () {
				that.createLi();
				that.$element.trigger('loaded' + EVENT_KEY);
			});
		},

		createDropdown: function () {
			// Options
			// If we are multiple or showTick option is set, then add the show-tick class
			var showTick = (this.multiple || this.options.showTick) ? ' show-tick' : '',
					inputGroup = '',
					autofocus = this.autofocus ? ' autofocus' : '';

			if (version.major < 4 && this.$element.parent().hasClass('input-group')) {
				inputGroup = ' input-group-btn';
			}

			// Elements
			var drop,
					header = '',
					searchbox = '',
					actionsbox = '',
					donebutton = '';

			if (this.options.header) {
				header =
					'<div class="' + classNames.POPOVERHEADER + '">' +
						'<button type="button" class="close" aria-hidden="true">&times;</button>' +
							this.options.header +
					'</div>';
			}

			if (this.options.liveSearch) {
				searchbox =
					'<div class="bs-searchbox">' +
						'<input type="text" class="form-control" autocomplete="off"' +
							(
								this.options.liveSearchPlaceholder === null ? ''
								:
								' placeholder="' + htmlEscape(this.options.liveSearchPlaceholder) + '"'
							) +
							' role="textbox" aria-label="Search">' +
					'</div>';
			}

			if (this.multiple && this.options.actionsBox) {
				actionsbox =
					'<div class="bs-actionsbox">' +
						'<div class="btn-group btn-group-sm btn-block">' +
							'<button type="button" class="actions-btn bs-select-all btn ' + classNames.BUTTONCLASS + '">' +
								this.options.selectAllText +
							'</button>' +
							'<button type="button" class="actions-btn bs-deselect-all btn ' + classNames.BUTTONCLASS + '">' +
								this.options.deselectAllText +
							'</button>' +
						'</div>' +
					'</div>';
			}

			if (this.multiple && this.options.doneButton) {
				donebutton =
					'<div class="bs-donebutton">' +
						'<div class="btn-group btn-block">' +
							'<button type="button" class="btn btn-sm ' + classNames.BUTTONCLASS + '">' +
								this.options.doneButtonText +
							'</button>' +
						'</div>' +
					'</div>';
			}

			drop =
				'<div class="dropdown bootstrap-select' + showTick + inputGroup + '">' +
					'<button type="button" class="' + this.options.styleBase + ' dropdown-toggle" ' + (this.options.display === 'static' ? 'data-display="static"' : '') + 'data-toggle="dropdown"' + autofocus + ' role="button">' +
						'<div class="filter-option">' +
							'<div class="filter-option-inner">' +
								'<div class="filter-option-inner-inner"></div>' +
							'</div> ' +
						'</div>' +
						(
							version.major === '4' ? ''
							:
							'<span class="bs-caret">' +
								this.options.template.caret +
							'</span>'
						) +
					'</button>' +
					'<div class="' + classNames.MENU + ' ' + (version.major === '4' ? '' : classNames.SHOW) + '" role="combobox">' +
						header +
						searchbox +
						actionsbox +
						'<div class="inner ' + classNames.SHOW + '" role="listbox" aria-expanded="false" tabindex="-1">' +
								'<ul class="' + classNames.MENU + ' inner ' + (version.major === '4' ? classNames.SHOW : '') + '">' +
								'</ul>' +
						'</div>' +
						donebutton +
					'</div>' +
				'</div>';

			return $(drop);
		},

		setPositionData: function () {
			this.selectpicker.view.canHighlight = [];

			for (var i = 0; i < this.selectpicker.current.data.length; i++) {
				var li = this.selectpicker.current.data[i],
						canHighlight = true;

				if (li.type === 'divider') {
					canHighlight = false;
					li.height = this.sizeInfo.dividerHeight;
				} else if (li.type === 'optgroup-label') {
					canHighlight = false;
					li.height = this.sizeInfo.dropdownHeaderHeight;
				} else {
					li.height = this.sizeInfo.liHeight;
				}

				if (li.disabled) canHighlight = false;

				this.selectpicker.view.canHighlight.push(canHighlight);

				li.position = (i === 0 ? 0 : this.selectpicker.current.data[i - 1].position) + li.height;
			}
		},

		isVirtual: function () {
			return (this.options.virtualScroll !== false) && (this.selectpicker.main.elements.length >= this.options.virtualScroll) || this.options.virtualScroll === true;
		},

		createView: function (isSearching, scrollTop) {
			scrollTop = scrollTop || 0;

			var that = this;

			this.selectpicker.current = isSearching ? this.selectpicker.search : this.selectpicker.main;

			var active = [];
			var selected;
			var prevActive;

			this.setPositionData();

			scroll(scrollTop, true);

			this.$menuInner.off('scroll.createView').on('scroll.createView', function (e, updateValue) {
				if (!that.noScroll) scroll(this.scrollTop, updateValue);
				that.noScroll = false;
			});

			function scroll (scrollTop, init) {
				var size = that.selectpicker.current.elements.length,
						chunks = [],
						chunkSize,
						chunkCount,
						firstChunk,
						lastChunk,
						currentChunk,
						prevPositions,
						positionIsDifferent,
						previousElements,
						menuIsDifferent = true,
						isVirtual = that.isVirtual();

				that.selectpicker.view.scrollTop = scrollTop;

				if (isVirtual === true) {
					// if an option that is encountered that is wider than the current menu width, update the menu width accordingly
					if (that.sizeInfo.hasScrollBar && that.$menu[0].offsetWidth > that.sizeInfo.totalMenuWidth) {
						that.sizeInfo.menuWidth = that.$menu[0].offsetWidth;
						that.sizeInfo.totalMenuWidth = that.sizeInfo.menuWidth + that.sizeInfo.scrollBarWidth;
						that.$menu.css('min-width', that.sizeInfo.menuWidth);
					}
				}

				chunkSize = Math.ceil(that.sizeInfo.menuInnerHeight / that.sizeInfo.liHeight * 1.5); // number of options in a chunk
				chunkCount = Math.round(size / chunkSize) || 1; // number of chunks

				for (var i = 0; i < chunkCount; i++) {
					var endOfChunk = (i + 1) * chunkSize;

					if (i === chunkCount - 1) {
						endOfChunk = size;
					}

					chunks[i] = [
						(i) * chunkSize + (!i ? 0 : 1),
						endOfChunk
					];

					if (!size) break;

					if (currentChunk === undefined && scrollTop <= that.selectpicker.current.data[endOfChunk - 1].position - that.sizeInfo.menuInnerHeight) {
						currentChunk = i;
					}
				}

				if (currentChunk === undefined) currentChunk = 0;

				prevPositions = [that.selectpicker.view.position0, that.selectpicker.view.position1];

				// always display previous, current, and next chunks
				firstChunk = Math.max(0, currentChunk - 1);
				lastChunk = Math.min(chunkCount - 1, currentChunk + 1);

				that.selectpicker.view.position0 = isVirtual === false ? 0 : (Math.max(0, chunks[firstChunk][0]) || 0);
				that.selectpicker.view.position1 = isVirtual === false ? size : (Math.min(size, chunks[lastChunk][1]) || 0);

				positionIsDifferent = prevPositions[0] !== that.selectpicker.view.position0 || prevPositions[1] !== that.selectpicker.view.position1;

				if (that.activeIndex !== undefined) {
					prevActive = that.selectpicker.main.elements[that.prevActiveIndex];
					active = that.selectpicker.main.elements[that.activeIndex];
					selected = that.selectpicker.main.elements[that.selectedIndex];

					if (init) {
						if (that.activeIndex !== that.selectedIndex && active && active.length) {
							active.classList.remove('active');
							if (active.firstChild) active.firstChild.classList.remove('active');
						}
						that.activeIndex = undefined;
					}

					if (that.activeIndex && that.activeIndex !== that.selectedIndex && selected && selected.length) {
						selected.classList.remove('active');
						if (selected.firstChild) selected.firstChild.classList.remove('active');
					}
				}

				if (that.prevActiveIndex !== undefined && that.prevActiveIndex !== that.activeIndex && that.prevActiveIndex !== that.selectedIndex && prevActive && prevActive.length) {
					prevActive.classList.remove('active');
					if (prevActive.firstChild) prevActive.firstChild.classList.remove('active');
				}

				if (init || positionIsDifferent) {
					previousElements = that.selectpicker.view.visibleElements ? that.selectpicker.view.visibleElements.slice() : [];

					if (isVirtual === false) {
						that.selectpicker.view.visibleElements = that.selectpicker.current.elements;
					} else {
						that.selectpicker.view.visibleElements = that.selectpicker.current.elements.slice(that.selectpicker.view.position0, that.selectpicker.view.position1);
					}

					that.setOptionStatus();

					// if searching, check to make sure the list has actually been updated before updating DOM
					// this prevents unnecessary repaints
					if (isSearching || (isVirtual === false && init)) menuIsDifferent = !isEqual(previousElements, that.selectpicker.view.visibleElements);

					// if virtual scroll is disabled and not searching,
					// menu should never need to be updated more than once
					if ((init || isVirtual === true) && menuIsDifferent) {
						var menuInner = that.$menuInner[0],
								menuFragment = document.createDocumentFragment(),
								emptyMenu = menuInner.firstChild.cloneNode(false),
								marginTop,
								marginBottom,
								elements = that.selectpicker.view.visibleElements,
								toSanitize = [];

						// replace the existing UL with an empty one - this is faster than $.empty()
						menuInner.replaceChild(emptyMenu, menuInner.firstChild);

						for (var i = 0, visibleElementsLen = elements.length; i < visibleElementsLen; i++) {
							var element = elements[i],
									elText,
									elementData;

							if (that.options.sanitize) {
								elText = element.lastChild;

								if (elText) {
									elementData = that.selectpicker.current.data[i + that.selectpicker.view.position0];

									if (elementData && elementData.content && !elementData.sanitized) {
										toSanitize.push(elText);
										elementData.sanitized = true;
									}
								}
							}

							menuFragment.appendChild(element);
						}

						if (that.options.sanitize && toSanitize.length) {
							sanitizeHtml(toSanitize, that.options.whiteList, that.options.sanitizeFn);
						}

						if (isVirtual === true) {
							marginTop = (that.selectpicker.view.position0 === 0 ? 0 : that.selectpicker.current.data[that.selectpicker.view.position0 - 1].position);
							marginBottom = (that.selectpicker.view.position1 > size - 1 ? 0 : that.selectpicker.current.data[size - 1].position - that.selectpicker.current.data[that.selectpicker.view.position1 - 1].position);

							menuInner.firstChild.style.marginTop = marginTop + 'px';
							menuInner.firstChild.style.marginBottom = marginBottom + 'px';
						}

						menuInner.firstChild.appendChild(menuFragment);
					}
				}

				that.prevActiveIndex = that.activeIndex;

				if (!that.options.liveSearch) {
					that.$menuInner.trigger('focus');
				} else if (isSearching && init) {
					var index = 0,
							newActive;

					if (!that.selectpicker.view.canHighlight[index]) {
						index = 1 + that.selectpicker.view.canHighlight.slice(1).indexOf(true);
					}

					newActive = that.selectpicker.view.visibleElements[index];

					if (that.selectpicker.view.currentActive) {
						that.selectpicker.view.currentActive.classList.remove('active');
						if (that.selectpicker.view.currentActive.firstChild) that.selectpicker.view.currentActive.firstChild.classList.remove('active');
					}

					if (newActive) {
						newActive.classList.add('active');
						if (newActive.firstChild) newActive.firstChild.classList.add('active');
					}

					that.activeIndex = (that.selectpicker.current.data[index] || {}).index;
				}
			}

			$(window)
				.off('resize' + EVENT_KEY + '.' + this.selectId + '.createView')
				.on('resize' + EVENT_KEY + '.' + this.selectId + '.createView', function () {
					var isActive = that.$newElement.hasClass(classNames.SHOW);

					if (isActive) scroll(that.$menuInner[0].scrollTop);
				});
		},

		setPlaceholder: function () {
			var updateIndex = false;

			if (this.options.title && !this.multiple) {
				if (!this.selectpicker.view.titleOption) this.selectpicker.view.titleOption = document.createElement('option');

				// this option doesn't create a new <li> element, but does add a new option at the start,
				// so startIndex should increase to prevent having to check every option for the bs-title-option class
				updateIndex = true;

				var element = this.$element[0],
						isSelected = false,
						titleNotAppended = !this.selectpicker.view.titleOption.parentNode;

				if (titleNotAppended) {
					// Use native JS to prepend option (faster)
					this.selectpicker.view.titleOption.className = 'bs-title-option';
					this.selectpicker.view.titleOption.value = '';

					// Check if selected or data-selected attribute is already set on an option. If not, select the titleOption option.
					// the selected item may have been changed by user or programmatically before the bootstrap select plugin runs,
					// if so, the select will have the data-selected attribute
					var $opt = $(element.options[element.selectedIndex]);
					isSelected = $opt.attr('selected') === undefined && this.$element.data('selected') === undefined;
				}

				if (titleNotAppended || this.selectpicker.view.titleOption.index !== 0) {
					element.insertBefore(this.selectpicker.view.titleOption, element.firstChild);
				}

				// Set selected *after* appending to select,
				// otherwise the option doesn't get selected in IE
				// set using selectedIndex, as setting the selected attr to true here doesn't work in IE11
				if (isSelected) element.selectedIndex = 0;
			}

			return updateIndex;
		},

		createLi: function () {
			var that = this,
					iconBase = this.options.iconBase,
					optionSelector = ':not([hidden]):not([data-hidden="true"])',
					mainElements = [],
					mainData = [],
					widestOptionLength = 0,
					optID = 0,
					startIndex = this.setPlaceholder() ? 1 : 0; // append the titleOption if necessary and skip the first option in the loop

			if (this.options.hideDisabled) optionSelector += ':not(:disabled)';

			if ((that.options.showTick || that.multiple) && !elementTemplates.checkMark.parentNode) {
				elementTemplates.checkMark.className = iconBase + ' ' + that.options.tickIcon + ' check-mark';
				elementTemplates.a.appendChild(elementTemplates.checkMark);
			}

			var selectOptions = this.$element[0].querySelectorAll('select > *' + optionSelector);

			function addDivider (config) {
				var previousData = mainData[mainData.length - 1];

				// ensure optgroup doesn't create back-to-back dividers
				if (
					previousData &&
					previousData.type === 'divider' &&
					(previousData.optID || config.optID)
				) {
					return;
				}

				config = config || {};
				config.type = 'divider';

				mainElements.push(
					generateOption.li(
						false,
						classNames.DIVIDER,
						(config.optID ? config.optID + 'div' : undefined)
					)
				);

				mainData.push(config);
			}

			function addOption (option, config) {
				config = config || {};

				config.divider = option.getAttribute('data-divider') === 'true';

				if (config.divider) {
					addDivider({
						optID: config.optID
					});
				} else {
					var liIndex = mainData.length,
							cssText = option.style.cssText,
							inlineStyle = cssText ? htmlEscape(cssText) : '',
							optionClass = (option.className || '') + (config.optgroupClass || '');

					if (config.optID) optionClass = 'opt ' + optionClass;

					config.text = option.textContent;

					config.content = option.getAttribute('data-content');
					config.tokens = option.getAttribute('data-tokens');
					config.subtext = option.getAttribute('data-subtext');
					config.icon = option.getAttribute('data-icon');
					config.iconBase = iconBase;

					var textElement = generateOption.text(config);

					mainElements.push(
						generateOption.li(
							generateOption.a(
								textElement,
								optionClass,
								inlineStyle
							),
							'',
							config.optID
						)
					);

					option.liIndex = liIndex;

					config.display = config.content || config.text;
					config.type = 'option';
					config.index = liIndex;
					config.option = option;
					config.disabled = config.disabled || option.disabled;

					mainData.push(config);

					var combinedLength = 0;

					// count the number of characters in the option - not perfect, but should work in most cases
					if (config.display) combinedLength += config.display.length;
					if (config.subtext) combinedLength += config.subtext.length;
					// if there is an icon, ensure this option's width is checked
					if (config.icon) combinedLength += 1;

					if (combinedLength > widestOptionLength) {
						widestOptionLength = combinedLength;

						// guess which option is the widest
						// use this when calculating menu width
						// not perfect, but it's fast, and the width will be updating accordingly when scrolling
						that.selectpicker.view.widestOption = mainElements[mainElements.length - 1];
					}
				}
			}

			function addOptgroup (index, selectOptions) {
				var optgroup = selectOptions[index],
						previous = selectOptions[index - 1],
						next = selectOptions[index + 1],
						options = optgroup.querySelectorAll('option' + optionSelector);

				if (!options.length) return;

				var config = {
							label: htmlEscape(optgroup.label),
							subtext: optgroup.getAttribute('data-subtext'),
							icon: optgroup.getAttribute('data-icon'),
							iconBase: iconBase
						},
						optgroupClass = ' ' + (optgroup.className || ''),
						headerIndex,
						lastIndex;

				optID++;

				if (previous) {
					addDivider({ optID: optID });
				}

				var labelElement = generateOption.label(config);

				mainElements.push(
					generateOption.li(labelElement, 'dropdown-header' + optgroupClass, optID)
				);

				mainData.push({
					display: config.label,
					subtext: config.subtext,
					type: 'optgroup-label',
					optID: optID
				});

				for (var j = 0, len = options.length; j < len; j++) {
					var option = options[j];

					if (j === 0) {
						headerIndex = mainData.length - 1;
						lastIndex = headerIndex + len;
					}

					addOption(option, {
						headerIndex: headerIndex,
						lastIndex: lastIndex,
						optID: optID,
						optgroupClass: optgroupClass,
						disabled: optgroup.disabled
					});
				}

				if (next) {
					addDivider({ optID: optID });
				}
			}

			for (var len = selectOptions.length; startIndex < len; startIndex++) {
				var item = selectOptions[startIndex];

				if (item.tagName !== 'OPTGROUP') {
					addOption(item, {});
				} else {
					addOptgroup(startIndex, selectOptions);
				}
			}

			this.selectpicker.main.elements = mainElements;
			this.selectpicker.main.data = mainData;

			this.selectpicker.current = this.selectpicker.main;
		},

		findLis: function () {
			return this.$menuInner.find('.inner > li');
		},

		render: function () {
			// ensure titleOption is appended and selected (if necessary) before getting selectedOptions
			this.setPlaceholder();

			var that = this,
					selectedOptions = this.$element[0].selectedOptions,
					selectedCount = selectedOptions.length,
					button = this.$button[0],
					buttonInner = button.querySelector('.filter-option-inner-inner'),
					multipleSeparator = document.createTextNode(this.options.multipleSeparator),
					titleFragment = elementTemplates.fragment.cloneNode(false),
					showCount,
					countMax,
					hasContent = false;

			this.togglePlaceholder();

			this.tabIndex();

			if (this.options.selectedTextFormat === 'static') {
				titleFragment = generateOption.text({ text: this.options.title }, true);
			} else {
				showCount = this.multiple && this.options.selectedTextFormat.indexOf('count') !== -1 && selectedCount > 1;

				// determine if the number of selected options will be shown (showCount === true)
				if (showCount) {
					countMax = this.options.selectedTextFormat.split('>');
					showCount = (countMax.length > 1 && selectedCount > countMax[1]) || (countMax.length === 1 && selectedCount >= 2);
				}

				// only loop through all selected options if the count won't be shown
				if (showCount === false) {
					for (var selectedIndex = 0; selectedIndex < selectedCount; selectedIndex++) {
						if (selectedIndex < 50) {
							var option = selectedOptions[selectedIndex],
									titleOptions = {},
									thisData = {
										content: option.getAttribute('data-content'),
										subtext: option.getAttribute('data-subtext'),
										icon: option.getAttribute('data-icon')
									};

							if (this.multiple && selectedIndex > 0) {
								titleFragment.appendChild(multipleSeparator.cloneNode(false));
							}

							if (option.title) {
								titleOptions.text = option.title;
							} else if (thisData.content && that.options.showContent) {
								titleOptions.content = thisData.content.toString();
								hasContent = true;
							} else {
								if (that.options.showIcon) {
									titleOptions.icon = thisData.icon;
									titleOptions.iconBase = this.options.iconBase;
								}
								if (that.options.showSubtext && !that.multiple && thisData.subtext) titleOptions.subtext = ' ' + thisData.subtext;
								titleOptions.text = option.textContent.trim();
							}

							titleFragment.appendChild(generateOption.text(titleOptions, true));
						} else {
							break;
						}
					}

					// add ellipsis
					if (selectedCount > 49) {
						titleFragment.appendChild(document.createTextNode('...'));
					}
				} else {
					var optionSelector = ':not([hidden]):not([data-hidden="true"]):not([data-divider="true"])';
					if (this.options.hideDisabled) optionSelector += ':not(:disabled)';

					// If this is a multiselect, and selectedTextFormat is count, then show 1 of 2 selected, etc.
					var totalCount = this.$element[0].querySelectorAll('select > option' + optionSelector + ', optgroup' + optionSelector + ' option' + optionSelector).length,
							tr8nText = (typeof this.options.countSelectedText === 'function') ? this.options.countSelectedText(selectedCount, totalCount) : this.options.countSelectedText;

					titleFragment = generateOption.text({
						text: tr8nText.replace('{0}', selectedCount.toString()).replace('{1}', totalCount.toString())
					}, true);
				}
			}

			if (this.options.title == undefined) {
				// use .attr to ensure undefined is returned if title attribute is not set
				this.options.title = this.$element.attr('title');
			}

			// If the select doesn't have a title, then use the default, or if nothing is set at all, use noneSelectedText
			if (!titleFragment.childNodes.length) {
				titleFragment = generateOption.text({
					text: typeof this.options.title !== 'undefined' ? this.options.title : this.options.noneSelectedText
				}, true);
			}

			// strip all HTML tags and trim the result, then unescape any escaped tags
			button.title = titleFragment.textContent.replace(/<[^>]*>?/g, '').trim();

			if (this.options.sanitize && hasContent) {
				sanitizeHtml([titleFragment], that.options.whiteList, that.options.sanitizeFn);
			}

			buttonInner.innerHTML = '';
			buttonInner.appendChild(titleFragment);

			if (version.major < 4 && this.$newElement[0].classList.contains('bs3-has-addon')) {
				var filterExpand = button.querySelector('.filter-expand'),
						clone = buttonInner.cloneNode(true);

				clone.className = 'filter-expand';

				if (filterExpand) {
					button.replaceChild(clone, filterExpand);
				} else {
					button.appendChild(clone);
				}
			}

			this.$element.trigger('rendered' + EVENT_KEY);
		},

		/**
		 * @param [style]
		 * @param [status]
		 */
		setStyle: function (newStyle, status) {
			var button = this.$button[0],
					newElement = this.$newElement[0],
					style = this.options.style.trim(),
					buttonClass;

			if (this.$element.attr('class')) {
				this.$newElement.addClass(this.$element.attr('class').replace(/selectpicker|mobile-device|bs-select-hidden|validate\[.*\]/gi, ''));
			}

			if (version.major < 4) {
				newElement.classList.add('bs3');

				if (newElement.parentNode.classList.contains('input-group') &&
						(newElement.previousElementSibling || newElement.nextElementSibling) &&
						(newElement.previousElementSibling || newElement.nextElementSibling).classList.contains('input-group-addon')
				) {
					newElement.classList.add('bs3-has-addon');
				}
			}

			if (newStyle) {
				buttonClass = newStyle.trim();
			} else {
				buttonClass = style;
			}

			if (status == 'add') {
				if (buttonClass) button.classList.add.apply(button.classList, buttonClass.split(' '));
			} else if (status == 'remove') {
				if (buttonClass) button.classList.remove.apply(button.classList, buttonClass.split(' '));
			} else {
				if (style) button.classList.remove.apply(button.classList, style.split(' '));
				if (buttonClass) button.classList.add.apply(button.classList, buttonClass.split(' '));
			}
		},

		liHeight: function (refresh) {
			if (!refresh && (this.options.size === false || this.sizeInfo)) return;

			if (!this.sizeInfo) this.sizeInfo = {};

			var newElement = document.createElement('div'),
					menu = document.createElement('div'),
					menuInner = document.createElement('div'),
					menuInnerInner = document.createElement('ul'),
					divider = document.createElement('li'),
					dropdownHeader = document.createElement('li'),
					li = document.createElement('li'),
					a = document.createElement('a'),
					text = document.createElement('span'),
					header = this.options.header && this.$menu.find('.' + classNames.POPOVERHEADER).length > 0 ? this.$menu.find('.' + classNames.POPOVERHEADER)[0].cloneNode(true) : null,
					search = this.options.liveSearch ? document.createElement('div') : null,
					actions = this.options.actionsBox && this.multiple && this.$menu.find('.bs-actionsbox').length > 0 ? this.$menu.find('.bs-actionsbox')[0].cloneNode(true) : null,
					doneButton = this.options.doneButton && this.multiple && this.$menu.find('.bs-donebutton').length > 0 ? this.$menu.find('.bs-donebutton')[0].cloneNode(true) : null,
					firstOption = this.$element.find('option')[0];

			this.sizeInfo.selectWidth = this.$newElement[0].offsetWidth;

			text.className = 'text';
			a.className = 'dropdown-item ' + (firstOption ? firstOption.className : '');
			newElement.className = this.$menu[0].parentNode.className + ' ' + classNames.SHOW;
			newElement.style.width = this.sizeInfo.selectWidth + 'px';
			if (this.options.width === 'auto') menu.style.minWidth = 0;
			menu.className = classNames.MENU + ' ' + classNames.SHOW;
			menuInner.className = 'inner ' + classNames.SHOW;
			menuInnerInner.className = classNames.MENU + ' inner ' + (version.major === '4' ? classNames.SHOW : '');
			divider.className = classNames.DIVIDER;
			dropdownHeader.className = 'dropdown-header';

			text.appendChild(document.createTextNode('\u200b'));
			a.appendChild(text);
			li.appendChild(a);
			dropdownHeader.appendChild(text.cloneNode(true));

			if (this.selectpicker.view.widestOption) {
				menuInnerInner.appendChild(this.selectpicker.view.widestOption.cloneNode(true));
			}

			menuInnerInner.appendChild(li);
			menuInnerInner.appendChild(divider);
			menuInnerInner.appendChild(dropdownHeader);
			if (header) menu.appendChild(header);
			if (search) {
				var input = document.createElement('input');
				search.className = 'bs-searchbox';
				input.className = 'form-control';
				search.appendChild(input);
				menu.appendChild(search);
			}
			if (actions) menu.appendChild(actions);
			menuInner.appendChild(menuInnerInner);
			menu.appendChild(menuInner);
			if (doneButton) menu.appendChild(doneButton);
			newElement.appendChild(menu);

			document.body.appendChild(newElement);

			var liHeight = li.offsetHeight,
					dropdownHeaderHeight = dropdownHeader ? dropdownHeader.offsetHeight : 0,
					headerHeight = header ? header.offsetHeight : 0,
					searchHeight = search ? search.offsetHeight : 0,
					actionsHeight = actions ? actions.offsetHeight : 0,
					doneButtonHeight = doneButton ? doneButton.offsetHeight : 0,
					dividerHeight = $(divider).outerHeight(true),
					// fall back to jQuery if getComputedStyle is not supported
					menuStyle = window.getComputedStyle ? window.getComputedStyle(menu) : false,
					menuWidth = menu.offsetWidth,
					$menu = menuStyle ? null : $(menu),
					menuPadding = {
						vert: toInteger(menuStyle ? menuStyle.paddingTop : $menu.css('paddingTop')) +
									toInteger(menuStyle ? menuStyle.paddingBottom : $menu.css('paddingBottom')) +
									toInteger(menuStyle ? menuStyle.borderTopWidth : $menu.css('borderTopWidth')) +
									toInteger(menuStyle ? menuStyle.borderBottomWidth : $menu.css('borderBottomWidth')),
						horiz: toInteger(menuStyle ? menuStyle.paddingLeft : $menu.css('paddingLeft')) +
									toInteger(menuStyle ? menuStyle.paddingRight : $menu.css('paddingRight')) +
									toInteger(menuStyle ? menuStyle.borderLeftWidth : $menu.css('borderLeftWidth')) +
									toInteger(menuStyle ? menuStyle.borderRightWidth : $menu.css('borderRightWidth'))
					},
					menuExtras = {
						vert: menuPadding.vert +
									toInteger(menuStyle ? menuStyle.marginTop : $menu.css('marginTop')) +
									toInteger(menuStyle ? menuStyle.marginBottom : $menu.css('marginBottom')) + 2,
						horiz: menuPadding.horiz +
									toInteger(menuStyle ? menuStyle.marginLeft : $menu.css('marginLeft')) +
									toInteger(menuStyle ? menuStyle.marginRight : $menu.css('marginRight')) + 2
					},
					scrollBarWidth;

			menuInner.style.overflowY = 'scroll';

			scrollBarWidth = menu.offsetWidth - menuWidth;

			document.body.removeChild(newElement);

			this.sizeInfo.liHeight = liHeight;
			this.sizeInfo.dropdownHeaderHeight = dropdownHeaderHeight;
			this.sizeInfo.headerHeight = headerHeight;
			this.sizeInfo.searchHeight = searchHeight;
			this.sizeInfo.actionsHeight = actionsHeight;
			this.sizeInfo.doneButtonHeight = doneButtonHeight;
			this.sizeInfo.dividerHeight = dividerHeight;
			this.sizeInfo.menuPadding = menuPadding;
			this.sizeInfo.menuExtras = menuExtras;
			this.sizeInfo.menuWidth = menuWidth;
			this.sizeInfo.totalMenuWidth = this.sizeInfo.menuWidth;
			this.sizeInfo.scrollBarWidth = scrollBarWidth;
			this.sizeInfo.selectHeight = this.$newElement[0].offsetHeight;

			this.setPositionData();
		},

		getSelectPosition: function () {
			var that = this,
					$window = $(window),
					pos = that.$newElement.offset(),
					$container = $(that.options.container),
					containerPos;

			if (that.options.container && $container.length && !$container.is('body')) {
				containerPos = $container.offset();
				containerPos.top += parseInt($container.css('borderTopWidth'));
				containerPos.left += parseInt($container.css('borderLeftWidth'));
			} else {
				containerPos = { top: 0, left: 0 };
			}

			var winPad = that.options.windowPadding;

			this.sizeInfo.selectOffsetTop = pos.top - containerPos.top - $window.scrollTop();
			this.sizeInfo.selectOffsetBot = $window.height() - this.sizeInfo.selectOffsetTop - this.sizeInfo.selectHeight - containerPos.top - winPad[2];
			this.sizeInfo.selectOffsetLeft = pos.left - containerPos.left - $window.scrollLeft();
			this.sizeInfo.selectOffsetRight = $window.width() - this.sizeInfo.selectOffsetLeft - this.sizeInfo.selectWidth - containerPos.left - winPad[1];
			this.sizeInfo.selectOffsetTop -= winPad[0];
			this.sizeInfo.selectOffsetLeft -= winPad[3];
		},

		setMenuSize: function (isAuto) {
			this.getSelectPosition();

			var selectWidth = this.sizeInfo.selectWidth,
					liHeight = this.sizeInfo.liHeight,
					headerHeight = this.sizeInfo.headerHeight,
					searchHeight = this.sizeInfo.searchHeight,
					actionsHeight = this.sizeInfo.actionsHeight,
					doneButtonHeight = this.sizeInfo.doneButtonHeight,
					divHeight = this.sizeInfo.dividerHeight,
					menuPadding = this.sizeInfo.menuPadding,
					menuInnerHeight,
					menuHeight,
					divLength = 0,
					minHeight,
					_minHeight,
					maxHeight,
					menuInnerMinHeight,
					estimate;

			if (this.options.dropupAuto) {
				// Get the estimated height of the menu without scrollbars.
				// This is useful for smaller menus, where there might be plenty of room
				// below the button without setting dropup, but we can't know
				// the exact height of the menu until createView is called later
				estimate = liHeight * this.selectpicker.current.elements.length + menuPadding.vert;
				this.$newElement.toggleClass(classNames.DROPUP, this.sizeInfo.selectOffsetTop - this.sizeInfo.selectOffsetBot > this.sizeInfo.menuExtras.vert && estimate + this.sizeInfo.menuExtras.vert + 50 > this.sizeInfo.selectOffsetBot);
			}

			if (this.options.size === 'auto') {
				_minHeight = this.selectpicker.current.elements.length > 3 ? this.sizeInfo.liHeight * 3 + this.sizeInfo.menuExtras.vert - 2 : 0;
				menuHeight = this.sizeInfo.selectOffsetBot - this.sizeInfo.menuExtras.vert;
				minHeight = _minHeight + headerHeight + searchHeight + actionsHeight + doneButtonHeight;
				menuInnerMinHeight = Math.max(_minHeight - menuPadding.vert, 0);

				if (this.$newElement.hasClass(classNames.DROPUP)) {
					menuHeight = this.sizeInfo.selectOffsetTop - this.sizeInfo.menuExtras.vert;
				}

				maxHeight = menuHeight;
				menuInnerHeight = menuHeight - headerHeight - searchHeight - actionsHeight - doneButtonHeight - menuPadding.vert;
			} else if (this.options.size && this.options.size != 'auto' && this.selectpicker.current.elements.length > this.options.size) {
				for (var i = 0; i < this.options.size; i++) {
					if (this.selectpicker.current.data[i].type === 'divider') divLength++;
				}

				menuHeight = liHeight * this.options.size + divLength * divHeight + menuPadding.vert;
				menuInnerHeight = menuHeight - menuPadding.vert;
				maxHeight = menuHeight + headerHeight + searchHeight + actionsHeight + doneButtonHeight;
				minHeight = menuInnerMinHeight = '';
			}

			if (this.options.dropdownAlignRight === 'auto') {
				this.$menu.toggleClass(classNames.MENURIGHT, this.sizeInfo.selectOffsetLeft > this.sizeInfo.selectOffsetRight && this.sizeInfo.selectOffsetRight < (this.sizeInfo.totalMenuWidth - selectWidth));
			}

			this.$menu.css({
				'max-height': maxHeight + 'px',
				'overflow': 'hidden',
				'min-height': minHeight + 'px'
			});

			this.$menuInner.css({
				'max-height': menuInnerHeight + 'px',
				'overflow-y': 'auto',
				'min-height': menuInnerMinHeight + 'px'
			});

			// ensure menuInnerHeight is always a positive number to prevent issues calculating chunkSize in createView
			this.sizeInfo.menuInnerHeight = Math.max(menuInnerHeight, 1);

			if (this.selectpicker.current.data.length && this.selectpicker.current.data[this.selectpicker.current.data.length - 1].position > this.sizeInfo.menuInnerHeight) {
				this.sizeInfo.hasScrollBar = true;
				this.sizeInfo.totalMenuWidth = this.sizeInfo.menuWidth + this.sizeInfo.scrollBarWidth;

				this.$menu.css('min-width', this.sizeInfo.totalMenuWidth);
			}

			if (this.dropdown && this.dropdown._popper) this.dropdown._popper.update();
		},

		setSize: function (refresh) {
			this.liHeight(refresh);

			if (this.options.header) this.$menu.css('padding-top', 0);
			if (this.options.size === false) return;

			var that = this,
					$window = $(window),
					selectedIndex,
					offset = 0;

			this.setMenuSize();

			if (this.options.liveSearch) {
				this.$searchbox
					.off('input.setMenuSize propertychange.setMenuSize')
					.on('input.setMenuSize propertychange.setMenuSize', function () {
						return that.setMenuSize();
					});
			}

			if (this.options.size === 'auto') {
				$window
					.off('resize' + EVENT_KEY + '.' + this.selectId + '.setMenuSize' + ' scroll' + EVENT_KEY + '.' + this.selectId + '.setMenuSize')
					.on('resize' + EVENT_KEY + '.' + this.selectId + '.setMenuSize' + ' scroll' + EVENT_KEY + '.' + this.selectId + '.setMenuSize', function () {
						return that.setMenuSize();
					});
			} else if (this.options.size && this.options.size != 'auto' && this.selectpicker.current.elements.length > this.options.size) {
				$window.off('resize' + EVENT_KEY + '.' + this.selectId + '.setMenuSize' + ' scroll' + EVENT_KEY + '.' + this.selectId + '.setMenuSize');
			}

			if (refresh) {
				offset = this.$menuInner[0].scrollTop;
			} else if (!that.multiple) {
				var element = that.$element[0];
				selectedIndex = (element.options[element.selectedIndex] || {}).liIndex;

				if (typeof selectedIndex === 'number' && that.options.size !== false) {
					offset = that.sizeInfo.liHeight * selectedIndex;
					offset = offset - (that.sizeInfo.menuInnerHeight / 2) + (that.sizeInfo.liHeight / 2);
				}
			}

			that.createView(false, offset);
		},

		setWidth: function () {
			var that = this;

			if (this.options.width === 'auto') {
				requestAnimationFrame(function () {
					that.$menu.css('min-width', '0');

					that.$element.on('loaded' + EVENT_KEY, function () {
						that.liHeight();
						that.setMenuSize();

						// Get correct width if element is hidden
						var $selectClone = that.$newElement.clone().appendTo('body'),
								btnWidth = $selectClone.css('width', 'auto').children('button').outerWidth();

						$selectClone.remove();

						// Set width to whatever's larger, button title or longest option
						that.sizeInfo.selectWidth = Math.max(that.sizeInfo.totalMenuWidth, btnWidth);
						that.$newElement.css('width', that.sizeInfo.selectWidth + 'px');
					});
				});
			} else if (this.options.width === 'fit') {
				// Remove inline min-width so width can be changed from 'auto'
				this.$menu.css('min-width', '');
				this.$newElement.css('width', '').addClass('fit-width');
			} else if (this.options.width) {
				// Remove inline min-width so width can be changed from 'auto'
				this.$menu.css('min-width', '');
				this.$newElement.css('width', this.options.width);
			} else {
				// Remove inline min-width/width so width can be changed
				this.$menu.css('min-width', '');
				this.$newElement.css('width', '');
			}
			// Remove fit-width class if width is changed programmatically
			if (this.$newElement.hasClass('fit-width') && this.options.width !== 'fit') {
				this.$newElement[0].classList.remove('fit-width');
			}
		},

		selectPosition: function () {
			this.$bsContainer = $('<div class="bs-container" />');

			var that = this,
					$container = $(this.options.container),
					pos,
					containerPos,
					actualHeight,
					getPlacement = function ($element) {
						var containerPosition = {},
								// fall back to dropdown's default display setting if display is not manually set
								display = that.options.display || (
									// Bootstrap 3 doesn't have $.fn.dropdown.Constructor.Default
									$.fn.dropdown.Constructor.Default ? $.fn.dropdown.Constructor.Default.display
									: false
								);

						that.$bsContainer.addClass($element.attr('class').replace(/form-control|fit-width/gi, '')).toggleClass(classNames.DROPUP, $element.hasClass(classNames.DROPUP));
						pos = $element.offset();

						if (!$container.is('body')) {
							containerPos = $container.offset();
							containerPos.top += parseInt($container.css('borderTopWidth')) - $container.scrollTop();
							containerPos.left += parseInt($container.css('borderLeftWidth')) - $container.scrollLeft();
						} else {
							containerPos = { top: 0, left: 0 };
						}

						actualHeight = $element.hasClass(classNames.DROPUP) ? 0 : $element[0].offsetHeight;

						// Bootstrap 4+ uses Popper for menu positioning
						if (version.major < 4 || display === 'static') {
							containerPosition.top = pos.top - containerPos.top + actualHeight;
							containerPosition.left = pos.left - containerPos.left;
						}

						containerPosition.width = $element[0].offsetWidth;

						that.$bsContainer.css(containerPosition);
					};

			this.$button.on('click.bs.dropdown.data-api', function () {
				if (that.isDisabled()) {
					return;
				}

				getPlacement(that.$newElement);

				that.$bsContainer
					.appendTo(that.options.container)
					.toggleClass(classNames.SHOW, !that.$button.hasClass(classNames.SHOW))
					.append(that.$menu);
			});

			$(window)
				.off('resize' + EVENT_KEY + '.' + this.selectId + ' scroll' + EVENT_KEY + '.' + this.selectId)
				.on('resize' + EVENT_KEY + '.' + this.selectId + ' scroll' + EVENT_KEY + '.' + this.selectId, function () {
					var isActive = that.$newElement.hasClass(classNames.SHOW);

					if (isActive) getPlacement(that.$newElement);
				});

			this.$element.on('hide' + EVENT_KEY, function () {
				that.$menu.data('height', that.$menu.height());
				that.$bsContainer.detach();
			});
		},

		setOptionStatus: function () {
			var that = this;

			that.noScroll = false;

			if (that.selectpicker.view.visibleElements && that.selectpicker.view.visibleElements.length) {
				for (var i = 0; i < that.selectpicker.view.visibleElements.length; i++) {
					var liData = that.selectpicker.current.data[i + that.selectpicker.view.position0],
							option = liData.option;

					if (option) {
						that.setDisabled(
							liData.index,
							liData.disabled
						);

						that.setSelected(
							liData.index,
							option.selected
						);
					}
				}
			}
		},

		/**
		 * @param {number} index - the index of the option that is being changed
		 * @param {boolean} selected - true if the option is being selected, false if being deselected
		 */
		setSelected: function (index, selected) {
			var li = this.selectpicker.main.elements[index],
					liData = this.selectpicker.main.data[index],
					activeIndexIsSet = this.activeIndex !== undefined,
					thisIsActive = this.activeIndex === index,
					prevActive,
					a,
					// if current option is already active
					// OR
					// if the current option is being selected, it's NOT multiple, and
					// activeIndex is undefined:
					//  - when the menu is first being opened, OR
					//  - after a search has been performed, OR
					//  - when retainActive is false when selecting a new option (i.e. index of the newly selected option is not the same as the current activeIndex)
					keepActive = thisIsActive || (selected && !this.multiple && !activeIndexIsSet);

			liData.selected = selected;

			a = li.firstChild;

			if (selected) {
				this.selectedIndex = index;
			}

			li.classList.toggle('selected', selected);
			li.classList.toggle('active', keepActive);

			if (keepActive) {
				this.selectpicker.view.currentActive = li;
				this.activeIndex = index;
			}

			if (a) {
				a.classList.toggle('selected', selected);
				a.classList.toggle('active', keepActive);
				a.setAttribute('aria-selected', selected);
			}

			if (!keepActive) {
				if (!activeIndexIsSet && selected && this.prevActiveIndex !== undefined) {
					prevActive = this.selectpicker.main.elements[this.prevActiveIndex];

					prevActive.classList.remove('active');
					if (prevActive.firstChild) {
						prevActive.firstChild.classList.remove('active');
					}
				}
			}
		},

		/**
		 * @param {number} index - the index of the option that is being disabled
		 * @param {boolean} disabled - true if the option is being disabled, false if being enabled
		 */
		setDisabled: function (index, disabled) {
			var li = this.selectpicker.main.elements[index],
					a;

			this.selectpicker.main.data[index].disabled = disabled;

			a = li.firstChild;

			li.classList.toggle(classNames.DISABLED, disabled);

			if (a) {
				if (version.major === '4') a.classList.toggle(classNames.DISABLED, disabled);

				a.setAttribute('aria-disabled', disabled);

				if (disabled) {
					a.setAttribute('tabindex', -1);
				} else {
					a.setAttribute('tabindex', 0);
				}
			}
		},

		isDisabled: function () {
			return this.$element[0].disabled;
		},

		checkDisabled: function () {
			var that = this;

			if (this.isDisabled()) {
				this.$newElement[0].classList.add(classNames.DISABLED);
				this.$button.addClass(classNames.DISABLED).attr('tabindex', -1).attr('aria-disabled', true);
			} else {
				if (this.$button[0].classList.contains(classNames.DISABLED)) {
					this.$newElement[0].classList.remove(classNames.DISABLED);
					this.$button.removeClass(classNames.DISABLED).attr('aria-disabled', false);
				}

				if (this.$button.attr('tabindex') == -1 && !this.$element.data('tabindex')) {
					this.$button.removeAttr('tabindex');
				}
			}

			this.$button.on('click', function () {
				return !that.isDisabled();
			});
		},

		togglePlaceholder: function () {
			// much faster than calling $.val()
			var element = this.$element[0],
					selectedIndex = element.selectedIndex,
					nothingSelected = selectedIndex === -1;

			if (!nothingSelected && !element.options[selectedIndex].value) nothingSelected = true;

			this.$button.toggleClass('bs-placeholder', nothingSelected);
		},

		tabIndex: function () {
			if (this.$element.data('tabindex') !== this.$element.attr('tabindex') &&
				(this.$element.attr('tabindex') !== -98 && this.$element.attr('tabindex') !== '-98')) {
				this.$element.data('tabindex', this.$element.attr('tabindex'));
				this.$button.attr('tabindex', this.$element.data('tabindex'));
			}

			this.$element.attr('tabindex', -98);
		},

		clickListener: function () {
			var that = this,
					$document = $(document);

			$document.data('spaceSelect', false);

			this.$button.on('keyup', function (e) {
				if (/(32)/.test(e.keyCode.toString(10)) && $document.data('spaceSelect')) {
					e.preventDefault();
					$document.data('spaceSelect', false);
				}
			});

			this.$newElement.on('show.bs.dropdown', function () {
				if (version.major > 3 && !that.dropdown) {
					that.dropdown = that.$button.data('bs.dropdown');
					that.dropdown._menu = that.$menu[0];
				}
			});

			this.$button.on('click.bs.dropdown.data-api', function () {
				if (!that.$newElement.hasClass(classNames.SHOW)) {
					that.setSize();
				}
			});

			function setFocus () {
				if (that.options.liveSearch) {
					that.$searchbox.trigger('focus');
				} else {
					that.$menuInner.trigger('focus');
				}
			}

			function checkPopperExists () {
				if (that.dropdown && that.dropdown._popper && that.dropdown._popper.state.isCreated) {
					setFocus();
				} else {
					requestAnimationFrame(checkPopperExists);
				}
			}

			this.$element.on('shown' + EVENT_KEY, function () {
				if (that.$menuInner[0].scrollTop !== that.selectpicker.view.scrollTop) {
					that.$menuInner[0].scrollTop = that.selectpicker.view.scrollTop;
				}

				if (version.major > 3) {
					requestAnimationFrame(checkPopperExists);
				} else {
					setFocus();
				}
			});

			this.$menuInner.on('click', 'li a', function (e, retainActive) {
				var $this = $(this),
						position0 = that.isVirtual() ? that.selectpicker.view.position0 : 0,
						clickedData = that.selectpicker.current.data[$this.parent().index() + position0],
						clickedIndex = clickedData.index,
						prevValue = getSelectValues(that.$element[0]),
						prevIndex = that.$element.prop('selectedIndex'),
						triggerChange = true;

				// Don't close on multi choice menu
				if (that.multiple && that.options.maxOptions !== 1) {
					e.stopPropagation();
				}

				e.preventDefault();

				// Don't run if the select is disabled
				if (!that.isDisabled() && !$this.parent().hasClass(classNames.DISABLED)) {
					var $options = that.$element.find('option'),
							option = clickedData.option,
							$option = $(option),
							state = option.selected,
							$optgroup = $option.parent('optgroup'),
							$optgroupOptions = $optgroup.find('option'),
							maxOptions = that.options.maxOptions,
							maxOptionsGrp = $optgroup.data('maxOptions') || false;

					if (clickedIndex === that.activeIndex) retainActive = true;

					if (!retainActive) {
						that.prevActiveIndex = that.activeIndex;
						that.activeIndex = undefined;
					}

					if (!that.multiple) { // Deselect all others if not multi select box
						$options.prop('selected', false);
						option.selected = true;
						that.setSelected(clickedIndex, true);
					} else { // Toggle the one we have chosen if we are multi select.
						option.selected = !state;

						that.setSelected(clickedIndex, !state);
						$this.trigger('blur');

						if (maxOptions !== false || maxOptionsGrp !== false) {
							var maxReached = maxOptions < $options.filter(':selected').length,
									maxReachedGrp = maxOptionsGrp < $optgroup.find('option:selected').length;

							if ((maxOptions && maxReached) || (maxOptionsGrp && maxReachedGrp)) {
								if (maxOptions && maxOptions == 1) {
									$options.prop('selected', false);
									$option.prop('selected', true);

									for (var i = 0; i < $options.length; i++) {
										that.setSelected(i, false);
									}

									that.setSelected(clickedIndex, true);
								} else if (maxOptionsGrp && maxOptionsGrp == 1) {
									$optgroup.find('option:selected').prop('selected', false);
									$option.prop('selected', true);

									for (var i = 0; i < $optgroupOptions.length; i++) {
										var option = $optgroupOptions[i];
										that.setSelected($options.index(option), false);
									}

									that.setSelected(clickedIndex, true);
								} else {
									var maxOptionsText = typeof that.options.maxOptionsText === 'string' ? [that.options.maxOptionsText, that.options.maxOptionsText] : that.options.maxOptionsText,
											maxOptionsArr = typeof maxOptionsText === 'function' ? maxOptionsText(maxOptions, maxOptionsGrp) : maxOptionsText,
											maxTxt = maxOptionsArr[0].replace('{n}', maxOptions),
											maxTxtGrp = maxOptionsArr[1].replace('{n}', maxOptionsGrp),
											$notify = $('<div class="notify"></div>');
									// If {var} is set in array, replace it
									/** @deprecated */
									if (maxOptionsArr[2]) {
										maxTxt = maxTxt.replace('{var}', maxOptionsArr[2][maxOptions > 1 ? 0 : 1]);
										maxTxtGrp = maxTxtGrp.replace('{var}', maxOptionsArr[2][maxOptionsGrp > 1 ? 0 : 1]);
									}

									$option.prop('selected', false);

									that.$menu.append($notify);

									if (maxOptions && maxReached) {
										$notify.append($('<div>' + maxTxt + '</div>'));
										triggerChange = false;
										that.$element.trigger('maxReached' + EVENT_KEY);
									}

									if (maxOptionsGrp && maxReachedGrp) {
										$notify.append($('<div>' + maxTxtGrp + '</div>'));
										triggerChange = false;
										that.$element.trigger('maxReachedGrp' + EVENT_KEY);
									}

									setTimeout(function () {
										that.setSelected(clickedIndex, false);
									}, 10);

									$notify.delay(750).fadeOut(300, function () {
										$(this).remove();
									});
								}
							}
						}
					}

					if (!that.multiple || (that.multiple && that.options.maxOptions === 1)) {
						that.$button.trigger('focus');
					} else if (that.options.liveSearch) {
						that.$searchbox.trigger('focus');
					}

					// Trigger select 'change'
					if (triggerChange) {
						if ((prevValue != getSelectValues(that.$element[0]) && that.multiple) || (prevIndex != that.$element.prop('selectedIndex') && !that.multiple)) {
							// $option.prop('selected') is current option state (selected/unselected). prevValue is the value of the select prior to being changed.
							changedArguments = [option.index, $option.prop('selected'), prevValue];
							that.$element
								.triggerNative('change');
						}
					}
				}
			});

			this.$menu.on('click', 'li.' + classNames.DISABLED + ' a, .' + classNames.POPOVERHEADER + ', .' + classNames.POPOVERHEADER + ' :not(.close)', function (e) {
				if (e.currentTarget == this) {
					e.preventDefault();
					e.stopPropagation();
					if (that.options.liveSearch && !$(e.target).hasClass('close')) {
						that.$searchbox.trigger('focus');
					} else {
						that.$button.trigger('focus');
					}
				}
			});

			this.$menuInner.on('click', '.divider, .dropdown-header', function (e) {
				e.preventDefault();
				e.stopPropagation();
				if (that.options.liveSearch) {
					that.$searchbox.trigger('focus');
				} else {
					that.$button.trigger('focus');
				}
			});

			this.$menu.on('click', '.' + classNames.POPOVERHEADER + ' .close', function () {
				that.$button.trigger('click');
			});

			this.$searchbox.on('click', function (e) {
				e.stopPropagation();
			});

			this.$menu.on('click', '.actions-btn', function (e) {
				if (that.options.liveSearch) {
					that.$searchbox.trigger('focus');
				} else {
					that.$button.trigger('focus');
				}

				e.preventDefault();
				e.stopPropagation();

				if ($(this).hasClass('bs-select-all')) {
					that.selectAll();
				} else {
					that.deselectAll();
				}
			});

			this.$element
				.on('change' + EVENT_KEY, function () {
					that.render();
					that.$element.trigger('changed' + EVENT_KEY, changedArguments);
					changedArguments = null;
				})
				.on('focus' + EVENT_KEY, function () {
					if (!that.options.mobile) that.$button.trigger('focus');
				});
		},

		liveSearchListener: function () {
			var that = this,
					noResults = document.createElement('li');

			this.$button.on('click.bs.dropdown.data-api', function () {
				if (!!that.$searchbox.val()) {
					that.$searchbox.val('');
				}
			});

			this.$searchbox.on('click.bs.dropdown.data-api focus.bs.dropdown.data-api touchend.bs.dropdown.data-api', function (e) {
				e.stopPropagation();
			});

			this.$searchbox.on('input propertychange', function () {
				var searchValue = that.$searchbox.val();

				that.selectpicker.search.elements = [];
				that.selectpicker.search.data = [];

				if (searchValue) {
					var i,
							searchMatch = [],
							q = searchValue.toUpperCase(),
							cache = {},
							cacheArr = [],
							searchStyle = that._searchStyle(),
							normalizeSearch = that.options.liveSearchNormalize;

					if (normalizeSearch) q = normalizeToBase(q);

					that._$lisSelected = that.$menuInner.find('.selected');

					for (var i = 0; i < that.selectpicker.main.data.length; i++) {
						var li = that.selectpicker.main.data[i];

						if (!cache[i]) {
							cache[i] = stringSearch(li, q, searchStyle, normalizeSearch);
						}

						if (cache[i] && li.headerIndex !== undefined && cacheArr.indexOf(li.headerIndex) === -1) {
							if (li.headerIndex > 0) {
								cache[li.headerIndex - 1] = true;
								cacheArr.push(li.headerIndex - 1);
							}

							cache[li.headerIndex] = true;
							cacheArr.push(li.headerIndex);

							cache[li.lastIndex + 1] = true;
						}

						if (cache[i] && li.type !== 'optgroup-label') cacheArr.push(i);
					}

					for (var i = 0, cacheLen = cacheArr.length; i < cacheLen; i++) {
						var index = cacheArr[i],
								prevIndex = cacheArr[i - 1],
								li = that.selectpicker.main.data[index],
								liPrev = that.selectpicker.main.data[prevIndex];

						if (li.type !== 'divider' || (li.type === 'divider' && liPrev && liPrev.type !== 'divider' && cacheLen - 1 !== i)) {
							that.selectpicker.search.data.push(li);
							searchMatch.push(that.selectpicker.main.elements[index]);
						}
					}

					that.activeIndex = undefined;
					that.noScroll = true;
					that.$menuInner.scrollTop(0);
					that.selectpicker.search.elements = searchMatch;
					that.createView(true);

					if (!searchMatch.length) {
						noResults.className = 'no-results';
						noResults.innerHTML = that.options.noneResultsText.replace('{0}', '"' + htmlEscape(searchValue) + '"');
						that.$menuInner[0].firstChild.appendChild(noResults);
					}
				} else {
					that.$menuInner.scrollTop(0);
					that.createView(false);
				}
			});
		},

		_searchStyle: function () {
			return this.options.liveSearchStyle || 'contains';
		},

		val: function (value) {
			if (typeof value !== 'undefined') {
				var prevValue = getSelectValues(this.$element[0]);

				changedArguments = [null, null, prevValue];

				this.$element
					.val(value)
					.trigger('changed' + EVENT_KEY, changedArguments);

				this.render();

				changedArguments = null;

				return this.$element;
			} else {
				return this.$element.val();
			}
		},

		changeAll: function (status) {
			if (!this.multiple) return;
			if (typeof status === 'undefined') status = true;

			var element = this.$element[0],
					previousSelected = 0,
					currentSelected = 0,
					prevValue = getSelectValues(element);

			element.classList.add('bs-select-hidden');

			for (var i = 0, len = this.selectpicker.current.elements.length; i < len; i++) {
				var liData = this.selectpicker.current.data[i],
						option = liData.option;

				if (option && !liData.disabled && liData.type !== 'divider') {
					if (liData.selected) previousSelected++;
					option.selected = status;
					if (status) currentSelected++;
				}
			}

			element.classList.remove('bs-select-hidden');

			if (previousSelected === currentSelected) return;

			this.setOptionStatus();

			this.togglePlaceholder();

			changedArguments = [null, null, prevValue];

			this.$element
				.triggerNative('change');
		},

		selectAll: function () {
			return this.changeAll(true);
		},

		deselectAll: function () {
			return this.changeAll(false);
		},

		toggle: function (e) {
			e = e || window.event;

			if (e) e.stopPropagation();

			this.$button.trigger('click.bs.dropdown.data-api');
		},

		keydown: function (e) {
			var $this = $(this),
					isToggle = $this.hasClass('dropdown-toggle'),
					$parent = isToggle ? $this.closest('.dropdown') : $this.closest(Selector.MENU),
					that = $parent.data('this'),
					$items = that.findLis(),
					index,
					isActive,
					liActive,
					activeLi,
					offset,
					updateScroll = false,
					downOnTab = e.which === keyCodes.TAB && !isToggle && !that.options.selectOnTab,
					isArrowKey = REGEXP_ARROW.test(e.which) || downOnTab,
					scrollTop = that.$menuInner[0].scrollTop,
					isVirtual = that.isVirtual(),
					position0 = isVirtual === true ? that.selectpicker.view.position0 : 0;

			isActive = that.$newElement.hasClass(classNames.SHOW);

			if (
				!isActive &&
				(
					isArrowKey ||
					(e.which >= 48 && e.which <= 57) ||
					(e.which >= 96 && e.which <= 105) ||
					(e.which >= 65 && e.which <= 90)
				)
			) {
				that.$button.trigger('click.bs.dropdown.data-api');

				if (that.options.liveSearch) {
					that.$searchbox.trigger('focus');
					return;
				}
			}

			if (e.which === keyCodes.ESCAPE && isActive) {
				e.preventDefault();
				that.$button.trigger('click.bs.dropdown.data-api').trigger('focus');
			}

			if (isArrowKey) { // if up or down
				if (!$items.length) return;

				// $items.index/.filter is too slow with a large list and no virtual scroll
				index = isVirtual === true ? $items.index($items.filter('.active')) : that.activeIndex;

				if (index === undefined) index = -1;

				if (index !== -1) {
					liActive = that.selectpicker.current.elements[index + position0];
					liActive.classList.remove('active');
					if (liActive.firstChild) liActive.firstChild.classList.remove('active');
				}

				if (e.which === keyCodes.ARROW_UP) { // up
					if (index !== -1) index--;
					if (index + position0 < 0) index += $items.length;

					if (!that.selectpicker.view.canHighlight[index + position0]) {
						index = that.selectpicker.view.canHighlight.slice(0, index + position0).lastIndexOf(true) - position0;
						if (index === -1) index = $items.length - 1;
					}
				} else if (e.which === keyCodes.ARROW_DOWN || downOnTab) { // down
					index++;
					if (index + position0 >= that.selectpicker.view.canHighlight.length) index = 0;

					if (!that.selectpicker.view.canHighlight[index + position0]) {
						index = index + 1 + that.selectpicker.view.canHighlight.slice(index + position0 + 1).indexOf(true);
					}
				}

				e.preventDefault();

				var liActiveIndex = position0 + index;

				if (e.which === keyCodes.ARROW_UP) { // up
					// scroll to bottom and highlight last option
					if (position0 === 0 && index === $items.length - 1) {
						that.$menuInner[0].scrollTop = that.$menuInner[0].scrollHeight;

						liActiveIndex = that.selectpicker.current.elements.length - 1;
					} else {
						activeLi = that.selectpicker.current.data[liActiveIndex];
						offset = activeLi.position - activeLi.height;

						updateScroll = offset < scrollTop;
					}
				} else if (e.which === keyCodes.ARROW_DOWN || downOnTab) { // down
					// scroll to top and highlight first option
					if (index === 0) {
						that.$menuInner[0].scrollTop = 0;

						liActiveIndex = 0;
					} else {
						activeLi = that.selectpicker.current.data[liActiveIndex];
						offset = activeLi.position - that.sizeInfo.menuInnerHeight;

						updateScroll = offset > scrollTop;
					}
				}

				liActive = that.selectpicker.current.elements[liActiveIndex];

				if (liActive) {
					liActive.classList.add('active');
					if (liActive.firstChild) liActive.firstChild.classList.add('active');
				}

				that.activeIndex = that.selectpicker.current.data[liActiveIndex].index;

				that.selectpicker.view.currentActive = liActive;

				if (updateScroll) that.$menuInner[0].scrollTop = offset;

				if (that.options.liveSearch) {
					that.$searchbox.trigger('focus');
				} else {
					$this.trigger('focus');
				}
			} else if (
				(!$this.is('input') && !REGEXP_TAB_OR_ESCAPE.test(e.which)) ||
				(e.which === keyCodes.SPACE && that.selectpicker.keydown.keyHistory)
			) {
				var searchMatch,
						matches = [],
						keyHistory;

				e.preventDefault();

				that.selectpicker.keydown.keyHistory += keyCodeMap[e.which];

				if (that.selectpicker.keydown.resetKeyHistory.cancel) clearTimeout(that.selectpicker.keydown.resetKeyHistory.cancel);
				that.selectpicker.keydown.resetKeyHistory.cancel = that.selectpicker.keydown.resetKeyHistory.start();

				keyHistory = that.selectpicker.keydown.keyHistory;

				// if all letters are the same, set keyHistory to just the first character when searching
				if (/^(.)\1+$/.test(keyHistory)) {
					keyHistory = keyHistory.charAt(0);
				}

				// find matches
				for (var i = 0; i < that.selectpicker.current.data.length; i++) {
					var li = that.selectpicker.current.data[i],
							hasMatch;

					hasMatch = stringSearch(li, keyHistory, 'startsWith', true);

					if (hasMatch && that.selectpicker.view.canHighlight[i]) {
						matches.push(li.index);
					}
				}

				if (matches.length) {
					var matchIndex = 0;

					$items.removeClass('active').find('a').removeClass('active');

					// either only one key has been pressed or they are all the same key
					if (keyHistory.length === 1) {
						matchIndex = matches.indexOf(that.activeIndex);

						if (matchIndex === -1 || matchIndex === matches.length - 1) {
							matchIndex = 0;
						} else {
							matchIndex++;
						}
					}

					searchMatch = matches[matchIndex];

					activeLi = that.selectpicker.main.data[searchMatch];

					if (scrollTop - activeLi.position > 0) {
						offset = activeLi.position - activeLi.height;
						updateScroll = true;
					} else {
						offset = activeLi.position - that.sizeInfo.menuInnerHeight;
						// if the option is already visible at the current scroll position, just keep it the same
						updateScroll = activeLi.position > scrollTop + that.sizeInfo.menuInnerHeight;
					}

					liActive = that.selectpicker.main.elements[searchMatch];
					liActive.classList.add('active');
					if (liActive.firstChild) liActive.firstChild.classList.add('active');
					that.activeIndex = matches[matchIndex];

					liActive.firstChild.focus();

					if (updateScroll) that.$menuInner[0].scrollTop = offset;

					$this.trigger('focus');
				}
			}

			// Select focused option if "Enter", "Spacebar" or "Tab" (when selectOnTab is true) are pressed inside the menu.
			if (
				isActive &&
				(
					(e.which === keyCodes.SPACE && !that.selectpicker.keydown.keyHistory) ||
					e.which === keyCodes.ENTER ||
					(e.which === keyCodes.TAB && that.options.selectOnTab)
				)
			) {
				if (e.which !== keyCodes.SPACE) e.preventDefault();

				if (!that.options.liveSearch || e.which !== keyCodes.SPACE) {
					that.$menuInner.find('.active a').trigger('click', true); // retain active class
					$this.trigger('focus');

					if (!that.options.liveSearch) {
						// Prevent screen from scrolling if the user hits the spacebar
						e.preventDefault();
						// Fixes spacebar selection of dropdown items in FF & IE
						$(document).data('spaceSelect', true);
					}
				}
			}
		},

		mobile: function () {
			this.$element[0].classList.add('mobile-device');
		},

		refresh: function () {
			// update options if data attributes have been changed
			var config = $.extend({}, this.options, this.$element.data());
			this.options = config;

			this.checkDisabled();
			this.setStyle();
			this.render();
			this.createLi();
			this.setWidth();

			this.setSize(true);

			this.$element.trigger('refreshed' + EVENT_KEY);
		},

		hide: function () {
			this.$newElement.hide();
		},

		show: function () {
			this.$newElement.show();
		},

		remove: function () {
			this.$newElement.remove();
			this.$element.remove();
		},

		destroy: function () {
			this.$newElement.before(this.$element).remove();

			if (this.$bsContainer) {
				this.$bsContainer.remove();
			} else {
				this.$menu.remove();
			}

			this.$element
				.off(EVENT_KEY)
				.removeData('selectpicker')
				.removeClass('bs-select-hidden selectpicker');

			$(window).off(EVENT_KEY + '.' + this.selectId);
		}
	};

	// SELECTPICKER PLUGIN DEFINITION
	// ==============================
	function Plugin (option) {
		// get the args of the outer function..
		var args = arguments;
		// The arguments of the function are explicitly re-defined from the argument list, because the shift causes them
		// to get lost/corrupted in android 2.3 and IE9 #715 #775
		var _option = option;

		[].shift.apply(args);

		// if the version was not set successfully
		if (!version.success) {
			// try to retreive it again
			try {
				version.full = ($.fn.dropdown.Constructor.VERSION || '').split(' ')[0].split('.');
			} catch (err) {
				// fall back to use BootstrapVersion if set
				if (Selectpicker.BootstrapVersion) {
					version.full = Selectpicker.BootstrapVersion.split(' ')[0].split('.');
				} else {
					version.full = [version.major, '0', '0'];

					console.warn(
						'There was an issue retrieving Bootstrap\'s version. ' +
						'Ensure Bootstrap is being loaded before bootstrap-select and there is no namespace collision. ' +
						'If loading Bootstrap asynchronously, the version may need to be manually specified via $.fn.selectpicker.Constructor.BootstrapVersion.',
						err
					);
				}
			}

			version.major = version.full[0];
			version.success = true;
		}

		if (version.major === '4') {
			// some defaults need to be changed if using Bootstrap 4
			// check to see if they have already been manually changed before forcing them to update
			var toUpdate = [];

			if (Selectpicker.DEFAULTS.style === classNames.BUTTONCLASS) toUpdate.push({ name: 'style', className: 'BUTTONCLASS' });
			if (Selectpicker.DEFAULTS.iconBase === classNames.ICONBASE) toUpdate.push({ name: 'iconBase', className: 'ICONBASE' });
			if (Selectpicker.DEFAULTS.tickIcon === classNames.TICKICON) toUpdate.push({ name: 'tickIcon', className: 'TICKICON' });

			classNames.DIVIDER = 'dropdown-divider';
			classNames.SHOW = 'show';
			classNames.BUTTONCLASS = 'btn-light';
			classNames.POPOVERHEADER = 'popover-header';
			classNames.ICONBASE = '';
			classNames.TICKICON = 'bs-ok-default';

			for (var i = 0; i < toUpdate.length; i++) {
				var option = toUpdate[i];
				Selectpicker.DEFAULTS[option.name] = classNames[option.className];
			}
		}

		var value;
		var chain = this.each(function () {
			var $this = $(this);
			if ($this.is('select')) {
				var data = $this.data('selectpicker'),
						options = typeof _option == 'object' && _option;

				if (!data) {
					var dataAttributes = $this.data();

					for (var dataAttr in dataAttributes) {
						if (dataAttributes.hasOwnProperty(dataAttr) && $.inArray(dataAttr, DISALLOWED_ATTRIBUTES) !== -1) {
							delete dataAttributes[dataAttr];
						}
					}

					var config = $.extend({}, Selectpicker.DEFAULTS, $.fn.selectpicker.defaults || {}, dataAttributes, options);
					config.template = $.extend({}, Selectpicker.DEFAULTS.template, ($.fn.selectpicker.defaults ? $.fn.selectpicker.defaults.template : {}), dataAttributes.template, options.template);
					$this.data('selectpicker', (data = new Selectpicker(this, config)));
				} else if (options) {
					for (var i in options) {
						if (options.hasOwnProperty(i)) {
							data.options[i] = options[i];
						}
					}
				}

				if (typeof _option == 'string') {
					if (data[_option] instanceof Function) {
						value = data[_option].apply(data, args);
					} else {
						value = data.options[_option];
					}
				}
			}
		});

		if (typeof value !== 'undefined') {
			// noinspection JSUnusedAssignment
			return value;
		} else {
			return chain;
		}
	}

	var old = $.fn.selectpicker;
	$.fn.selectpicker = Plugin;
	$.fn.selectpicker.Constructor = Selectpicker;

	// SELECTPICKER NO CONFLICT
	// ========================
	$.fn.selectpicker.noConflict = function () {
		$.fn.selectpicker = old;
		return this;
	};

	$(document)
		.off('keydown.bs.dropdown.data-api')
		.on('keydown' + EVENT_KEY, '.bootstrap-select [data-toggle="dropdown"], .bootstrap-select [role="listbox"], .bootstrap-select .bs-searchbox input', Selectpicker.prototype.keydown)
		.on('focusin.modal', '.bootstrap-select [data-toggle="dropdown"], .bootstrap-select [role="listbox"], .bootstrap-select .bs-searchbox input', function (e) {
			e.stopPropagation();
		});

	// SELECTPICKER DATA-API
	// =====================
	$(window).on('load' + EVENT_KEY + '.data-api', function () {
		$('.selectpicker').each(function () {
			var $selectpicker = $(this);
			Plugin.call($selectpicker, $selectpicker.data());
		})
	});
})(jQuery);
}));

/******************************************************************************/

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
		// IE does not support 'new Event()'
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

/******************************************************************************/

/**!
 * Sortable 1.10.0-rc3
 * @author  RubaXa   <trash@rubaxa.org>
 * @author  owenm    <owen23355@gmail.com>
 * @license MIT
 */
(function (global, factory) {
	typeof exports === 'object' && typeof module !== 'undefined' ? module.exports = factory() :
	typeof define === 'function' && define.amd ? define(factory) :
	(global = global || self, global.Sortable = factory());
}(this, function () { 'use strict';

	function _typeof(obj) {
		if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") {
			_typeof = function (obj) {
				return typeof obj;
			};
		} else {
			_typeof = function (obj) {
				return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj;
			};
		}

		return _typeof(obj);
	}

	function _defineProperty(obj, key, value) {
		if (key in obj) {
			Object.defineProperty(obj, key, {
				value: value,
				enumerable: true,
				configurable: true,
				writable: true
			});
		} else {
			obj[key] = value;
		}

		return obj;
	}

	function _extends() {
		_extends = Object.assign || function (target) {
			for (var i = 1; i < arguments.length; i++) {
				var source = arguments[i];

				for (var key in source) {
					if (Object.prototype.hasOwnProperty.call(source, key)) {
						target[key] = source[key];
					}
				}
			}

			return target;
		};

		return _extends.apply(this, arguments);
	}

	function _objectSpread(target) {
		for (var i = 1; i < arguments.length; i++) {
			var source = arguments[i] != null ? arguments[i] : {};
			var ownKeys = Object.keys(source);

			if (typeof Object.getOwnPropertySymbols === 'function') {
				ownKeys = ownKeys.concat(Object.getOwnPropertySymbols(source).filter(function (sym) {
					return Object.getOwnPropertyDescriptor(source, sym).enumerable;
				}));
			}

			ownKeys.forEach(function (key) {
				_defineProperty(target, key, source[key]);
			});
		}

		return target;
	}

	function _objectWithoutPropertiesLoose(source, excluded) {
		if (source == null) return {};
		var target = {};
		var sourceKeys = Object.keys(source);
		var key, i;

		for (i = 0; i < sourceKeys.length; i++) {
			key = sourceKeys[i];
			if (excluded.indexOf(key) >= 0) continue;
			target[key] = source[key];
		}

		return target;
	}

	function _objectWithoutProperties(source, excluded) {
		if (source == null) return {};

		var target = _objectWithoutPropertiesLoose(source, excluded);

		var key, i;

		if (Object.getOwnPropertySymbols) {
			var sourceSymbolKeys = Object.getOwnPropertySymbols(source);

			for (i = 0; i < sourceSymbolKeys.length; i++) {
				key = sourceSymbolKeys[i];
				if (excluded.indexOf(key) >= 0) continue;
				if (!Object.prototype.propertyIsEnumerable.call(source, key)) continue;
				target[key] = source[key];
			}
		}

		return target;
	}

	function _toConsumableArray(arr) {
		return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _nonIterableSpread();
	}

	function _arrayWithoutHoles(arr) {
		if (Array.isArray(arr)) {
			for (var i = 0, arr2 = new Array(arr.length); i < arr.length; i++) arr2[i] = arr[i];

			return arr2;
		}
	}

	function _iterableToArray(iter) {
		if (Symbol.iterator in Object(iter) || Object.prototype.toString.call(iter) === "[object Arguments]") return Array.from(iter);
	}

	function _nonIterableSpread() {
		throw new TypeError("Invalid attempt to spread non-iterable instance");
	}

	var version = "1.10.0-rc3";

	function userAgent(pattern) {
		return !!
		/*@__PURE__*/
		navigator.userAgent.match(pattern);
	}

	var IE11OrLess = userAgent(/(?:Trident.*rv[ :]?11\.|msie|iemobile|Windows Phone)/i);
	var Edge = userAgent(/Edge/i);
	var FireFox = userAgent(/firefox/i);
	var Safari = userAgent(/safari/i) && !userAgent(/chrome/i) && !userAgent(/android/i);
	var IOS = userAgent(/iP(ad|od|hone)/i);
	var ChromeForAndroid = userAgent(/chrome/i) && userAgent(/android/i);

	var captureMode = {
		capture: false,
		passive: false
	};

	function on(el, event, fn) {
		el.addEventListener(event, fn, !IE11OrLess && captureMode);
	}

	function off(el, event, fn) {
		el.removeEventListener(event, fn, !IE11OrLess && captureMode);
	}

	function matches(
	/**HTMLElement*/
	el,
	/**String*/
	selector) {
		if (!selector) return;
		selector[0] === '>' && (selector = selector.substring(1));

		if (el) {
			try {
				if (el.matches) {
					return el.matches(selector);
				} else if (el.msMatchesSelector) {
					return el.msMatchesSelector(selector);
				} else if (el.webkitMatchesSelector) {
					return el.webkitMatchesSelector(selector);
				}
			} catch (_) {
				return false;
			}
		}

		return false;
	}

	function getParentOrHost(el) {
		return el.host && el !== document && el.host.nodeType ? el.host : el.parentNode;
	}

	function closest(
	/**HTMLElement*/
	el,
	/**String*/
	selector,
	/**HTMLElement*/
	ctx, includeCTX) {
		if (el) {
			ctx = ctx || document;

			do {
				if (selector != null && (selector[0] === '>' ? el.parentNode === ctx && matches(el, selector) : matches(el, selector)) || includeCTX && el === ctx) {
					return el;
				}

				if (el === ctx) break;
				/* jshint boss:true */
			} while (el = getParentOrHost(el));
		}

		return null;
	}

	var R_SPACE = /\s+/g;

	function toggleClass(el, name, state) {
		if (el && name) {
			if (el.classList) {
				el.classList[state ? 'add' : 'remove'](name);
			} else {
				var className = (' ' + el.className + ' ').replace(R_SPACE, ' ').replace(' ' + name + ' ', ' ');
				el.className = (className + (state ? ' ' + name : '')).replace(R_SPACE, ' ');
			}
		}
	}

	function css(el, prop, val) {
		var style = el && el.style;

		if (style) {
			if (val === void 0) {
				if (document.defaultView && document.defaultView.getComputedStyle) {
					val = document.defaultView.getComputedStyle(el, '');
				} else if (el.currentStyle) {
					val = el.currentStyle;
				}

				return prop === void 0 ? val : val[prop];
			} else {
				if (!(prop in style) && prop.indexOf('webkit') === -1) {
					prop = '-webkit-' + prop;
				}

				style[prop] = val + (typeof val === 'string' ? '' : 'px');
			}
		}
	}

	function matrix(el, selfOnly) {
		var appliedTransforms = '';

		do {
			var transform = css(el, 'transform');

			if (transform && transform !== 'none') {
				appliedTransforms = transform + ' ' + appliedTransforms;
			}
			/* jshint boss:true */

		} while (!selfOnly && (el = el.parentNode));

		var matrixFn = window.DOMMatrix || window.WebKitCSSMatrix || window.CSSMatrix;
		/*jshint -W056 */

		return matrixFn && new matrixFn(appliedTransforms);
	}

	function find(ctx, tagName, iterator) {
		if (ctx) {
			var list = ctx.getElementsByTagName(tagName),
					i = 0,
					n = list.length;

			if (iterator) {
				for (; i < n; i++) {
					iterator(list[i], i);
				}
			}

			return list;
		}

		return [];
	}

	function getWindowScrollingElement() {
		if (IE11OrLess) {
			return document.documentElement;
		} else {
			return document.scrollingElement;
		}
	}
	/**
	 * Returns the "bounding client rect" of given element
	 * @param  {HTMLElement} el                       The element whose boundingClientRect is wanted
	 * @param  {[Boolean]} relativeToContainingBlock  Whether the rect should be relative to the containing block of (including) the container
	 * @param  {[Boolean]} relativeToNonStaticParent  Whether the rect should be relative to the relative parent of (including) the contaienr
	 * @param  {[Boolean]} undoScale                  Whether the container's scale() should be undone
	 * @param  {[HTMLElement]} container              The parent the element will be placed in
	 * @return {Object}                               The boundingClientRect of el, with specified adjustments
	 */


	function getRect(el, relativeToContainingBlock, relativeToNonStaticParent, undoScale, container) {
		if (!el.getBoundingClientRect && el !== window) return;
		var elRect, top, left, bottom, right, height, width;

		if (el !== window && el !== getWindowScrollingElement()) {
			elRect = el.getBoundingClientRect();
			top = elRect.top;
			left = elRect.left;
			bottom = elRect.bottom;
			right = elRect.right;
			height = elRect.height;
			width = elRect.width;
		} else {
			top = 0;
			left = 0;
			bottom = window.innerHeight;
			right = window.innerWidth;
			height = window.innerHeight;
			width = window.innerWidth;
		}

		if ((relativeToContainingBlock || relativeToNonStaticParent) && el !== window) {
			// Adjust for translate()
			container = container || el.parentNode; // solves #1123 (see: https://stackoverflow.com/a/37953806/6088312)
			// Not needed on <= IE11

			if (!IE11OrLess) {
				do {
					if (container && container.getBoundingClientRect && (css(container, 'transform') !== 'none' || relativeToNonStaticParent && css(container, 'position') !== 'static')) {
						var containerRect = container.getBoundingClientRect(); // Set relative to edges of padding box of container

						top -= containerRect.top + parseInt(css(container, 'border-top-width'));
						left -= containerRect.left + parseInt(css(container, 'border-left-width'));
						bottom = top + elRect.height;
						right = left + elRect.width;
						break;
					}
					/* jshint boss:true */

				} while (container = container.parentNode);
			}
		}

		if (undoScale && el !== window) {
			// Adjust for scale()
			var elMatrix = matrix(container || el),
					scaleX = elMatrix && elMatrix.a,
					scaleY = elMatrix && elMatrix.d;

			if (elMatrix) {
				top /= scaleY;
				left /= scaleX;
				width /= scaleX;
				height /= scaleY;
				bottom = top + height;
				right = left + width;
			}
		}

		return {
			top: top,
			left: left,
			bottom: bottom,
			right: right,
			width: width,
			height: height
		};
	}
	/**
	 * Checks if a side of an element is scrolled past a side of its parents
	 * @param  {HTMLElement}  el           The element who's side being scrolled out of view is in question
	 * @param  {[DOMRect]}    rect         Optional rect of el to use
	 * @param  {String}       elSide       Side of the element in question ('top', 'left', 'right', 'bottom')
	 * @param  {String}       parentSide   Side of the parent in question ('top', 'left', 'right', 'bottom')
	 * @return {HTMLElement}               The parent scroll element that the el's side is scrolled past, or null if there is no such element
	 */


	function isScrolledPast(el, rect, elSide, parentSide) {
		var parent = getParentAutoScrollElement(el, true),
				elSideVal = (rect ? rect : getRect(el))[elSide];
		/* jshint boss:true */

		while (parent) {
			var parentSideVal = getRect(parent)[parentSide],
					visible = void 0;

			if (parentSide === 'top' || parentSide === 'left') {
				visible = elSideVal >= parentSideVal;
			} else {
				visible = elSideVal <= parentSideVal;
			}

			if (!visible) return parent;
			if (parent === getWindowScrollingElement()) break;
			parent = getParentAutoScrollElement(parent, false);
		}

		return false;
	}
	/**
	 * Gets nth child of el, ignoring hidden children, sortable's elements (does not ignore clone if it's visible)
	 * and non-draggable elements
	 * @param  {HTMLElement} el       The parent element
	 * @param  {Number} childNum      The index of the child
	 * @param  {Object} options       Parent Sortable's options
	 * @return {HTMLElement}          The child at index childNum, or null if not found
	 */


	function getChild(el, childNum, options) {
		var currentChild = 0,
				i = 0,
				children = el.children;

		while (i < children.length) {
			if (children[i].style.display !== 'none' && children[i] !== Sortable.ghost && children[i] !== Sortable.dragged && closest(children[i], options.draggable, el, false)) {
				if (currentChild === childNum) {
					return children[i];
				}

				currentChild++;
			}

			i++;
		}

		return null;
	}
	/**
	 * Gets the last child in the el, ignoring ghostEl or invisible elements (clones)
	 * @param  {HTMLElement} el       Parent element
	 * @param  {selector} selector    Any other elements that should be ignored
	 * @return {HTMLElement}          The last child, ignoring ghostEl
	 */


	function lastChild(el, selector) {
		var last = el.lastElementChild;

		while (last && (last === Sortable.ghost || css(last, 'display') === 'none' || selector && !matches(last, selector))) {
			last = last.previousElementSibling;
		}

		return last || null;
	}
	/**
	 * Returns the index of an element within its parent for a selected set of
	 * elements
	 * @param  {HTMLElement} el
	 * @param  {selector} selector
	 * @return {number}
	 */


	function index(el, selector) {
		var index = 0;

		if (!el || !el.parentNode) {
			return -1;
		}
		/* jshint boss:true */


		while (el = el.previousElementSibling) {
			if (el.nodeName.toUpperCase() !== 'TEMPLATE' && el !== Sortable.clone && (!selector || matches(el, selector))) {
				index++;
			}
		}

		return index;
	}
	/**
	 * Returns the scroll offset of the given element, added with all the scroll offsets of parent elements.
	 * The value is returned in real pixels.
	 * @param  {HTMLElement} el
	 * @return {Array}             Offsets in the format of [left, top]
	 */


	function getRelativeScrollOffset(el) {
		var offsetLeft = 0,
				offsetTop = 0,
				winScroller = getWindowScrollingElement();

		if (el) {
			do {
				var elMatrix = matrix(el),
						scaleX = elMatrix.a,
						scaleY = elMatrix.d;
				offsetLeft += el.scrollLeft * scaleX;
				offsetTop += el.scrollTop * scaleY;
			} while (el !== winScroller && (el = el.parentNode));
		}

		return [offsetLeft, offsetTop];
	}
	/**
	 * Returns the index of the object within the given array
	 * @param  {Array} arr   Array that may or may not hold the object
	 * @param  {Object} obj  An object that has a key-value pair unique to and identical to a key-value pair in the object you want to find
	 * @return {Number}      The index of the object in the array, or -1
	 */


	function indexOfObject(arr, obj) {
		for (var i in arr) {
			if (!arr.hasOwnProperty(i)) continue;

			for (var key in obj) {
				if (obj.hasOwnProperty(key) && obj[key] === arr[i][key]) return Number(i);
			}
		}

		return -1;
	}

	function getParentAutoScrollElement(el, includeSelf) {
		// skip to window
		if (!el || !el.getBoundingClientRect) return getWindowScrollingElement();
		var elem = el;
		var gotSelf = false;

		do {
			// we don't need to get elem css if it isn't even overflowing in the first place (performance)
			if (elem.clientWidth < elem.scrollWidth || elem.clientHeight < elem.scrollHeight) {
				var elemCSS = css(elem);

				if (elem.clientWidth < elem.scrollWidth && (elemCSS.overflowX == 'auto' || elemCSS.overflowX == 'scroll') || elem.clientHeight < elem.scrollHeight && (elemCSS.overflowY == 'auto' || elemCSS.overflowY == 'scroll')) {
					if (!elem.getBoundingClientRect || elem === document.body) return getWindowScrollingElement();
					if (gotSelf || includeSelf) return elem;
					gotSelf = true;
				}
			}
			/* jshint boss:true */

		} while (elem = elem.parentNode);

		return getWindowScrollingElement();
	}

	function extend(dst, src) {
		if (dst && src) {
			for (var key in src) {
				if (src.hasOwnProperty(key)) {
					dst[key] = src[key];
				}
			}
		}

		return dst;
	}

	function isRectEqual(rect1, rect2) {
		return Math.round(rect1.top) === Math.round(rect2.top) && Math.round(rect1.left) === Math.round(rect2.left) && Math.round(rect1.height) === Math.round(rect2.height) && Math.round(rect1.width) === Math.round(rect2.width);
	}

	var _throttleTimeout;

	function throttle(callback, ms) {
		return function () {
			if (!_throttleTimeout) {
				var args = arguments,
						_this = this;

				if (args.length === 1) {
					callback.call(_this, args[0]);
				} else {
					callback.apply(_this, args);
				}

				_throttleTimeout = setTimeout(function () {
					_throttleTimeout = void 0;
				}, ms);
			}
		};
	}

	function cancelThrottle() {
		clearTimeout(_throttleTimeout);
		_throttleTimeout = void 0;
	}

	function scrollBy(el, x, y) {
		el.scrollLeft += x;
		el.scrollTop += y;
	}

	function clone(el) {
		var Polymer = window.Polymer;
		var $ = window.jQuery || window.Zepto;

		if (Polymer && Polymer.dom) {
			return Polymer.dom(el).cloneNode(true);
		} else if ($) {
			return $(el).clone(true)[0];
		} else {
			return el.cloneNode(true);
		}
	}

	function setRect(el, rect) {
		css(el, 'position', 'absolute');
		css(el, 'top', rect.top);
		css(el, 'left', rect.left);
		css(el, 'width', rect.width);
		css(el, 'height', rect.height);
	}

	function unsetRect(el) {
		css(el, 'position', '');
		css(el, 'top', '');
		css(el, 'left', '');
		css(el, 'width', '');
		css(el, 'height', '');
	}

	var expando = 'Sortable' + new Date().getTime();

	function AnimationStateManager() {
		var animationStates = [],
				animationCallbackId;
		return {
			captureAnimationState: function captureAnimationState() {
				animationStates = [];
				if (!this.options.animation) return;
				var children = [].slice.call(this.el.children);
				children.forEach(function (child) {
					if (css(child, 'display') === 'none' || child === Sortable.ghost) return;
					animationStates.push({
						target: child,
						rect: getRect(child)
					});
					var fromRect = getRect(child); // If animating: compensate for current animation

					if (child.thisAnimationDuration) {
						var childMatrix = matrix(child, true);

						if (childMatrix) {
							fromRect.top -= childMatrix.f;
							fromRect.left -= childMatrix.e;
						}
					}

					child.fromRect = fromRect;
				});
			},
			addAnimationState: function addAnimationState(state) {
				animationStates.push(state);
			},
			removeAnimationState: function removeAnimationState(target) {
				animationStates.splice(indexOfObject(animationStates, {
					target: target
				}), 1);
			},
			animateAll: function animateAll(callback) {
				var _this = this;

				if (!this.options.animation) {
					clearTimeout(animationCallbackId);
					if (typeof callback === 'function') callback();
					return;
				}

				var animating = false,
						animationTime = 0;
				animationStates.forEach(function (state) {
					var time = 0,
							target = state.target,
							fromRect = target.fromRect,
							toRect = getRect(target),
							prevFromRect = target.prevFromRect,
							prevToRect = target.prevToRect,
							animatingRect = state.rect,
							targetMatrix = matrix(target, true);

					if (targetMatrix) {
						// Compensate for current animation
						toRect.top -= targetMatrix.f;
						toRect.left -= targetMatrix.e;
					}

					target.toRect = toRect; // If element is scrolled out of view: Do not animate

					if ((isScrolledPast(target, toRect, 'bottom', 'top') || isScrolledPast(target, toRect, 'top', 'bottom') || isScrolledPast(target, toRect, 'right', 'left') || isScrolledPast(target, toRect, 'left', 'right')) && (isScrolledPast(target, animatingRect, 'bottom', 'top') || isScrolledPast(target, animatingRect, 'top', 'bottom') || isScrolledPast(target, animatingRect, 'right', 'left') || isScrolledPast(target, animatingRect, 'left', 'right')) && (isScrolledPast(target, fromRect, 'bottom', 'top') || isScrolledPast(target, fromRect, 'top', 'bottom') || isScrolledPast(target, fromRect, 'right', 'left') || isScrolledPast(target, fromRect, 'left', 'right'))) return;

					if (target.thisAnimationDuration) {
						// Could also check if animatingRect is between fromRect and toRect
						if (isRectEqual(prevFromRect, toRect) && !isRectEqual(fromRect, toRect) && // Make sure animatingRect is on line between toRect & fromRect
						(animatingRect.top - toRect.top) / (animatingRect.left - toRect.left) === (fromRect.top - toRect.top) / (fromRect.left - toRect.left)) {
							// If returning to same place as started from animation and on same axis
							time = calculateRealTime(animatingRect, prevFromRect, prevToRect, _this.options);
						}
					} // if fromRect != toRect: animate


					if (!isRectEqual(toRect, fromRect)) {
						target.prevFromRect = fromRect;
						target.prevToRect = toRect;

						if (!time) {
							time = _this.options.animation;
						}

						_this.animate(target, animatingRect, time);
					}

					if (time) {
						animating = true;
						animationTime = Math.max(animationTime, time);
						clearTimeout(target.animationResetTimer);
						target.animationResetTimer = setTimeout(function () {
							target.animationTime = 0;
							target.prevFromRect = null;
							target.fromRect = null;
							target.prevToRect = null;
							target.thisAnimationDuration = null;
						}, time);
						target.thisAnimationDuration = time;
					}
				});
				clearTimeout(animationCallbackId);

				if (!animating) {
					if (typeof callback === 'function') callback();
				} else {
					animationCallbackId = setTimeout(function () {
						if (typeof callback === 'function') callback();
					}, animationTime);
				}

				animationStates = [];
			},
			animate: function animate(target, prev, duration) {
				if (duration) {
					css(target, 'transition', '');
					css(target, 'transform', '');
					var currentRect = getRect(target),
							elMatrix = matrix(this.el),
							scaleX = elMatrix && elMatrix.a,
							scaleY = elMatrix && elMatrix.d,
							translateX = (prev.left - currentRect.left) / (scaleX || 1),
							translateY = (prev.top - currentRect.top) / (scaleY || 1);
					target.animatingX = !!translateX;
					target.animatingY = !!translateY;
					css(target, 'transform', 'translate3d(' + translateX + 'px,' + translateY + 'px,0)');
					repaint(target); // repaint

					css(target, 'transition', 'transform ' + duration + 'ms' + (this.options.easing ? ' ' + this.options.easing : ''));
					css(target, 'transform', 'translate3d(0,0,0)');
					typeof target.animated === 'number' && clearTimeout(target.animated);
					target.animated = setTimeout(function () {
						css(target, 'transition', '');
						css(target, 'transform', '');
						target.animated = false;
						target.animatingX = false;
						target.animatingY = false;
					}, duration);
				}
			}
		};
	}

	function repaint(target) {
		return target.offsetWidth;
	}

	function calculateRealTime(animatingRect, fromRect, toRect, options) {
		return Math.sqrt(Math.pow(fromRect.top - animatingRect.top, 2) + Math.pow(fromRect.left - animatingRect.left, 2)) / Math.sqrt(Math.pow(fromRect.top - toRect.top, 2) + Math.pow(fromRect.left - toRect.left, 2)) * options.animation;
	}

	var plugins = [];
	var defaults = {
		initializeByDefault: true
	};
	var PluginManager = {
		mount: function mount(plugin) {
			// Set default static properties
			for (var option in defaults) {
				if (defaults.hasOwnProperty(option) && !(option in plugin)) {
					plugin[option] = defaults[option];
				}
			}

			plugins.push(plugin);
		},
		pluginEvent: function pluginEvent(eventName, sortable, evt) {
			var _this = this;

			this.eventCanceled = false;
			var eventNameGlobal = eventName + 'Global';
			plugins.forEach(function (plugin) {
				if (!sortable[plugin.pluginName]) return; // Fire global events if it exists in this sortable

				if (sortable[plugin.pluginName][eventNameGlobal]) {
					_this.eventCanceled = !!sortable[plugin.pluginName][eventNameGlobal](_objectSpread({
						sortable: sortable
					}, evt));
				} // Only fire plugin event if plugin is enabled in this sortable,
				// and plugin has event defined


				if (sortable.options[plugin.pluginName] && sortable[plugin.pluginName][eventName]) {
					_this.eventCanceled = _this.eventCanceled || !!sortable[plugin.pluginName][eventName](_objectSpread({
						sortable: sortable
					}, evt));
				}
			});
		},
		initializePlugins: function initializePlugins(sortable, el, defaults) {
			plugins.forEach(function (plugin) {
				var pluginName = plugin.pluginName;
				if (!sortable.options[pluginName] && !plugin.initializeByDefault) return;
				var initialized = new plugin(sortable, el);
				initialized.sortable = sortable;
				sortable[pluginName] = initialized; // Add default options from plugin

				_extends(defaults, initialized.options);
			});

			for (var option in sortable.options) {
				if (!sortable.options.hasOwnProperty(option)) continue;
				var modified = this.modifyOption(sortable, option, sortable.options[option]);

				if (typeof modified !== 'undefined') {
					sortable.options[option] = modified;
				}
			}
		},
		getEventOptions: function getEventOptions(name, sortable) {
			var eventOptions = {};
			plugins.forEach(function (plugin) {
				if (typeof plugin.eventOptions !== 'function') return;

				_extends(eventOptions, plugin.eventOptions.call(sortable, name));
			});
			return eventOptions;
		},
		modifyOption: function modifyOption(sortable, name, value) {
			var modifiedValue;
			plugins.forEach(function (plugin) {
				// Plugin must exist on the Sortable
				if (!sortable[plugin.pluginName]) return; // If static option listener exists for this option, call in the context of the Sortable's instance of this plugin

				if (plugin.optionListeners && typeof plugin.optionListeners[name] === 'function') {
					modifiedValue = plugin.optionListeners[name].call(sortable[plugin.pluginName], value);
				}
			});
			return modifiedValue;
		}
	};

	function dispatchEvent(_ref) {
		var sortable = _ref.sortable,
				rootEl = _ref.rootEl,
				name = _ref.name,
				targetEl = _ref.targetEl,
				cloneEl = _ref.cloneEl,
				toEl = _ref.toEl,
				fromEl = _ref.fromEl,
				oldIndex = _ref.oldIndex,
				newIndex = _ref.newIndex,
				oldDraggableIndex = _ref.oldDraggableIndex,
				newDraggableIndex = _ref.newDraggableIndex,
				originalEvent = _ref.originalEvent,
				putSortable = _ref.putSortable,
				eventOptions = _ref.eventOptions;
		sortable = sortable || rootEl[expando];
		var evt,
				options = sortable.options,
				onName = 'on' + name.charAt(0).toUpperCase() + name.substr(1); // Support for new CustomEvent feature

		if (window.CustomEvent && !IE11OrLess && !Edge) {
			evt = new CustomEvent(name, {
				bubbles: true,
				cancelable: true
			});
		} else {
			evt = document.createEvent('Event');
			evt.initEvent(name, true, true);
		}

		evt.to = toEl || rootEl;
		evt.from = fromEl || rootEl;
		evt.item = targetEl || rootEl;
		evt.clone = cloneEl;
		evt.oldIndex = oldIndex;
		evt.newIndex = newIndex;
		evt.oldDraggableIndex = oldDraggableIndex;
		evt.newDraggableIndex = newDraggableIndex;
		evt.originalEvent = originalEvent;
		evt.pullMode = putSortable ? putSortable.lastPutMode : undefined;

		var allEventOptions = _objectSpread({}, eventOptions, PluginManager.getEventOptions(name, sortable));

		for (var option in allEventOptions) {
			evt[option] = allEventOptions[option];
		}

		if (rootEl) {
			rootEl.dispatchEvent(evt);
		}

		if (options[onName]) {
			options[onName].call(sortable, evt);
		}
	}

	var pluginEvent = function pluginEvent(eventName, sortable) {
		var _ref = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {},
				originalEvent = _ref.evt,
				data = _objectWithoutProperties(_ref, ["evt"]);

		PluginManager.pluginEvent.bind(Sortable)(eventName, sortable, _objectSpread({
			dragEl: dragEl,
			parentEl: parentEl,
			ghostEl: ghostEl,
			rootEl: rootEl,
			nextEl: nextEl,
			lastDownEl: lastDownEl,
			cloneEl: cloneEl,
			cloneHidden: cloneHidden,
			dragStarted: moved,
			putSortable: putSortable,
			activeSortable: Sortable.active,
			originalEvent: originalEvent,
			oldIndex: oldIndex,
			oldDraggableIndex: oldDraggableIndex,
			newIndex: newIndex,
			newDraggableIndex: newDraggableIndex,
			hideGhostForTarget: _hideGhostForTarget,
			unhideGhostForTarget: _unhideGhostForTarget,
			cloneNowHidden: function cloneNowHidden() {
				cloneHidden = true;
			},
			cloneNowShown: function cloneNowShown() {
				cloneHidden = false;
			},
			dispatchSortableEvent: function dispatchSortableEvent(name) {
				_dispatchEvent({
					sortable: sortable,
					name: name,
					originalEvent: originalEvent
				});
			}
		}, data));
	};

	function _dispatchEvent(info) {
		dispatchEvent(_objectSpread({
			putSortable: putSortable,
			cloneEl: cloneEl,
			targetEl: dragEl,
			rootEl: rootEl,
			oldIndex: oldIndex,
			oldDraggableIndex: oldDraggableIndex,
			newIndex: newIndex,
			newDraggableIndex: newDraggableIndex
		}, info));
	}

	if (typeof window === "undefined" || !window.document) {
		throw new Error("Sortable.js requires a window with a document");
	}

	var dragEl,
			parentEl,
			ghostEl,
			rootEl,
			nextEl,
			lastDownEl,
			cloneEl,
			cloneHidden,
			oldIndex,
			newIndex,
			oldDraggableIndex,
			newDraggableIndex,
			activeGroup,
			putSortable,
			awaitingDragStarted = false,
			ignoreNextClick = false,
			sortables = [],
			tapEvt,
			touchEvt,
			moved,
			lastTarget,
			lastDirection,
			pastFirstInvertThresh = false,
			isCircumstantialInvert = false,
			targetMoveDistance,
			// For positioning ghost absolutely
	ghostRelativeParent,
			ghostRelativeParentInitialScroll = [],
			// (left, top)
	_silent = false,
			savedInputChecked = [];
	/** @const */

	var PositionGhostAbsolutely = IOS,
			CSSFloatProperty = Edge || IE11OrLess ? 'cssFloat' : 'float',
			// This will not pass for IE9, because IE9 DnD only works on anchors
	supportDraggable = !ChromeForAndroid && !IOS && 'draggable' in document.createElement('div'),
			supportCssPointerEvents = function () {
		// false when <= IE11
		if (IE11OrLess) {
			return false;
		}

		var el = document.createElement('x');
		el.style.cssText = 'pointer-events:auto';
		return el.style.pointerEvents === 'auto';
	}(),
			_detectDirection = function _detectDirection(el, options) {
		var elCSS = css(el),
				elWidth = parseInt(elCSS.width) - parseInt(elCSS.paddingLeft) - parseInt(elCSS.paddingRight) - parseInt(elCSS.borderLeftWidth) - parseInt(elCSS.borderRightWidth),
				child1 = getChild(el, 0, options),
				child2 = getChild(el, 1, options),
				firstChildCSS = child1 && css(child1),
				secondChildCSS = child2 && css(child2),
				firstChildWidth = firstChildCSS && parseInt(firstChildCSS.marginLeft) + parseInt(firstChildCSS.marginRight) + getRect(child1).width,
				secondChildWidth = secondChildCSS && parseInt(secondChildCSS.marginLeft) + parseInt(secondChildCSS.marginRight) + getRect(child2).width;

		if (elCSS.display === 'flex') {
			return elCSS.flexDirection === 'column' || elCSS.flexDirection === 'column-reverse' ? 'vertical' : 'horizontal';
		}

		if (elCSS.display === 'grid') {
			return elCSS.gridTemplateColumns.split(' ').length <= 1 ? 'vertical' : 'horizontal';
		}

		if (child1 && firstChildCSS["float"] !== 'none') {
			var touchingSideChild2 = firstChildCSS["float"] === 'left' ? 'left' : 'right';
			return child2 && (secondChildCSS.clear === 'both' || secondChildCSS.clear === touchingSideChild2) ? 'vertical' : 'horizontal';
		}

		return child1 && (firstChildCSS.display === 'block' || firstChildCSS.display === 'flex' || firstChildCSS.display === 'table' || firstChildCSS.display === 'grid' || firstChildWidth >= elWidth && elCSS[CSSFloatProperty] === 'none' || child2 && elCSS[CSSFloatProperty] === 'none' && firstChildWidth + secondChildWidth > elWidth) ? 'vertical' : 'horizontal';
	},
			_dragElInRowColumn = function _dragElInRowColumn(dragRect, targetRect, vertical) {
		var dragElS1Opp = vertical ? dragRect.left : dragRect.top,
				dragElS2Opp = vertical ? dragRect.right : dragRect.bottom,
				dragElOppLength = vertical ? dragRect.width : dragRect.height,
				targetS1Opp = vertical ? targetRect.left : targetRect.top,
				targetS2Opp = vertical ? targetRect.right : targetRect.bottom,
				targetOppLength = vertical ? targetRect.width : targetRect.height;
		return dragElS1Opp === targetS1Opp || dragElS2Opp === targetS2Opp || dragElS1Opp + dragElOppLength / 2 === targetS1Opp + targetOppLength / 2;
	},

	/**
	 * Detects first nearest empty sortable to X and Y position using emptyInsertThreshold.
	 * @param  {Number} x      X position
	 * @param  {Number} y      Y position
	 * @return {HTMLElement}   Element of the first found nearest Sortable
	 */
	_detectNearestEmptySortable = function _detectNearestEmptySortable(x, y) {
		var ret;
		sortables.some(function (sortable) {
			if (lastChild(sortable)) return;
			var rect = getRect(sortable),
					threshold = sortable[expando].options.emptyInsertThreshold,
					insideHorizontally = x >= rect.left - threshold && x <= rect.right + threshold,
					insideVertically = y >= rect.top - threshold && y <= rect.bottom + threshold;

			if (threshold && insideHorizontally && insideVertically) {
				return ret = sortable;
			}
		});
		return ret;
	},
			_prepareGroup = function _prepareGroup(options) {
		function toFn(value, pull) {
			return function (to, from, dragEl, evt) {
				var sameGroup = to.options.group.name && from.options.group.name && to.options.group.name === from.options.group.name;

				if (value == null && (pull || sameGroup)) {
					// Default pull value
					// Default pull and put value if same group
					return true;
				} else if (value == null || value === false) {
					return false;
				} else if (pull && value === 'clone') {
					return value;
				} else if (typeof value === 'function') {
					return toFn(value(to, from, dragEl, evt), pull)(to, from, dragEl, evt);
				} else {
					var otherGroup = (pull ? to : from).options.group.name;
					return value === true || typeof value === 'string' && value === otherGroup || value.join && value.indexOf(otherGroup) > -1;
				}
			};
		}

		var group = {};
		var originalGroup = options.group;

		if (!originalGroup || _typeof(originalGroup) != 'object') {
			originalGroup = {
				name: originalGroup
			};
		}

		group.name = originalGroup.name;
		group.checkPull = toFn(originalGroup.pull, true);
		group.checkPut = toFn(originalGroup.put);
		group.revertClone = originalGroup.revertClone;
		options.group = group;
	},
			_hideGhostForTarget = function _hideGhostForTarget() {
		if (!supportCssPointerEvents && ghostEl) {
			css(ghostEl, 'display', 'none');
		}
	},
			_unhideGhostForTarget = function _unhideGhostForTarget() {
		if (!supportCssPointerEvents && ghostEl) {
			css(ghostEl, 'display', '');
		}
	}; // #1184 fix - Prevent click event on fallback if dragged but item not changed position


	document.addEventListener('click', function (evt) {
		if (ignoreNextClick) {
			evt.preventDefault();
			evt.stopPropagation && evt.stopPropagation();
			evt.stopImmediatePropagation && evt.stopImmediatePropagation();
			ignoreNextClick = false;
			return false;
		}
	}, true);

	var nearestEmptyInsertDetectEvent = function nearestEmptyInsertDetectEvent(evt) {
		if (dragEl) {
			evt = evt.touches ? evt.touches[0] : evt;

			var nearest = _detectNearestEmptySortable(evt.clientX, evt.clientY);

			if (nearest) {
				// Create imitation event
				var event = {};

				for (var i in evt) {
					if (evt.hasOwnProperty(i)) {
						event[i] = evt[i];
					}
				}

				event.target = event.rootEl = nearest;
				event.preventDefault = void 0;
				event.stopPropagation = void 0;

				nearest[expando]._onDragOver(event);
			}
		}
	};

	var _checkOutsideTargetEl = function _checkOutsideTargetEl(evt) {
		if (dragEl) {
			dragEl.parentNode[expando]._isOutsideThisEl(evt.target);
		}
	};
	/**
	 * @class  Sortable
	 * @param  {HTMLElement}  el
	 * @param  {Object}       [options]
	 */


	function Sortable(el, options) {
		if (!(el && el.nodeType && el.nodeType === 1)) {
			throw "Sortable: el must be an HTMLElement, not ".concat({}.toString.call(el));
		}

		this.el = el; // root element

		this.options = options = _extends({}, options); // Export instance

		el[expando] = this;
		var defaults = {
			group: null,
			sort: true,
			disabled: false,
			store: null,
			handle: null,
			draggable: /^[uo]l$/i.test(el.nodeName) ? '>li' : '>*',
			swapThreshold: 1,
			// percentage; 0 <= x <= 1
			invertSwap: false,
			// invert always
			invertedSwapThreshold: null,
			// will be set to same as swapThreshold if default
			removeCloneOnHide: true,
			direction: function direction() {
				return _detectDirection(el, this.options);
			},
			ghostClass: 'sortable-ghost',
			chosenClass: 'sortable-chosen',
			dragClass: 'sortable-drag',
			ignore: 'a, img',
			filter: null,
			preventOnFilter: true,
			animation: 0,
			easing: null,
			setData: function setData(dataTransfer, dragEl) {
				dataTransfer.setData('Text', dragEl.textContent);
			},
			dropBubble: false,
			dragoverBubble: false,
			dataIdAttr: 'data-id',
			delay: 0,
			delayOnTouchOnly: false,
			touchStartThreshold: (Number.parseInt ? Number : window).parseInt(window.devicePixelRatio, 10) || 1,
			forceFallback: false,
			fallbackClass: 'sortable-fallback',
			fallbackOnBody: false,
			fallbackTolerance: 0,
			fallbackOffset: {
				x: 0,
				y: 0
			},
			supportPointer: Sortable.supportPointer !== false && 'PointerEvent' in window,
			emptyInsertThreshold: 5
		};
		PluginManager.initializePlugins(this, el, defaults); // Set default options

		for (var name in defaults) {
			!(name in options) && (options[name] = defaults[name]);
		}

		_prepareGroup(options); // Bind all private methods


		for (var fn in this) {
			if (fn.charAt(0) === '_' && typeof this[fn] === 'function') {
				this[fn] = this[fn].bind(this);
			}
		} // Setup drag mode


		this.nativeDraggable = options.forceFallback ? false : supportDraggable;

		if (this.nativeDraggable) {
			// Touch start threshold cannot be greater than the native dragstart threshold
			this.options.touchStartThreshold = 1;
		} // Bind events


		if (options.supportPointer) {
			on(el, 'pointerdown', this._onTapStart);
		} else {
			on(el, 'mousedown', this._onTapStart);
			on(el, 'touchstart', this._onTapStart);
		}

		if (this.nativeDraggable) {
			on(el, 'dragover', this);
			on(el, 'dragenter', this);
		}

		sortables.push(this.el); // Restore sorting

		options.store && options.store.get && this.sort(options.store.get(this) || []); // Add animation state manager

		_extends(this, AnimationStateManager());
	}

	Sortable.prototype =
	/** @lends Sortable.prototype */
	{
		constructor: Sortable,
		_isOutsideThisEl: function _isOutsideThisEl(target) {
			if (!this.el.contains(target) && target !== this.el) {
				lastTarget = null;
			}
		},
		_getDirection: function _getDirection(evt, target) {
			return typeof this.options.direction === 'function' ? this.options.direction.call(this, evt, target, dragEl) : this.options.direction;
		},
		_onTapStart: function _onTapStart(
		/** Event|TouchEvent */
		evt) {
			if (!evt.cancelable) return;

			var _this = this,
					el = this.el,
					options = this.options,
					preventOnFilter = options.preventOnFilter,
					type = evt.type,
					touch = evt.touches && evt.touches[0],
					target = (touch || evt).target,
					originalTarget = evt.target.shadowRoot && (evt.path && evt.path[0] || evt.composedPath && evt.composedPath()[0]) || target,
					filter = options.filter;

			_saveInputCheckedState(el); // Don't trigger start event when an element is been dragged, otherwise the evt.oldindex always wrong when set option.group.


			if (dragEl) {
				return;
			}

			if (/mousedown|pointerdown/.test(type) && evt.button !== 0 || options.disabled) {
				return; // only left button and enabled
			} // cancel dnd if original target is content editable


			if (originalTarget.isContentEditable) {
				return;
			}

			target = closest(target, options.draggable, el, false);

			if (target && target.animated) {
				return;
			}

			if (lastDownEl === target) {
				// Ignoring duplicate down
				return;
			} // Get the index of the dragged element within its parent


			oldIndex = index(target);
			oldDraggableIndex = index(target, options.draggable); // Check filter

			if (typeof filter === 'function') {
				if (filter.call(this, evt, target, this)) {
					_dispatchEvent({
						sortable: _this,
						rootEl: originalTarget,
						name: 'filter',
						targetEl: target,
						toEl: el,
						fromEl: el
					});

					pluginEvent('filter', _this, {
						evt: evt
					});
					preventOnFilter && evt.cancelable && evt.preventDefault();
					return; // cancel dnd
				}
			} else if (filter) {
				filter = filter.split(',').some(function (criteria) {
					criteria = closest(originalTarget, criteria.trim(), el, false);

					if (criteria) {
						_dispatchEvent({
							sortable: _this,
							rootEl: criteria,
							name: 'filter',
							targetEl: target,
							fromEl: el,
							toEl: el
						});

						pluginEvent('filter', _this, {
							evt: evt
						});
						return true;
					}
				});

				if (filter) {
					preventOnFilter && evt.cancelable && evt.preventDefault();
					return; // cancel dnd
				}
			}

			if (options.handle && !closest(originalTarget, options.handle, el, false)) {
				return;
			} // Prepare


			this._prepareDragStart(evt, touch, target);
		},
		_prepareDragStart: function _prepareDragStart(
		/** Event */
		evt,
		/** Touch */
		touch,
		/** HTMLElement */
		target) {
			var _this = this,
					el = _this.el,
					options = _this.options,
					ownerDocument = el.ownerDocument,
					dragStartFn;

			if (target && !dragEl && target.parentNode === el) {
				rootEl = el;
				dragEl = target;
				parentEl = dragEl.parentNode;
				nextEl = dragEl.nextSibling;
				lastDownEl = target;
				activeGroup = options.group;
				Sortable.dragged = dragEl;
				tapEvt = {
					target: dragEl,
					clientX: (touch || evt).clientX,
					clientY: (touch || evt).clientY
				};
				this._lastX = (touch || evt).clientX;
				this._lastY = (touch || evt).clientY;
				dragEl.style['will-change'] = 'all';

				dragStartFn = function dragStartFn() {
					pluginEvent('delayEnded', _this, {
						evt: evt
					});

					if (Sortable.eventCanceled) {
						_this._onDrop();

						return;
					} // Delayed drag has been triggered
					// we can re-enable the events: touchmove/mousemove


					_this._disableDelayedDragEvents();

					if (!FireFox && _this.nativeDraggable) {
						dragEl.draggable = true;
					} // Bind the events: dragstart/dragend


					_this._triggerDragStart(evt, touch); // Drag start event


					_dispatchEvent({
						sortable: _this,
						name: 'choose',
						originalEvent: evt
					}); // Chosen item


					toggleClass(dragEl, options.chosenClass, true);
				}; // Disable "draggable"


				options.ignore.split(',').forEach(function (criteria) {
					find(dragEl, criteria.trim(), _disableDraggable);
				});
				on(ownerDocument, 'dragover', nearestEmptyInsertDetectEvent);
				on(ownerDocument, 'mousemove', nearestEmptyInsertDetectEvent);
				on(ownerDocument, 'touchmove', nearestEmptyInsertDetectEvent);
				on(ownerDocument, 'mouseup', _this._onDrop);
				on(ownerDocument, 'touchend', _this._onDrop);
				on(ownerDocument, 'touchcancel', _this._onDrop); // Make dragEl draggable (must be before delay for FireFox)

				if (FireFox && this.nativeDraggable) {
					this.options.touchStartThreshold = 4;
					dragEl.draggable = true;
				}

				pluginEvent('delayStart', this, {
					evt: evt
				}); // Delay is impossible for native DnD in Edge or IE

				if (options.delay && (!options.delayOnTouchOnly || touch) && (!this.nativeDraggable || !(Edge || IE11OrLess))) {
					if (Sortable.eventCanceled) {
						this._onDrop();

						return;
					} // If the user moves the pointer or let go the click or touch
					// before the delay has been reached:
					// disable the delayed drag


					on(ownerDocument, 'mouseup', _this._disableDelayedDrag);
					on(ownerDocument, 'touchend', _this._disableDelayedDrag);
					on(ownerDocument, 'touchcancel', _this._disableDelayedDrag);
					on(ownerDocument, 'mousemove', _this._delayedDragTouchMoveHandler);
					on(ownerDocument, 'touchmove', _this._delayedDragTouchMoveHandler);
					options.supportPointer && on(ownerDocument, 'pointermove', _this._delayedDragTouchMoveHandler);
					_this._dragStartTimer = setTimeout(dragStartFn, options.delay);
				} else {
					dragStartFn();
				}
			}
		},
		_delayedDragTouchMoveHandler: function _delayedDragTouchMoveHandler(
		/** TouchEvent|PointerEvent **/
		e) {
			var touch = e.touches ? e.touches[0] : e;

			if (Math.max(Math.abs(touch.clientX - this._lastX), Math.abs(touch.clientY - this._lastY)) >= Math.floor(this.options.touchStartThreshold / (this.nativeDraggable && window.devicePixelRatio || 1))) {
				this._disableDelayedDrag();
			}
		},
		_disableDelayedDrag: function _disableDelayedDrag() {
			dragEl && _disableDraggable(dragEl);
			clearTimeout(this._dragStartTimer);

			this._disableDelayedDragEvents();
		},
		_disableDelayedDragEvents: function _disableDelayedDragEvents() {
			var ownerDocument = this.el.ownerDocument;
			off(ownerDocument, 'mouseup', this._disableDelayedDrag);
			off(ownerDocument, 'touchend', this._disableDelayedDrag);
			off(ownerDocument, 'touchcancel', this._disableDelayedDrag);
			off(ownerDocument, 'mousemove', this._delayedDragTouchMoveHandler);
			off(ownerDocument, 'touchmove', this._delayedDragTouchMoveHandler);
			off(ownerDocument, 'pointermove', this._delayedDragTouchMoveHandler);
		},
		_triggerDragStart: function _triggerDragStart(
		/** Event */
		evt,
		/** Touch */
		touch) {
			touch = touch || evt.pointerType == 'touch' && evt;

			if (!this.nativeDraggable || touch) {
				if (this.options.supportPointer) {
					on(document, 'pointermove', this._onTouchMove);
				} else if (touch) {
					on(document, 'touchmove', this._onTouchMove);
				} else {
					on(document, 'mousemove', this._onTouchMove);
				}
			} else {
				on(dragEl, 'dragend', this);
				on(rootEl, 'dragstart', this._onDragStart);
			}

			try {
				if (document.selection) {
					// Timeout neccessary for IE9
					_nextTick(function () {
						document.selection.empty();
					});
				} else {
					window.getSelection().removeAllRanges();
				}
			} catch (err) {}
		},
		_dragStarted: function _dragStarted(fallback, evt) {

			awaitingDragStarted = false;

			if (rootEl && dragEl) {
				pluginEvent('dragStarted', this, {
					evt: evt
				});

				if (this.nativeDraggable) {
					on(document, 'dragover', _checkOutsideTargetEl);
				}

				var options = this.options; // Apply effect

				!fallback && toggleClass(dragEl, options.dragClass, false);
				toggleClass(dragEl, options.ghostClass, true);
				Sortable.active = this;
				fallback && this._appendGhost(); // Drag start event

				_dispatchEvent({
					sortable: this,
					name: 'start',
					originalEvent: evt
				});
			} else {
				this._nulling();
			}
		},
		_emulateDragOver: function _emulateDragOver() {
			if (touchEvt) {
				this._lastX = touchEvt.clientX;
				this._lastY = touchEvt.clientY;

				_hideGhostForTarget();

				var target = document.elementFromPoint(touchEvt.clientX, touchEvt.clientY);
				var parent = target;

				while (target && target.shadowRoot) {
					target = target.shadowRoot.elementFromPoint(touchEvt.clientX, touchEvt.clientY);
					if (target === parent) break;
					parent = target;
				}

				dragEl.parentNode[expando]._isOutsideThisEl(target);

				if (parent) {
					do {
						if (parent[expando]) {
							var inserted = void 0;
							inserted = parent[expando]._onDragOver({
								clientX: touchEvt.clientX,
								clientY: touchEvt.clientY,
								target: target,
								rootEl: parent
							});

							if (inserted && !this.options.dragoverBubble) {
								break;
							}
						}

						target = parent; // store last element
					}
					/* jshint boss:true */
					while (parent = parent.parentNode);
				}

				_unhideGhostForTarget();
			}
		},
		_onTouchMove: function _onTouchMove(
		/**TouchEvent*/
		evt) {
			if (tapEvt) {
				var options = this.options,
						fallbackTolerance = options.fallbackTolerance,
						fallbackOffset = options.fallbackOffset,
						touch = evt.touches ? evt.touches[0] : evt,
						ghostMatrix = ghostEl && matrix(ghostEl),
						scaleX = ghostEl && ghostMatrix && ghostMatrix.a,
						scaleY = ghostEl && ghostMatrix && ghostMatrix.d,
						relativeScrollOffset = PositionGhostAbsolutely && ghostRelativeParent && getRelativeScrollOffset(ghostRelativeParent),
						dx = (touch.clientX - tapEvt.clientX + fallbackOffset.x) / (scaleX || 1) + (relativeScrollOffset ? relativeScrollOffset[0] - ghostRelativeParentInitialScroll[0] : 0) / (scaleX || 1),
						dy = (touch.clientY - tapEvt.clientY + fallbackOffset.y) / (scaleY || 1) + (relativeScrollOffset ? relativeScrollOffset[1] - ghostRelativeParentInitialScroll[1] : 0) / (scaleY || 1),
						translate3d = evt.touches ? 'translate3d(' + dx + 'px,' + dy + 'px,0)' : 'translate(' + dx + 'px,' + dy + 'px)'; // only set the status to dragging, when we are actually dragging

				if (!Sortable.active && !awaitingDragStarted) {
					if (fallbackTolerance && Math.max(Math.abs(touch.clientX - this._lastX), Math.abs(touch.clientY - this._lastY)) < fallbackTolerance) {
						return;
					}

					this._onDragStart(evt, true);
				}

				touchEvt = touch;
				css(ghostEl, 'webkitTransform', translate3d);
				css(ghostEl, 'mozTransform', translate3d);
				css(ghostEl, 'msTransform', translate3d);
				css(ghostEl, 'transform', translate3d);
				evt.cancelable && evt.preventDefault();
			}
		},
		_appendGhost: function _appendGhost() {
			// Bug if using scale(): https://stackoverflow.com/questions/2637058
			// Not being adjusted for
			if (!ghostEl) {
				var container = this.options.fallbackOnBody ? document.body : rootEl,
						rect = getRect(dragEl, true, PositionGhostAbsolutely, true, container),
						options = this.options; // Position absolutely

				if (PositionGhostAbsolutely) {
					// Get relatively positioned parent
					ghostRelativeParent = container;

					while (css(ghostRelativeParent, 'position') === 'static' && css(ghostRelativeParent, 'transform') === 'none' && ghostRelativeParent !== document) {
						ghostRelativeParent = ghostRelativeParent.parentNode;
					}

					if (ghostRelativeParent !== document.body && ghostRelativeParent !== document.documentElement) {
						if (ghostRelativeParent === document) ghostRelativeParent = getWindowScrollingElement();
						rect.top += ghostRelativeParent.scrollTop;
						rect.left += ghostRelativeParent.scrollLeft;
					} else {
						ghostRelativeParent = getWindowScrollingElement();
					}

					ghostRelativeParentInitialScroll = getRelativeScrollOffset(ghostRelativeParent);
				}

				ghostEl = dragEl.cloneNode(true);
				toggleClass(ghostEl, options.ghostClass, false);
				toggleClass(ghostEl, options.fallbackClass, true);
				toggleClass(ghostEl, options.dragClass, true);
				css(ghostEl, 'transition', '');
				css(ghostEl, 'transform', '');
				css(ghostEl, 'box-sizing', 'border-box');
				css(ghostEl, 'margin', 0);
				css(ghostEl, 'top', rect.top);
				css(ghostEl, 'left', rect.left);
				css(ghostEl, 'width', rect.width);
				css(ghostEl, 'height', rect.height);
				css(ghostEl, 'opacity', '0.8');
				css(ghostEl, 'position', PositionGhostAbsolutely ? 'absolute' : 'fixed');
				css(ghostEl, 'zIndex', '100000');
				css(ghostEl, 'pointerEvents', 'none');
				Sortable.ghost = ghostEl;
				container.appendChild(ghostEl);
			}
		},
		_onDragStart: function _onDragStart(
		/**Event*/
		evt,
		/**boolean*/
		fallback) {
			var _this = this;

			var dataTransfer = evt.dataTransfer;
			var options = _this.options;
			pluginEvent('dragStart', this, {
				evt: evt
			});

			if (Sortable.eventCanceled) {
				this._onDrop();

				return;
			}

			pluginEvent('setupClone', this);

			if (!Sortable.eventCanceled) {
				cloneEl = clone(dragEl);
				cloneEl.draggable = false;
				cloneEl.style['will-change'] = '';

				this._hideClone();

				toggleClass(cloneEl, this.options.chosenClass, false);
				Sortable.clone = cloneEl;
			} // #1143: IFrame support workaround


			_this.cloneId = _nextTick(function () {
				pluginEvent('clone', _this);
				if (Sortable.eventCanceled) return;

				if (!_this.options.removeCloneOnHide) {
					rootEl.insertBefore(cloneEl, dragEl);
				}

				_this._hideClone();

				_dispatchEvent({
					sortable: _this,
					name: 'clone'
				});
			});
			!fallback && toggleClass(dragEl, options.dragClass, true); // Set proper drop events

			if (fallback) {
				ignoreNextClick = true;
				_this._loopId = setInterval(_this._emulateDragOver, 50);
			} else {
				// Undo what was set in _prepareDragStart before drag started
				off(document, 'mouseup', _this._onDrop);
				off(document, 'touchend', _this._onDrop);
				off(document, 'touchcancel', _this._onDrop);

				if (dataTransfer) {
					dataTransfer.effectAllowed = 'move';
					options.setData && options.setData.call(_this, dataTransfer, dragEl);
				}

				on(document, 'drop', _this); // #1276 fix:

				css(dragEl, 'transform', 'translateZ(0)');
			}

			awaitingDragStarted = true;
			_this._dragStartId = _nextTick(_this._dragStarted.bind(_this, fallback, evt));
			on(document, 'selectstart', _this);
			moved = true;

			if (Safari) {
				css(document.body, 'user-select', 'none');
			}
		},
		// Returns true - if no further action is needed (either inserted or another condition)
		_onDragOver: function _onDragOver(
		/**Event*/
		evt) {
			var el = this.el,
					target = evt.target,
					dragRect,
					targetRect,
					revert,
					options = this.options,
					group = options.group,
					activeSortable = Sortable.active,
					isOwner = activeGroup === group,
					canSort = options.sort,
					fromSortable = putSortable || activeSortable,
					vertical,
					_this = this,
					completedFired = false;

			if (_silent) return;

			function dragOverEvent(name, extra) {
				pluginEvent(name, _this, _objectSpread({
					evt: evt,
					isOwner: isOwner,
					axis: vertical ? 'vertical' : 'horizontal',
					revert: revert,
					dragRect: dragRect,
					targetRect: targetRect,
					canSort: canSort,
					fromSortable: fromSortable,
					target: target,
					completed: completed,
					onMove: function onMove(target, after) {
						return _onMove(rootEl, el, dragEl, dragRect, target, getRect(target), evt, after);
					},
					changed: changed
				}, extra));
			} // Capture animation state


			function capture() {
				dragOverEvent('dragOverAnimationCapture');

				_this.captureAnimationState();

				if (_this !== fromSortable) {
					fromSortable.captureAnimationState();
				}
			} // Return invocation when dragEl is inserted (or completed)


			function completed(insertion) {
				dragOverEvent('dragOverCompleted', {
					insertion: insertion
				});

				if (insertion) {
					// Clones must be hidden before folding animation to capture dragRectAbsolute properly
					if (isOwner) {
						activeSortable._hideClone();
					} else {
						activeSortable._showClone(_this);
					}

					if (_this !== fromSortable) {
						// Set ghost class to new sortable's ghost class
						toggleClass(dragEl, putSortable ? putSortable.options.ghostClass : activeSortable.options.ghostClass, false);
						toggleClass(dragEl, options.ghostClass, true);
					}

					if (putSortable !== _this && _this !== Sortable.active) {
						putSortable = _this;
					} else if (_this === Sortable.active && putSortable) {
						putSortable = null;
					} // Animation


					if (fromSortable === _this) {
						_this._ignoreWhileAnimating = target;
					}

					_this.animateAll(function () {
						dragOverEvent('dragOverAnimationComplete');
						_this._ignoreWhileAnimating = null;
					});

					if (_this !== fromSortable) {
						fromSortable.animateAll();
						fromSortable._ignoreWhileAnimating = null;
					}
				} // Null lastTarget if it is not inside a previously swapped element


				if (target === dragEl && !dragEl.animated || target === el && !target.animated) {
					lastTarget = null;
				} // no bubbling and not fallback


				if (!options.dragoverBubble && !evt.rootEl && target !== document) {
					dragEl.parentNode[expando]._isOutsideThisEl(evt.target); // Do not detect for empty insert if already inserted


					!insertion && nearestEmptyInsertDetectEvent(evt);
				}

				!options.dragoverBubble && evt.stopPropagation && evt.stopPropagation();
				return completedFired = true;
			} // Call when dragEl has been inserted


			function changed() {
				newIndex = index(dragEl);
				newDraggableIndex = index(dragEl, options.draggable);

				_dispatchEvent({
					sortable: _this,
					name: 'change',
					toEl: el,
					newIndex: newIndex,
					newDraggableIndex: newDraggableIndex,
					originalEvent: evt
				});
			}

			if (evt.preventDefault !== void 0) {
				evt.cancelable && evt.preventDefault();
			}

			target = closest(target, options.draggable, el, true);
			dragOverEvent('dragOver');
			if (Sortable.eventCanceled) return completedFired;

			if (dragEl.contains(evt.target) || target.animated && target.animatingX && target.animatingY || _this._ignoreWhileAnimating === target) {
				return completed(false);
			}

			ignoreNextClick = false;

			if (activeSortable && !options.disabled && (isOwner ? canSort || (revert = !rootEl.contains(dragEl)) // Reverting item into the original list
			: putSortable === this || (this.lastPutMode = activeGroup.checkPull(this, activeSortable, dragEl, evt)) && group.checkPut(this, activeSortable, dragEl, evt))) {
				vertical = this._getDirection(evt, target) === 'vertical';
				dragRect = getRect(dragEl);
				dragOverEvent('dragOverValid');
				if (Sortable.eventCanceled) return completedFired;

				if (revert) {
					parentEl = rootEl; // actualization

					capture();

					this._hideClone();

					dragOverEvent('revert');

					if (!Sortable.eventCanceled) {
						if (nextEl) {
							rootEl.insertBefore(dragEl, nextEl);
						} else {
							rootEl.appendChild(dragEl);
						}
					}

					return completed(true);
				}

				var elLastChild = lastChild(el, options.draggable);

				if (!elLastChild || _ghostIsLast(evt, vertical, this) && !elLastChild.animated) {
					// If already at end of list: Do not insert
					if (elLastChild === dragEl) {
						return completed(false);
					} // assign target only if condition is true


					if (elLastChild && el === evt.target) {
						target = elLastChild;
					}

					if (target) {
						targetRect = getRect(target);
					}

					if (_onMove(rootEl, el, dragEl, dragRect, target, targetRect, evt, !!target) !== false) {
						capture();
						el.appendChild(dragEl);
						parentEl = el; // actualization

						changed();
						return completed(true);
					}
				} else if (target.parentNode === el) {
					targetRect = getRect(target);
					var direction = 0,
							targetBeforeFirstSwap,
							differentLevel = dragEl.parentNode !== el,
							differentRowCol = !_dragElInRowColumn(dragEl.animated && dragEl.toRect || dragRect, target.animated && target.toRect || targetRect, vertical),
							side1 = vertical ? 'top' : 'left',
							scrolledPastTop = isScrolledPast(target, null, 'top', 'top') || isScrolledPast(dragEl, null, 'top', 'top'),
							scrollBefore = scrolledPastTop ? scrolledPastTop.scrollTop : void 0;

					if (lastTarget !== target) {
						targetBeforeFirstSwap = targetRect[side1];
						pastFirstInvertThresh = false;
						isCircumstantialInvert = !differentRowCol && options.invertSwap || differentLevel;
					}

					direction = _getSwapDirection(evt, target, targetRect, vertical, differentRowCol ? 1 : options.swapThreshold, options.invertedSwapThreshold == null ? options.swapThreshold : options.invertedSwapThreshold, isCircumstantialInvert, lastTarget === target);
					var sibling;

					if (direction !== 0) {
						// Check if target is beside dragEl in respective direction (ignoring hidden elements)
						var dragIndex = index(dragEl);

						do {
							dragIndex -= direction;
							sibling = parentEl.children[dragIndex];
						} while (sibling && (css(sibling, 'display') === 'none' || sibling === ghostEl));
					} // If dragEl is already beside target: Do not insert


					if (direction === 0 || sibling === target) {
						return completed(false);
					}

					lastTarget = target;
					lastDirection = direction;
					var nextSibling = target.nextElementSibling,
							after = false;
					after = direction === 1;

					var moveVector = _onMove(rootEl, el, dragEl, dragRect, target, targetRect, evt, after);

					if (moveVector !== false) {
						if (moveVector === 1 || moveVector === -1) {
							after = moveVector === 1;
						}

						_silent = true;
						setTimeout(_unsilent, 30);
						capture();

						if (after && !nextSibling) {
							el.appendChild(dragEl);
						} else {
							target.parentNode.insertBefore(dragEl, after ? nextSibling : target);
						} // Undo chrome's scroll adjustment (has no effect on other browsers)


						if (scrolledPastTop) {
							scrollBy(scrolledPastTop, 0, scrollBefore - scrolledPastTop.scrollTop);
						}

						parentEl = dragEl.parentNode; // actualization
						// must be done before animation

						if (targetBeforeFirstSwap !== undefined && !isCircumstantialInvert) {
							targetMoveDistance = Math.abs(targetBeforeFirstSwap - getRect(target)[side1]);
						}

						changed();
						return completed(true);
					}
				}

				if (el.contains(dragEl)) {
					return completed(false);
				}
			}

			return false;
		},
		_ignoreWhileAnimating: null,
		_offMoveEvents: function _offMoveEvents() {
			off(document, 'mousemove', this._onTouchMove);
			off(document, 'touchmove', this._onTouchMove);
			off(document, 'pointermove', this._onTouchMove);
			off(document, 'dragover', nearestEmptyInsertDetectEvent);
			off(document, 'mousemove', nearestEmptyInsertDetectEvent);
			off(document, 'touchmove', nearestEmptyInsertDetectEvent);
		},
		_offUpEvents: function _offUpEvents() {
			var ownerDocument = this.el.ownerDocument;
			off(ownerDocument, 'mouseup', this._onDrop);
			off(ownerDocument, 'touchend', this._onDrop);
			off(ownerDocument, 'pointerup', this._onDrop);
			off(ownerDocument, 'touchcancel', this._onDrop);
			off(document, 'selectstart', this);
		},
		_onDrop: function _onDrop(
		/**Event*/
		evt) {
			var el = this.el,
					options = this.options; // Get the index of the dragged element within its parent

			newIndex = index(dragEl);
			newDraggableIndex = index(dragEl, options.draggable);
			pluginEvent('drop', this, {
				evt: evt
			}); // Get again after plugin event

			newIndex = index(dragEl);
			newDraggableIndex = index(dragEl, options.draggable);

			if (Sortable.eventCanceled) {
				this._nulling();

				return;
			}

			awaitingDragStarted = false;
			isCircumstantialInvert = false;
			pastFirstInvertThresh = false;
			clearInterval(this._loopId);
			clearTimeout(this._dragStartTimer);

			_cancelNextTick(this.cloneId);

			_cancelNextTick(this._dragStartId); // Unbind events


			if (this.nativeDraggable) {
				off(document, 'drop', this);
				off(el, 'dragstart', this._onDragStart);
			}

			this._offMoveEvents();

			this._offUpEvents();

			if (Safari) {
				css(document.body, 'user-select', '');
			}

			if (evt) {
				if (moved) {
					evt.cancelable && evt.preventDefault();
					!options.dropBubble && evt.stopPropagation();
				}

				ghostEl && ghostEl.parentNode && ghostEl.parentNode.removeChild(ghostEl);

				if (rootEl === parentEl || putSortable && putSortable.lastPutMode !== 'clone') {
					// Remove clone(s)
					cloneEl && cloneEl.parentNode && cloneEl.parentNode.removeChild(cloneEl);
				}

				if (dragEl) {
					if (this.nativeDraggable) {
						off(dragEl, 'dragend', this);
					}

					_disableDraggable(dragEl);

					dragEl.style['will-change'] = ''; // Remove classes
					// ghostClass is added in dragStarted

					if (moved && !awaitingDragStarted) {
						toggleClass(dragEl, putSortable ? putSortable.options.ghostClass : this.options.ghostClass, false);
					}

					toggleClass(dragEl, this.options.chosenClass, false); // Drag stop event

					_dispatchEvent({
						sortable: this,
						name: 'unchoose',
						toEl: parentEl,
						newIndex: null,
						newDraggableIndex: null,
						originalEvent: evt
					});

					if (rootEl !== parentEl) {
						if (newIndex >= 0) {
							// Add event
							_dispatchEvent({
								rootEl: parentEl,
								name: 'add',
								toEl: parentEl,
								fromEl: rootEl,
								originalEvent: evt
							}); // Remove event


							_dispatchEvent({
								sortable: this,
								name: 'remove',
								toEl: parentEl,
								originalEvent: evt
							}); // drag from one list and drop into another


							_dispatchEvent({
								rootEl: parentEl,
								name: 'sort',
								toEl: parentEl,
								fromEl: rootEl,
								originalEvent: evt
							});

							_dispatchEvent({
								sortable: this,
								name: 'sort',
								toEl: parentEl,
								originalEvent: evt
							});
						}

						putSortable && putSortable.save();
					} else {
						if (newIndex !== oldIndex) {
							if (newIndex >= 0) {
								// drag & drop within the same list
								_dispatchEvent({
									sortable: this,
									name: 'update',
									toEl: parentEl,
									originalEvent: evt
								});

								_dispatchEvent({
									sortable: this,
									name: 'sort',
									toEl: parentEl,
									originalEvent: evt
								});
							}
						}
					}

					if (Sortable.active) {
						/* jshint eqnull:true */
						if (newIndex == null || newIndex === -1) {
							newIndex = oldIndex;
							newDraggableIndex = oldDraggableIndex;
						}

						_dispatchEvent({
							sortable: this,
							name: 'end',
							toEl: parentEl,
							originalEvent: evt
						}); // Save sorting


						this.save();
					}
				}
			}

			this._nulling();
		},
		_nulling: function _nulling() {
			pluginEvent('nulling', this);
			rootEl = dragEl = parentEl = ghostEl = nextEl = cloneEl = lastDownEl = cloneHidden = tapEvt = touchEvt = moved = newIndex = newDraggableIndex = oldIndex = oldDraggableIndex = lastTarget = lastDirection = putSortable = activeGroup = Sortable.dragged = Sortable.ghost = Sortable.clone = Sortable.active = null;
			savedInputChecked.forEach(function (el) {
				el.checked = true;
			});
			savedInputChecked.length = 0;
		},
		handleEvent: function handleEvent(
		/**Event*/
		evt) {
			switch (evt.type) {
				case 'drop':
				case 'dragend':
					this._onDrop(evt);

					break;

				case 'dragenter':
				case 'dragover':
					if (dragEl) {
						this._onDragOver(evt);

						_globalDragOver(evt);
					}

					break;

				case 'selectstart':
					evt.preventDefault();
					break;
			}
		},

		/**
		 * Serializes the item into an array of string.
		 * @returns {String[]}
		 */
		toArray: function toArray() {
			var order = [],
					el,
					children = this.el.children,
					i = 0,
					n = children.length,
					options = this.options;

			for (; i < n; i++) {
				el = children[i];

				if (closest(el, options.draggable, this.el, false)) {
					order.push(el.getAttribute(options.dataIdAttr) || _generateId(el));
				}
			}

			return order;
		},

		/**
		 * Sorts the elements according to the array.
		 * @param  {String[]}  order  order of the items
		 */
		sort: function sort(order) {
			var items = {},
					rootEl = this.el;
			this.toArray().forEach(function (id, i) {
				var el = rootEl.children[i];

				if (closest(el, this.options.draggable, rootEl, false)) {
					items[id] = el;
				}
			}, this);
			order.forEach(function (id) {
				if (items[id]) {
					rootEl.removeChild(items[id]);
					rootEl.appendChild(items[id]);
				}
			});
		},

		/**
		 * Save the current sorting
		 */
		save: function save() {
			var store = this.options.store;
			store && store.set && store.set(this);
		},

		/**
		 * For each element in the set, get the first element that matches the selector by testing the element itself and traversing up through its ancestors in the DOM tree.
		 * @param   {HTMLElement}  el
		 * @param   {String}       [selector]  default: options.draggable
		 * @returns {HTMLElement|null}
		 */
		closest: function closest$1(el, selector) {
			return closest(el, selector || this.options.draggable, this.el, false);
		},

		/**
		 * Set/get option
		 * @param   {string} name
		 * @param   {*}      [value]
		 * @returns {*}
		 */
		option: function option(name, value) {
			var options = this.options;

			if (value === void 0) {
				return options[name];
			} else {
				var modifiedValue = PluginManager.modifyOption(this, name, value);

				if (typeof modifiedValue !== 'undefined') {
					options[name] = modifiedValue;
				} else {
					options[name] = value;
				}

				if (name === 'group') {
					_prepareGroup(options);
				}
			}
		},

		/**
		 * Destroy
		 */
		destroy: function destroy() {
			pluginEvent('destroy', this);
			var el = this.el;
			el[expando] = null;
			off(el, 'mousedown', this._onTapStart);
			off(el, 'touchstart', this._onTapStart);
			off(el, 'pointerdown', this._onTapStart);

			if (this.nativeDraggable) {
				off(el, 'dragover', this);
				off(el, 'dragenter', this);
			} // Remove draggable attributes


			Array.prototype.forEach.call(el.querySelectorAll('[draggable]'), function (el) {
				el.removeAttribute('draggable');
			});

			this._onDrop();

			sortables.splice(sortables.indexOf(this.el), 1);
			this.el = el = null;
		},
		_hideClone: function _hideClone() {
			if (!cloneHidden) {
				pluginEvent('hideClone', this);
				if (Sortable.eventCanceled) return;
				css(cloneEl, 'display', 'none');

				if (this.options.removeCloneOnHide && cloneEl.parentNode) {
					cloneEl.parentNode.removeChild(cloneEl);
				}

				cloneHidden = true;
			}
		},
		_showClone: function _showClone(putSortable) {
			if (putSortable.lastPutMode !== 'clone') {
				this._hideClone();

				return;
			}

			if (cloneHidden) {
				pluginEvent('showClone', this);
				if (Sortable.eventCanceled) return; // show clone at dragEl or original position

				if (rootEl.contains(dragEl) && !this.options.group.revertClone) {
					rootEl.insertBefore(cloneEl, dragEl);
				} else if (nextEl) {
					rootEl.insertBefore(cloneEl, nextEl);
				} else {
					rootEl.appendChild(cloneEl);
				}

				if (this.options.group.revertClone) {
					this._animate(dragEl, cloneEl);
				}

				css(cloneEl, 'display', '');
				cloneHidden = false;
			}
		}
	};

	function _globalDragOver(
	/**Event*/
	evt) {
		if (evt.dataTransfer) {
			evt.dataTransfer.dropEffect = 'move';
		}

		evt.cancelable && evt.preventDefault();
	}

	function _onMove(fromEl, toEl, dragEl, dragRect, targetEl, targetRect, originalEvent, willInsertAfter) {
		var evt,
				sortable = fromEl[expando],
				onMoveFn = sortable.options.onMove,
				retVal; // Support for new CustomEvent feature

		if (window.CustomEvent && !IE11OrLess && !Edge) {
			evt = new CustomEvent('move', {
				bubbles: true,
				cancelable: true
			});
		} else {
			evt = document.createEvent('Event');
			evt.initEvent('move', true, true);
		}

		evt.to = toEl;
		evt.from = fromEl;
		evt.dragged = dragEl;
		evt.draggedRect = dragRect;
		evt.related = targetEl || toEl;
		evt.relatedRect = targetRect || getRect(toEl);
		evt.willInsertAfter = willInsertAfter;
		evt.originalEvent = originalEvent;
		fromEl.dispatchEvent(evt);

		if (onMoveFn) {
			retVal = onMoveFn.call(sortable, evt, originalEvent);
		}

		return retVal;
	}

	function _disableDraggable(el) {
		el.draggable = false;
	}

	function _unsilent() {
		_silent = false;
	}

	function _ghostIsLast(evt, vertical, sortable) {
		var rect = getRect(lastChild(sortable.el, sortable.options.draggable));
		var spacer = 10;
		return vertical ? evt.clientX > rect.right + spacer || evt.clientX <= rect.right && evt.clientY > rect.bottom && evt.clientX >= rect.left : evt.clientX > rect.right && evt.clientY > rect.top || evt.clientX <= rect.right && evt.clientY > rect.bottom + spacer;
	}

	function _getSwapDirection(evt, target, targetRect, vertical, swapThreshold, invertedSwapThreshold, invertSwap, isLastTarget) {
		var mouseOnAxis = vertical ? evt.clientY : evt.clientX,
				targetLength = vertical ? targetRect.height : targetRect.width,
				targetS1 = vertical ? targetRect.top : targetRect.left,
				targetS2 = vertical ? targetRect.bottom : targetRect.right,
				invert = false;

		if (!invertSwap) {
			// Never invert or create dragEl shadow when target movemenet causes mouse to move past the end of regular swapThreshold
			if (isLastTarget && targetMoveDistance < targetLength * swapThreshold) {
				// multiplied only by swapThreshold because mouse will already be inside target by (1 - threshold) * targetLength / 2
				// check if past first invert threshold on side opposite of lastDirection
				if (!pastFirstInvertThresh && (lastDirection === 1 ? mouseOnAxis > targetS1 + targetLength * invertedSwapThreshold / 2 : mouseOnAxis < targetS2 - targetLength * invertedSwapThreshold / 2)) {
					// past first invert threshold, do not restrict inverted threshold to dragEl shadow
					pastFirstInvertThresh = true;
				}

				if (!pastFirstInvertThresh) {
					// dragEl shadow (target move distance shadow)
					if (lastDirection === 1 ? mouseOnAxis < targetS1 + targetMoveDistance // over dragEl shadow
					: mouseOnAxis > targetS2 - targetMoveDistance) {
						return -lastDirection;
					}
				} else {
					invert = true;
				}
			} else {
				// Regular
				if (mouseOnAxis > targetS1 + targetLength * (1 - swapThreshold) / 2 && mouseOnAxis < targetS2 - targetLength * (1 - swapThreshold) / 2) {
					return _getInsertDirection(target);
				}
			}
		}

		invert = invert || invertSwap;

		if (invert) {
			// Invert of regular
			if (mouseOnAxis < targetS1 + targetLength * invertedSwapThreshold / 2 || mouseOnAxis > targetS2 - targetLength * invertedSwapThreshold / 2) {
				return mouseOnAxis > targetS1 + targetLength / 2 ? 1 : -1;
			}
		}

		return 0;
	}
	/**
	 * Gets the direction dragEl must be swapped relative to target in order to make it
	 * seem that dragEl has been "inserted" into that element's position
	 * @param  {HTMLElement} target       The target whose position dragEl is being inserted at
	 * @return {Number}                   Direction dragEl must be swapped
	 */


	function _getInsertDirection(target) {
		if (index(dragEl) < index(target)) {
			return 1;
		} else {
			return -1;
		}
	}
	/**
	 * Generate id
	 * @param   {HTMLElement} el
	 * @returns {String}
	 * @private
	 */


	function _generateId(el) {
		var str = el.tagName + el.className + el.src + el.href + el.textContent,
				i = str.length,
				sum = 0;

		while (i--) {
			sum += str.charCodeAt(i);
		}

		return sum.toString(36);
	}

	function _saveInputCheckedState(root) {
		savedInputChecked.length = 0;
		var inputs = root.getElementsByTagName('input');
		var idx = inputs.length;

		while (idx--) {
			var _el = inputs[idx];
			_el.checked && savedInputChecked.push(_el);
		}
	}

	function _nextTick(fn) {
		return setTimeout(fn, 0);
	}

	function _cancelNextTick(id) {
		return clearTimeout(id);
	} // Fixed #973:


	on(document, 'touchmove', function (evt) {
		if ((Sortable.active || awaitingDragStarted) && evt.cancelable) {
			evt.preventDefault();
		}
	}); // Export utils

	Sortable.utils = {
		on: on,
		off: off,
		css: css,
		find: find,
		is: function is(el, selector) {
			return !!closest(el, selector, el, false);
		},
		extend: extend,
		throttle: throttle,
		closest: closest,
		toggleClass: toggleClass,
		clone: clone,
		index: index,
		nextTick: _nextTick,
		cancelNextTick: _cancelNextTick,
		detectDirection: _detectDirection,
		getChild: getChild
	};
	/**
	 * Mount a plugin to Sortable
	 * @param  {...SortablePlugin|SortablePlugin[]} plugins       Plugins being mounted
	 */

	Sortable.mount = function () {
		for (var _len = arguments.length, plugins = new Array(_len), _key = 0; _key < _len; _key++) {
			plugins[_key] = arguments[_key];
		}

		if (plugins[0].constructor === Array) plugins = plugins[0];
		plugins.forEach(function (plugin) {
			if (!plugin.prototype || !plugin.prototype.constructor) {
				throw "Sortable: Mounted plugin must be a constructor function, not ".concat({}.toString.call(el));
			}

			if (plugin.utils) Sortable.utils = _objectSpread({}, Sortable.utils, plugin.utils);
			PluginManager.mount(plugin);
		});
	};
	/**
	 * Create sortable instance
	 * @param {HTMLElement}  el
	 * @param {Object}      [options]
	 */


	Sortable.create = function (el, options) {
		return new Sortable(el, options);
	}; // Export


	Sortable.version = version;

	var autoScrolls = [],
			scrollEl,
			scrollRootEl,
			scrolling = false,
			lastAutoScrollX,
			lastAutoScrollY,
			touchEvt$1,
			pointerElemChangedInterval;

	function AutoScrollPlugin() {
		function AutoScroll() {
			this.options = {
				scroll: true,
				scrollSensitivity: 30,
				scrollSpeed: 10,
				bubbleScroll: true
			}; // Bind all private methods

			for (var fn in this) {
				if (fn.charAt(0) === '_' && typeof this[fn] === 'function') {
					this[fn] = this[fn].bind(this);
				}
			}
		}

		AutoScroll.prototype = {
			dragStarted: function dragStarted(_ref) {
				var originalEvent = _ref.originalEvent;

				if (this.sortable.nativeDraggable) {
					on(document, 'dragover', this._handleAutoScroll);
				} else {
					if (this.sortable.options.supportPointer) {
						on(document, 'pointermove', this._handleFallbackAutoScroll);
					} else if (originalEvent.touches) {
						on(document, 'touchmove', this._handleFallbackAutoScroll);
					} else {
						on(document, 'mousemove', this._handleFallbackAutoScroll);
					}
				}
			},
			dragOverCompleted: function dragOverCompleted(_ref2) {
				var originalEvent = _ref2.originalEvent;

				// For when bubbling is canceled and using fallback (fallback 'touchmove' always reached)
				if (!this.sortable.options.dragOverBubble && !originalEvent.rootEl) {
					this._handleAutoScroll(originalEvent);
				}
			},
			drop: function drop() {
				if (this.sortable.nativeDraggable) {
					off(document, 'dragover', this._handleAutoScroll);
				} else {
					off(document, 'pointermove', this._handleFallbackAutoScroll);
					off(document, 'touchmove', this._handleFallbackAutoScroll);
					off(document, 'mousemove', this._handleFallbackAutoScroll);
				}

				clearPointerElemChangedInterval();
				clearAutoScrolls();
				cancelThrottle();
			},
			nulling: function nulling() {
				touchEvt$1 = scrollRootEl = scrollEl = scrolling = pointerElemChangedInterval = lastAutoScrollX = lastAutoScrollY = null;
				autoScrolls.length = 0;
			},
			_handleFallbackAutoScroll: function _handleFallbackAutoScroll(evt) {
				this._handleAutoScroll(evt, true);
			},
			_handleAutoScroll: function _handleAutoScroll(evt, fallback) {
				var _this = this;

				var x = evt.clientX,
						y = evt.clientY,
						elem = document.elementFromPoint(x, y);
				touchEvt$1 = evt; // IE does not seem to have native autoscroll,
				// Edge's autoscroll seems too conditional,
				// MACOS Safari does not have autoscroll,
				// Firefox and Chrome are good

				if (fallback || Edge || IE11OrLess || Safari) {
					autoScroll(evt, this.options, elem, fallback); // Listener for pointer element change

					var ogElemScroller = getParentAutoScrollElement(elem, true);

					if (scrolling && (!pointerElemChangedInterval || x !== lastAutoScrollX || y !== lastAutoScrollY)) {
						pointerElemChangedInterval && clearPointerElemChangedInterval(); // Detect for pointer elem change, emulating native DnD behaviour

						pointerElemChangedInterval = setInterval(function () {
							var newElem = getParentAutoScrollElement(document.elementFromPoint(x, y), true);

							if (newElem !== ogElemScroller) {
								ogElemScroller = newElem;
								clearAutoScrolls();
							}

							autoScroll(evt, _this.options, newElem, fallback);
						}, 10);
						lastAutoScrollX = x;
						lastAutoScrollY = y;
					}
				} else {
					// if DnD is enabled (and browser has good autoscrolling), first autoscroll will already scroll, so get parent autoscroll of first autoscroll
					if (!this.sortable.options.bubbleScroll || getParentAutoScrollElement(elem, true) === getWindowScrollingElement()) {
						clearAutoScrolls();
						return;
					}

					autoScroll(evt, this.options, getParentAutoScrollElement(elem, false), false);
				}
			}
		};
		return _extends(AutoScroll, {
			pluginName: 'scroll',
			initializeByDefault: true
		});
	}

	function clearAutoScrolls() {
		autoScrolls.forEach(function (autoScroll) {
			clearInterval(autoScroll.pid);
		});
		autoScrolls = [];
	}

	function clearPointerElemChangedInterval() {
		clearInterval(pointerElemChangedInterval);
	}

	var autoScroll = throttle(function (evt, options, rootEl, isFallback) {
		// Bug: https://bugzilla.mozilla.org/show_bug.cgi?id=505521
		if (!options.scroll) return;
		var sens = options.scrollSensitivity,
				speed = options.scrollSpeed,
				winScroller = getWindowScrollingElement();
		var scrollThisInstance = false,
				scrollCustomFn; // New scroll root, set scrollEl

		if (scrollRootEl !== rootEl) {
			scrollRootEl = rootEl;
			clearAutoScrolls();
			scrollEl = options.scroll;
			scrollCustomFn = options.scrollFn;

			if (scrollEl === true) {
				scrollEl = getParentAutoScrollElement(rootEl, true);
			}
		}

		var layersOut = 0;
		var currentParent = scrollEl;

		do {
			var el = currentParent,
					rect = getRect(el),
					top = rect.top,
					bottom = rect.bottom,
					left = rect.left,
					right = rect.right,
					width = rect.width,
					height = rect.height,
					canScrollX = void 0,
					canScrollY = void 0,
					scrollWidth = el.scrollWidth,
					scrollHeight = el.scrollHeight,
					elCSS = css(el),
					scrollPosX = el.scrollLeft,
					scrollPosY = el.scrollTop;

			if (el === winScroller) {
				canScrollX = width < scrollWidth && (elCSS.overflowX === 'auto' || elCSS.overflowX === 'scroll' || elCSS.overflowX === 'visible');
				canScrollY = height < scrollHeight && (elCSS.overflowY === 'auto' || elCSS.overflowY === 'scroll' || elCSS.overflowY === 'visible');
			} else {
				canScrollX = width < scrollWidth && (elCSS.overflowX === 'auto' || elCSS.overflowX === 'scroll');
				canScrollY = height < scrollHeight && (elCSS.overflowY === 'auto' || elCSS.overflowY === 'scroll');
			}

			var vx = canScrollX && (Math.abs(right - evt.clientX) <= sens && scrollPosX + width < scrollWidth) - (Math.abs(left - evt.clientX) <= sens && !!scrollPosX);
			var vy = canScrollY && (Math.abs(bottom - evt.clientY) <= sens && scrollPosY + height < scrollHeight) - (Math.abs(top - evt.clientY) <= sens && !!scrollPosY);

			if (!autoScrolls[layersOut]) {
				for (var i = 0; i <= layersOut; i++) {
					if (!autoScrolls[i]) {
						autoScrolls[i] = {};
					}
				}
			}

			if (autoScrolls[layersOut].vx != vx || autoScrolls[layersOut].vy != vy || autoScrolls[layersOut].el !== el) {
				autoScrolls[layersOut].el = el;
				autoScrolls[layersOut].vx = vx;
				autoScrolls[layersOut].vy = vy;
				clearInterval(autoScrolls[layersOut].pid);

				if (vx != 0 || vy != 0) {
					scrollThisInstance = true;
					/* jshint loopfunc:true */

					autoScrolls[layersOut].pid = setInterval(function () {
						// emulate drag over during autoscroll (fallback), emulating native DnD behaviour
						if (isFallback && this.layer === 0) {
							Sortable.active._onTouchMove(touchEvt$1); // To move ghost if it is positioned absolutely

						}

						var scrollOffsetY = autoScrolls[this.layer].vy ? autoScrolls[this.layer].vy * speed : 0;
						var scrollOffsetX = autoScrolls[this.layer].vx ? autoScrolls[this.layer].vx * speed : 0;

						if (typeof scrollCustomFn === 'function') {
							if (scrollCustomFn.call(Sortable.dragged.parentNode[expando], scrollOffsetX, scrollOffsetY, evt, touchEvt$1, autoScrolls[this.layer].el) !== 'continue') {
								return;
							}
						}

						scrollBy(autoScrolls[this.layer].el, scrollOffsetX, scrollOffsetY);
					}.bind({
						layer: layersOut
					}), 24);
				}
			}

			layersOut++;
		} while (options.bubbleScroll && currentParent !== winScroller && (currentParent = getParentAutoScrollElement(currentParent, false)));

		scrolling = scrollThisInstance; // in case another function catches scrolling as false in between when it is not
	}, 30);

	var drop = function drop(_ref) {
		var originalEvent = _ref.originalEvent,
				putSortable = _ref.putSortable,
				dragEl = _ref.dragEl,
				activeSortable = _ref.activeSortable,
				dispatchSortableEvent = _ref.dispatchSortableEvent,
				hideGhostForTarget = _ref.hideGhostForTarget,
				unhideGhostForTarget = _ref.unhideGhostForTarget;
		var toSortable = putSortable || activeSortable;
		hideGhostForTarget();
		var target = document.elementFromPoint(originalEvent.clientX, originalEvent.clientY);
		unhideGhostForTarget();

		if (toSortable && !toSortable.el.contains(target)) {
			dispatchSortableEvent('spill');
			this.onSpill(dragEl);
		}
	};

	function Revert() {}

	Revert.prototype = {
		startIndex: null,
		dragStart: function dragStart(_ref2) {
			var oldDraggableIndex = _ref2.oldDraggableIndex;
			this.startIndex = oldDraggableIndex;
		},
		onSpill: function onSpill(dragEl) {
			this.sortable.captureAnimationState();
			var nextSibling = getChild(this.sortable.el, this.startIndex, this.sortable.options);

			if (nextSibling) {
				this.sortable.el.insertBefore(dragEl, nextSibling);
			} else {
				this.sortable.el.appendChild(dragEl);
			}

			this.sortable.animateAll();
		},
		drop: drop
	};

	_extends(Revert, {
		pluginName: 'revertOnSpill'
	});

	function Remove() {}

	Remove.prototype = {
		onSpill: function onSpill(dragEl) {
			this.sortable.captureAnimationState();
			dragEl.parentNode && dragEl.parentNode.removeChild(dragEl);
			this.sortable.animateAll();
		},
		drop: drop
	};

	_extends(Remove, {
		pluginName: 'removeOnSpill'
	});

	var lastSwapEl;

	function SwapPlugin() {
		function Swap() {
			this.options = {
				swapClass: 'sortable-swap-highlight'
			};
		}

		Swap.prototype = {
			dragStart: function dragStart(_ref) {
				var dragEl = _ref.dragEl;
				lastSwapEl = dragEl;
			},
			dragOverValid: function dragOverValid(_ref2) {
				var completed = _ref2.completed,
						target = _ref2.target,
						onMove = _ref2.onMove,
						activeSortable = _ref2.activeSortable,
						changed = _ref2.changed;
				if (!activeSortable.options.swap) return;
				var el = this.sortable.el,
						options = this.sortable.options;

				if (target && target !== el) {
					var prevSwapEl = lastSwapEl;

					if (onMove(target) !== false) {
						toggleClass(target, options.swapClass, true);
						lastSwapEl = target;
					} else {
						lastSwapEl = null;
					}

					if (prevSwapEl && prevSwapEl !== lastSwapEl) {
						toggleClass(prevSwapEl, options.swapClass, false);
					}
				}

				changed();
				return completed(true);
			},
			drop: function drop(_ref3) {
				var activeSortable = _ref3.activeSortable,
						putSortable = _ref3.putSortable,
						dragEl = _ref3.dragEl;
				var toSortable = putSortable || this.sortable;
				var options = this.sortable.options;
				lastSwapEl && toggleClass(lastSwapEl, options.swapClass, false);

				if (lastSwapEl && (options.swap || putSortable && putSortable.options.swap)) {
					if (dragEl !== lastSwapEl) {
						toSortable.captureAnimationState();
						if (toSortable !== activeSortable) activeSortable.captureAnimationState();
						swapNodes(dragEl, lastSwapEl);
						toSortable.animateAll();
						if (toSortable !== activeSortable) activeSortable.animateAll();
					}
				}
			},
			nulling: function nulling() {
				lastSwapEl = null;
			}
		};
		return _extends(Swap, {
			pluginName: 'swap',
			eventOptions: function eventOptions() {
				return {
					swapItem: lastSwapEl
				};
			}
		});
	}

	function swapNodes(n1, n2) {
		var p1 = n1.parentNode,
				p2 = n2.parentNode,
				i1,
				i2;
		if (!p1 || !p2 || p1.isEqualNode(n2) || p2.isEqualNode(n1)) return;
		i1 = index(n1);
		i2 = index(n2);

		if (p1.isEqualNode(p2) && i1 < i2) {
			i2++;
		}

		p1.insertBefore(n2, p1.children[i1]);
		p2.insertBefore(n1, p2.children[i2]);
	}

	var multiDragElements = [],
			multiDragClones = [],
			lastMultiDragSelect,
			// for selection with modifier key down (SHIFT)
	multiDragSortable,
			initialFolding = false,
			// Initial multi-drag fold when drag started
	folding = false,
			// Folding any other time
	dragStarted = false,
			dragEl$1,
			clonesFromRect,
			clonesHidden;

	function MultiDragPlugin() {
		function MultiDrag(sortable) {
			// Bind all private methods
			for (var fn in this) {
				if (fn.charAt(0) === '_' && typeof this[fn] === 'function') {
					this[fn] = this[fn].bind(this);
				}
			}

			if (sortable.options.supportPointer) {
				on(document, 'pointerup', this._deselectMultiDrag);
			} else {
				on(document, 'mouseup', this._deselectMultiDrag);
				on(document, 'touchend', this._deselectMultiDrag);
			}

			on(document, 'keydown', this._checkKeyDown);
			on(document, 'keyup', this._checkKeyUp);
			this.options = {
				selectedClass: 'sortable-selected',
				multiDragKey: null,
				setData: function setData(dataTransfer, dragEl) {
					var data = '';

					if (multiDragElements.length && multiDragSortable === sortable) {
						multiDragElements.forEach(function (multiDragElement, i) {
							data += (!i ? '' : ', ') + multiDragElement.textContent;
						});
					} else {
						data = dragEl.textContent;
					}

					dataTransfer.setData('Text', data);
				}
			};
		}

		MultiDrag.prototype = {
			multiDragKeyDown: false,
			isMultiDrag: false,
			delayStartGlobal: function delayStartGlobal(_ref) {
				var dragged = _ref.dragEl;
				dragEl$1 = dragged;
			},
			delayEnded: function delayEnded() {
				this.isMultiDrag = ~multiDragElements.indexOf(dragEl$1);
			},
			setupClone: function setupClone(_ref2) {
				var sortable = _ref2.sortable;
				if (!this.isMultiDrag) return;

				for (var _i = 0; _i < multiDragElements.length; _i++) {
					multiDragClones.push(clone(multiDragElements[_i]));
					multiDragClones[_i].sortableIndex = multiDragElements[_i].sortableIndex;
					multiDragClones[_i].draggable = false;
					multiDragClones[_i].style['will-change'] = '';
					toggleClass(multiDragClones[_i], sortable.options.selectedClass, false);
					multiDragElements[_i] === dragEl$1 && toggleClass(multiDragClones[_i], sortable.options.chosenClass, false);
				}

				sortable._hideClone();

				return true;
			},
			clone: function clone(_ref3) {
				var sortable = _ref3.sortable,
						rootEl = _ref3.rootEl,
						dispatchSortableEvent = _ref3.dispatchSortableEvent;
				if (!this.isMultiDrag) return;

				if (!sortable.options.removeCloneOnHide) {
					if (multiDragElements.length && multiDragSortable === sortable) {
						insertMultiDragClones(true, rootEl);
						dispatchSortableEvent('clone');
						return true;
					}
				}
			},
			showClone: function showClone(_ref4) {
				var cloneNowShown = _ref4.cloneNowShown,
						rootEl = _ref4.rootEl;
				if (!this.isMultiDrag) return;
				insertMultiDragClones(false, rootEl);
				multiDragClones.forEach(function (clone) {
					css(clone, 'display', '');
				});
				cloneNowShown();
				clonesHidden = false;
				return true;
			},
			hideClone: function hideClone(_ref5) {
				var sortable = _ref5.sortable,
						cloneNowHidden = _ref5.cloneNowHidden;
				if (!this.isMultiDrag) return;
				multiDragClones.forEach(function (clone) {
					css(clone, 'display', 'none');

					if (sortable.options.removeCloneOnHide && clone.parentNode) {
						clone.parentNode.removeChild(clone);
					}
				});
				cloneNowHidden();
				clonesHidden = true;
				return true;
			},
			dragStartGlobal: function dragStartGlobal(_ref6) {
				var sortable = _ref6.sortable;

				if (!this.isMultiDrag && multiDragSortable) {
					multiDragSortable.multiDrag._deselectMultiDrag();
				}

				multiDragElements.forEach(function (multiDragElement) {
					multiDragElement.sortableIndex = index(multiDragElement);
				}); // Sort multi-drag elements

				multiDragElements = multiDragElements.sort(function (a, b) {
					return a.sortableIndex - b.sortableIndex;
				});
				dragStarted = true;
			},
			dragStarted: function dragStarted(_ref7) {
				var sortable = _ref7.sortable;
				if (!this.isMultiDrag) return;

				if (sortable.options.sort) {
					// Capture rects,
					// hide multi drag elements (by positioning them absolute),
					// set multi drag elements rects to dragRect,
					// show multi drag elements,
					// animate to rects,
					// unset rects & remove from DOM
					sortable.captureAnimationState();

					if (sortable.options.animation) {
						multiDragElements.forEach(function (multiDragElement) {
							if (multiDragElement === dragEl$1) return;
							css(multiDragElement, 'position', 'absolute');
						});
						var dragRect = getRect(dragEl$1, false, true, true);
						multiDragElements.forEach(function (multiDragElement) {
							if (multiDragElement === dragEl$1) return;
							setRect(multiDragElement, dragRect);
						});
						folding = true;
						initialFolding = true;
					}
				}

				sortable.animateAll(function () {
					folding = false;
					initialFolding = false;

					if (sortable.options.animation) {
						multiDragElements.forEach(function (multiDragElement) {
							unsetRect(multiDragElement);
						});
					} // Remove all auxiliary multidrag items from el, if sorting enabled


					if (sortable.options.sort) {
						removeMultiDragElements();
					}
				});
			},
			dragOver: function dragOver(_ref8) {
				var target = _ref8.target,
						completed = _ref8.completed;

				if (folding && ~multiDragElements.indexOf(target)) {
					return completed(false);
				}
			},
			revert: function revert(_ref9) {
				var fromSortable = _ref9.fromSortable,
						rootEl = _ref9.rootEl,
						sortable = _ref9.sortable,
						dragRect = _ref9.dragRect;

				if (multiDragElements.length > 1) {
					// Setup unfold animation
					multiDragElements.forEach(function (multiDragElement) {
						sortable.addAnimationState({
							target: multiDragElement,
							rect: folding ? getRect(multiDragElement) : dragRect
						});
						unsetRect(multiDragElement);
						multiDragElement.fromRect = dragRect;
						fromSortable.removeAnimationState(multiDragElement);
					});
					folding = false;
					insertMultiDragElements(!sortable.options.removeCloneOnHide, rootEl);
				}
			},
			dragOverCompleted: function dragOverCompleted(_ref10) {
				var sortable = _ref10.sortable,
						isOwner = _ref10.isOwner,
						insertion = _ref10.insertion,
						activeSortable = _ref10.activeSortable,
						parentEl = _ref10.parentEl,
						putSortable = _ref10.putSortable;
				var options = sortable.options;

				if (insertion) {
					// Clones must be hidden before folding animation to capture dragRectAbsolute properly
					if (isOwner) {
						activeSortable._hideClone();
					}

					initialFolding = false; // If leaving sort:false root, or already folding - Fold to new location

					if (options.animation && multiDragElements.length > 1 && (folding || !isOwner && !activeSortable.options.sort && !putSortable)) {
						// Fold: Set all multi drag elements's rects to dragEl's rect when multi-drag elements are invisible
						var dragRectAbsolute = getRect(dragEl$1, false, true, true);
						multiDragElements.forEach(function (multiDragElement) {
							if (multiDragElement === dragEl$1) return;
							setRect(multiDragElement, dragRectAbsolute); // Move element(s) to end of parentEl so that it does not interfere with multi-drag clones insertion if they are inserted
							// while folding, and so that we can capture them again because old sortable will no longer be fromSortable

							parentEl.appendChild(multiDragElement);
						});
						folding = true;
					} // Clones must be shown (and check to remove multi drags) after folding when interfering multiDragElements are moved out


					if (!isOwner) {
						// Only remove if not folding (folding will remove them anyways)
						if (!folding) {
							removeMultiDragElements();
						}

						if (multiDragElements.length > 1) {
							var clonesHiddenBefore = clonesHidden;

							activeSortable._showClone(sortable); // Unfold animation for clones if showing from hidden


							if (activeSortable.options.animation && !clonesHidden && clonesHiddenBefore) {
								multiDragClones.forEach(function (clone) {
									activeSortable.addAnimationState({
										target: clone,
										rect: clonesFromRect
									});
									clone.fromRect = clonesFromRect;
									clone.thisAnimationDuration = null;
								});
							}
						} else {
							activeSortable._showClone(sortable);
						}
					}
				}
			},
			dragOverAnimationCapture: function dragOverAnimationCapture(_ref11) {
				var dragRect = _ref11.dragRect,
						isOwner = _ref11.isOwner,
						activeSortable = _ref11.activeSortable;
				multiDragElements.forEach(function (multiDragElement) {
					multiDragElement.thisAnimationDuration = null;
				});

				if (activeSortable.options.animation && !isOwner && activeSortable.multiDrag.isMultiDrag) {
					clonesFromRect = _extends({}, dragRect);
					var dragMatrix = matrix(dragEl$1, true);
					clonesFromRect.top -= dragMatrix.f;
					clonesFromRect.left -= dragMatrix.e;
				}
			},
			dragOverAnimationComplete: function dragOverAnimationComplete() {
				if (folding) {
					folding = false;
					removeMultiDragElements();
				}
			},
			drop: function drop(_ref12) {
				var evt = _ref12.originalEvent,
						rootEl = _ref12.rootEl,
						parentEl = _ref12.parentEl,
						sortable = _ref12.sortable,
						dispatchSortableEvent = _ref12.dispatchSortableEvent,
						oldIndex = _ref12.oldIndex,
						putSortable = _ref12.putSortable;
				var toSortable = putSortable || this.sortable;
				if (!evt) return;
				var options = sortable.options,
						children = parentEl.children; // Multi-drag selection

				if (!dragStarted) {
					if (options.multiDragKey && !this.multiDragKeyDown) {
						this._deselectMultiDrag();
					}

					toggleClass(dragEl$1, options.selectedClass, !~multiDragElements.indexOf(dragEl$1));

					if (!~multiDragElements.indexOf(dragEl$1)) {
						multiDragElements.push(dragEl$1);
						dispatchEvent({
							sortable: sortable,
							rootEl: rootEl,
							name: 'select',
							targetEl: dragEl$1,
							originalEvt: evt
						}); // Modifier activated, select from last to dragEl

						if ((!options.multiDragKey || this.multiDragKeyDown) && evt.shiftKey && lastMultiDragSelect && sortable.el.contains(lastMultiDragSelect)) {
							var lastIndex = index(lastMultiDragSelect),
									currentIndex = index(dragEl$1);

							if (~lastIndex && ~currentIndex && lastIndex !== currentIndex) {
								// Must include lastMultiDragSelect (select it), in case modified selection from no selection
								// (but previous selection existed)
								var n, _i2;

								if (currentIndex > lastIndex) {
									_i2 = lastIndex;
									n = currentIndex;
								} else {
									_i2 = currentIndex;
									n = lastIndex + 1;
								}

								for (; _i2 < n; _i2++) {
									if (~multiDragElements.indexOf(children[_i2])) continue;
									toggleClass(children[_i2], options.selectedClass, true);
									multiDragElements.push(children[_i2]);
									dispatchEvent({
										sortable: sortable,
										rootEl: rootEl,
										name: 'select',
										targetEl: children[_i2],
										originalEvt: evt
									});
								}
							}
						} else {
							lastMultiDragSelect = dragEl$1;
						}

						multiDragSortable = toSortable;
					} else {
						multiDragElements.splice(multiDragElements.indexOf(dragEl$1), 1);
						lastMultiDragSelect = null;
						dispatchEvent({
							sortable: sortable,
							rootEl: rootEl,
							name: 'deselect',
							targetEl: dragEl$1,
							originalEvt: evt
						});
					}
				} // Multi-drag drop


				if (dragStarted && this.isMultiDrag) {
					// Do not "unfold" after around dragEl if reverted
					if ((parentEl[expando].options.sort || parentEl !== rootEl) && multiDragElements.length > 1) {
						var dragRect = getRect(dragEl$1),
								multiDragIndex = index(dragEl$1, ':not(.' + this.options.selectedClass + ')');
						if (!initialFolding && options.animation) dragEl$1.thisAnimationDuration = null;
						toSortable.captureAnimationState();

						if (!initialFolding) {
							if (options.animation) {
								dragEl$1.fromRect = dragRect;
								multiDragElements.forEach(function (multiDragElement) {
									multiDragElement.thisAnimationDuration = null;

									if (multiDragElement !== dragEl$1) {
										var rect = folding ? getRect(multiDragElement) : dragRect;
										multiDragElement.fromRect = rect; // Prepare unfold animation

										toSortable.addAnimationState({
											target: multiDragElement,
											rect: rect
										});
									}
								});
							} // Multi drag elements are not necessarily removed from the DOM on drop, so to reinsert
							// properly they must all be removed


							removeMultiDragElements();
							multiDragElements.forEach(function (multiDragElement) {
								if (children[multiDragIndex]) {
									parentEl.insertBefore(multiDragElement, children[multiDragIndex]);
								} else {
									parentEl.appendChild(multiDragElement);
								}

								multiDragIndex++;
							}); // If initial folding is done, the elements may have changed position because they are now
							// unfolding around dragEl, even though dragEl may not have his index changed, so update event
							// must be fired here as Sortable will not.

							if (oldIndex === index(dragEl$1)) {
								var update = false;
								multiDragElements.forEach(function (multiDragElement) {
									if (multiDragElement.sortableIndex !== index(multiDragElement)) {
										update = true;
										return;
									}
								});

								if (update) {
									dispatchSortableEvent('update');
								}
							}
						} // Must be done after capturing individual rects (scroll bar)


						multiDragElements.forEach(function (multiDragElement) {
							unsetRect(multiDragElement);
						});
						toSortable.animateAll();
					}

					multiDragSortable = toSortable;
				} // Remove clones if necessary


				if (rootEl === parentEl || putSortable && putSortable.lastPutMode !== 'clone') {
					multiDragClones.forEach(function (clone) {
						clone.parentNode && clone.parentNode.removeChild(clone);
					});
				}
			},
			nullingGlobal: function nullingGlobal() {
				this.isMultiDrag = dragStarted = false;
				multiDragClones.length = 0;
			},
			destroy: function destroy() {
				this._deselectMultiDrag();

				off(document, 'pointerup', this._deselectMultiDrag);
				off(document, 'mouseup', this._deselectMultiDrag);
				off(document, 'touchend', this._deselectMultiDrag);
				off(document, 'keydown', this._checkKeyDown);
				off(document, 'keyup', this._checkKeyUp);
			},
			_deselectMultiDrag: function _deselectMultiDrag(evt) {
				if (dragStarted) return; // Only deselect if selection is in this sortable

				if (multiDragSortable !== this.sortable) return; // Only deselect if target is not item in this sortable

				if (evt && closest(evt.target, this.sortable.options.draggable, this.sortable.el, false)) return; // Only deselect if left click

				if (evt && evt.button !== 0) return;

				while (multiDragElements.length) {
					var el = multiDragElements[0];
					toggleClass(el, this.sortable.options.selectedClass, false);
					multiDragElements.shift();
					dispatchEvent({
						sortable: this.sortable,
						rootEl: this.sortable.el,
						name: 'deselect',
						targetEl: el,
						originalEvt: evt
					});
				}
			},
			_checkKeyDown: function _checkKeyDown(evt) {
				if (evt.key === this.sortable.options.multiDragKey) {
					this.multiDragKeyDown = true;
				}
			},
			_checkKeyUp: function _checkKeyUp(evt) {
				if (evt.key === this.sortable.options.multiDragKey) {
					this.multiDragKeyDown = false;
				}
			}
		};
		return _extends(MultiDrag, {
			// Static methods & properties
			pluginName: 'multiDrag',
			utils: {
				/**
				 * Selects the provided multi-drag item
				 * @param  {HTMLElement} el    The element to be selected
				 */
				select: function select(el) {
					var sortable = el.parentNode[expando];
					if (!sortable || !sortable.options.multiDrag || ~multiDragElements.indexOf(el)) return;

					if (multiDragSortable && multiDragSortable !== sortable) {
						multiDragSortable.multiDrag._deselectMultiDrag();

						multiDragSortable = sortable;
					}

					toggleClass(el, sortable.options.selectedClass, true);
					multiDragElements.push(el);
				},

				/**
				 * Deselects the provided multi-drag item
				 * @param  {HTMLElement} el    The element to be deselected
				 */
				deselect: function deselect(el) {
					var sortable = el.parentNode[expando],
							index = multiDragElements.indexOf(el);
					if (!sortable || !sortable.options.multiDrag || !~index) return;
					toggleClass(el, sortable.options.selectedClass, false);
					multiDragElements.splice(index, 1);
				}
			},
			eventOptions: function eventOptions() {
				var _this = this;

				var oldIndicies = [],
						newIndicies = [];
				multiDragElements.forEach(function (multiDragElement) {
					oldIndicies.push({
						multiDragElement: multiDragElement,
						index: multiDragElement.sortableIndex
					}); // multiDragElements will already be sorted if folding

					var newIndex;

					if (folding && multiDragElement !== dragEl$1) {
						newIndex = -1;
					} else if (folding) {
						newIndex = index(multiDragElement, ':not(.' + _this.options.selectedClass + ')');
					} else {
						newIndex = index(multiDragElement);
					}

					newIndicies.push({
						multiDragElement: multiDragElement,
						index: newIndex
					});
				});
				return {
					items: _toConsumableArray(multiDragElements),
					clones: [].concat(multiDragClones),
					oldIndicies: oldIndicies,
					newIndicies: newIndicies
				};
			},
			optionListeners: {
				multiDragKey: function multiDragKey(key) {
					key = key.toLowerCase();

					if (key === 'ctrl') {
						key = 'Control';
					} else if (key.length > 1) {
						key = key.charAt(0).toUpperCase() + key.substr(1);
					}

					return key;
				}
			}
		});
	}

	function insertMultiDragElements(clonesInserted, rootEl) {
		multiDragElements.forEach(function (multiDragElement) {
			var target = rootEl.children[multiDragElement.sortableIndex + (clonesInserted ? Number(i) : 0)];

			if (target) {
				rootEl.insertBefore(multiDragElement, target);
			} else {
				rootEl.appendChild(multiDragElement);
			}
		});
	}
	/**
	 * Insert multi-drag clones
	 * @param  {[Boolean]} elementsInserted  Whether the multi-drag elements are inserted
	 * @param  {HTMLElement} rootEl
	 */


	function insertMultiDragClones(elementsInserted, rootEl) {
		multiDragClones.forEach(function (clone) {
			var target = rootEl.children[clone.sortableIndex + (elementsInserted ? Number(i) : 0)];

			if (target) {
				rootEl.insertBefore(clone, target);
			} else {
				rootEl.appendChild(clone);
			}
		});
	}

	function removeMultiDragElements() {
		multiDragElements.forEach(function (multiDragElement) {
			if (multiDragElement === dragEl$1) return;
			multiDragElement.parentNode && multiDragElement.parentNode.removeChild(multiDragElement);
		});
	}

	Sortable.mount(new AutoScrollPlugin());
	Sortable.mount(Remove, Revert);

	Sortable.mount(new SwapPlugin());
	Sortable.mount(new MultiDragPlugin());

	return Sortable;
}));

/******************************************************************************/

(function(window, $) {
	var fave = function(window, $) {
		// Private
		var FormDataWasChanged = false;

		function IsDebugMode() {
			return window.fave_debug && window.fave_debug === true;
		};

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
			} else {
				ShowSystemMsgModal(title, message, error);
			}
		};

		function ShowSystemMsgModal(title, message, error) {
			$('#sys-modal-system-message-placeholder').html('');
			var html = '<div class="modal fade" id="sys-modal-system-message" tabindex="-1" role="dialog" aria-labelledby="sysModalSystemMessageLabel" aria-hidden="true"> \
				<div class="modal-dialog modal-dialog-centered" role="document"> \
					<div class="modal-content"> \
							<div class="modal-header"> \
								<h5 class="modal-title" id="sysModalSystemMessageLabel">' + title + '</h5> \
								<button type="button" class="close" data-dismiss="modal" aria-label="Close"> \
									<span aria-hidden="true">&times;</span> \
								</button> \
							</div> \
							<div class="modal-body text-left">' + message + '</div> \
							<div class="modal-footer"> \
								<button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button> \
							</div> \
					</div> \
				</div> \
			</div>';
			$('#sys-modal-system-message-placeholder').html(html);
			$('#sys-modal-system-message').modal({
				backdrop: 'static',
				keyboard: true,
				show: false,
			});
			$('#sys-modal-system-message').on('hidden.bs.modal', function(e) {
				$('#sys-modal-system-message-placeholder').html('');
			});
			$('#sys-modal-system-message').modal('show');
		};

		function AjaxEval(data) {
			try {
				eval(data);
			} catch(e) {
				if(e instanceof SyntaxError) {
					console.log(data);
					console.log('Error: JavaScript code eval error', e.message);
				}
			}
		};

		function AjaxDone(data) {
			AjaxEval(data);
		};

		function AjaxFail(data, status, error) {
			if(status.toLowerCase() === "error" && error.toLowerCase() === "not found") {
				AjaxEval(data);
			} else {
				console.log('Error: data sending error, page will be reloaded', data, status, error);
				setTimeout(function() {
					window.location.reload(false);
				}, 1000);
			}
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
					if(IsDebugMode()) console.log('done', data);
					AjaxDone(data);
				}).fail(function(xhr, status, error) {
					if(IsDebugMode()) console.log('fail', xhr, status, error);
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
						if(!$(this).hasClass('ignore-lost-data')) {
							FormDataWasChanged = true;
						}
					}
				});
			}
		};

		function PreventDataLost() {
			if(!FormDataWasChanged) {
				FormDataWasChanged = true;
			}
		};

		function FormDataIsChanged() {
			return FormDataWasChanged;
		}

		function HtmlDecode(value) {
			var doc = new DOMParser().parseFromString(value, "text/html");
			return doc.documentElement.textContent;
		};

		function HtmlFixEditorHtml(value) {
			newValue = value;
			newValue = newValue.replace(/&nbsp;/gi, '');
			return newValue;
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
		};

		function MakeTextAreasWysiwyg() {
			$('textarea.wysiwyg').each(function() {
				var area = $(this)[0];
				var area_id = area.id;
				var area_name = area.name;
				var area_html = area.innerHTML;

				// Wrap editro by additional DIV and remove target textarea
				$(area).wrap('<div id="' + area_id + '_wysiwyg" class="form-control wysiwyg" style="height:auto;padding:0px"></div>').remove();

				var wysiwyg = document.getElementById(area_id + '_wysiwyg');
				wysiwyg.id = area_id;

				// Create and init editor
				var editor = window.pell.init({
					element: wysiwyg,
					onChange: function(html) {
						area.innerHTML = HtmlFixEditorHtml(html);
						$(area).val(HtmlFixEditorHtml(html));
						if(!FormDataWasChanged) {
							FormDataWasChanged = true;
						}
					},
					defaultParagraphSeparator: 'p',
					styleWithCSS: false,
					actions: [
						'paragraph',
						'heading1',
						'heading2',
						'bold',
						'italic',
						'underline',
						'strikethrough',
						'ulist',
						'olist',
						'link',
						{
							name: 'htmlcode',
							icon: 'HTML',
							title: 'HTML Source',
							result: function(edt, ctn, btn) {
								var jedt = $(edt);
								var jctn = $(ctn);
								var jbtn = $(btn);
								if(!jbtn.hasClass('pell-button-html-pressed')) {
									jbtn.addClass('pell-button-html-pressed');
									jedt.addClass('pell-html-mode');

									jedt.find('.pell-actionbar .pell-button').prop('disabled', true);
									jbtn.prop('disabled', false);

									setTimeout(function() {
										jedt.find('textarea.form-control').focus();
									}, 0);
								} else {
									jbtn.removeClass('pell-button-html-pressed');
									jedt.removeClass('pell-html-mode');
									jedt.find('.pell-actionbar .pell-button').prop('disabled', false);

									var srcValue = jedt.find('textarea.form-control').val();
									ctn.innerHTML = HtmlFixEditorHtml(srcValue);
									$(ctn).val(HtmlFixEditorHtml(srcValue));

									setTimeout(function() {
										jctn.focus();
									}, 0);
								}
							},
						},
					],
					classes: {
						actionbar: 'pell-actionbar',
						button: 'pell-button',
						content: 'pell-content',
						selected: 'pell-button-selected'
					}
				});

				editor.onfocusin = function() {
					$(wysiwyg).addClass('focused');
				};

				editor.onfocusout = function() {
					$(wysiwyg).find('.pell-actionbar button.pell-button-selected').removeClass('pell-button-selected');
					$(wysiwyg).removeClass('focused');
				};

				// Re-add textarea
				$(wysiwyg).append('<textarea class="form-control" id="' + area_id + '_wysiwyg' + '" name="' + area_name + '" style="display:none"></textarea>');
				area = document.getElementById(area_id + '_wysiwyg');

				// Prevent data lost if HTML was changed
				$(area).on('input', function() {
					if(!FormDataWasChanged) {
						FormDataWasChanged = true;
					}
				});

				// Copy HTML to textarea and editor
				area.innerHTML = HtmlFixEditorHtml(HtmlDecode(area_html));
				$(area).val(HtmlFixEditorHtml(HtmlDecode(area_html)));
				editor.content.innerHTML = HtmlFixEditorHtml(HtmlDecode(area_html));
			});
		};

		function MakeTextAreasTmplEditor() {
			var IgnoreDataLost = true;
			$('textarea.tmpl-editor').each(function() {
				var targetTextArea = $(this)[0];
				var targetFileExt = $(this).data('emode');
				var targetEditorMode = 'text/html';
				if(targetFileExt == 'js') {
					targetEditorMode = 'javascript';
				} else if(targetFileExt == 'css') {
					targetEditorMode = 'css';
				}
				CodeMirror.fromTextArea(targetTextArea, {
					lineNumbers: true,
					lineWrapping: true,
					viewportMargin: Infinity,
					mode: targetEditorMode,
				}).on('change', function(editor){
					targetTextArea.value = editor.getValue();
					if(!IgnoreDataLost) {
						if(!FormDataWasChanged) {
							FormDataWasChanged = true;
						}
					}
				});
			});
			IgnoreDataLost = false;
		};

		function MakeTextAreasNotReactOnTab() {
			$('textarea.use-tab-key').each(function() {
				$(this).keydown(function(e) {
					if(e.keyCode === 9) {
						var start = this.selectionStart;
						var end = this.selectionEnd;
						var $this = $(this);
						var value = $this.val();
						$this.val(value.substring(0, start) + "\t" + value.substring(end));
						this.selectionStart = this.selectionEnd = start + 1;
						e.preventDefault();
						if(!FormDataWasChanged) {
							FormDataWasChanged = true;
						}
					}
				});
			});
		};

		function Initialize() {
			// Check if jQuery was loaded
			if(typeof $ == 'function') {
				AllFormsToAjax();
				BindWindowBeforeUnload();
				MakeTextAreasAutoSized();
				MakeTextAreasWysiwyg();
				MakeTextAreasTmplEditor();
				MakeTextAreasNotReactOnTab();
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

			FormDataWasChanged: function() {
				PreventDataLost();
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

			ShopProductsAdd: function() {
				var selText = $('#lbl_attributes option:selected').text();
				var selValue = $('#lbl_attributes').val();
				if(selValue == '0') { return; }
				$('#lbl_attributes')[0].selectedIndex = 0;
				$('#lbl_attributes').selectpicker('refresh');
				if($('#prod_attr_' + selValue).length > 0) { return; }
				$('#list').append('<div class="form-group" id="prod_attr_' + selValue + '"><div><b>' + selText + '</b></div><div class="position-relative"><select class="form-control" name="value.' + selValue + '" autocomplete="off" required disabled><option value="0">Loading values...</option></select><button type="button" class="btn btn-danger btn-dynamic-remove" onclick="fave.ShopProductsRemove(this);" disabled>&times;</button></div></div>');
				PreventDataLost();
				$.ajax({
					type: 'POST',
					url: '/cp/',
					data: {
						action: 'shop-get-attribute-values',
						id: selValue
					}
				}).done(function(data) {
					try {
						eval(data);
					} catch(e) {
						if(e instanceof SyntaxError) {
							console.log(data);
							console.log('Error: JavaScript code eval error', e.message);
						}
					}
				}).fail(function(xhr, status, error) {
					$('#prod_attr_' + selValue).remove();
					try {
						eval(xhr.responseText);
					} catch(e) {
						if(e instanceof SyntaxError) {
							console.log(xhr.responseText);
							console.log('Error: JavaScript code eval error', e.message);
						}
					}
				});
			},

			ShopProductsRemove: function(button) {
				$(button).parent().parent().remove();
				PreventDataLost();
			},

			ShopAttributesAdd: function() {
				$('#list').append('<div class="form-group position-relative"><input class="form-control" type="text" name="value.0" value="" placeholder="" autocomplete="off" required><button type="button" class="btn btn-danger btn-dynamic-remove" onclick="fave.ShopAttributesRemove(this);">&times;</button></div>');
				PreventDataLost();
				setTimeout(function() {
					$('#list input').last().focus();
				}, 100);
			},

			ShopAttributesRemove: function(button) {
				$(button).parent().remove();
				PreventDataLost();
			},

			ShopProductsUploadImage: function(action_name, product_id, input_id) {
				var file_el = $('#' + input_id)[0];
				if(!file_el.files) return;
				if(file_el.files.length <= 0) return;

				$('#img-upload-block input').prop('disabled', true);
				$('#upload-msg').css('display', 'block');

				var fd = new FormData();
				fd.append('action', action_name);
				fd.append('id', product_id);
				fd.append('count', file_el.files.length);
				for(var i = 0; i < file_el.files.length; i++) {
					fd.append('file_' + i, file_el.files[i]);
				}

				$.ajax({
					url: '/cp/',
					method: 'POST',
					type: 'POST',
					data: fd,
					contentType: false,
					processData: false
				}).done(function(data) {
					try {
						eval(data);
					} catch(e) {
						if(e instanceof SyntaxError) {
							console.log(data);
							console.log('Error: JavaScript code eval error', e.message);
						}
					}
				}).fail(function(xhr, status, error) {
					try {
						eval(xhr.responseText);
					} catch(e) {
						if(e instanceof SyntaxError) {
							console.log(xhr.responseText);
							console.log('Error: JavaScript code eval error', e.message);
						}
					}
				}).always(function() {
					file_el.value = '';
					$('#img-upload-block input').prop('disabled', false);
					$('#upload-msg').css('display', 'none');
				});
			},

			ShopProductsDeleteImage: function(button, product_id, filename) {
				if($(button).hasClass('in-progress')) return;
				$(button).addClass('in-progress');
				$.ajax({
					type: "POST",
					url: '/cp/',
					data: {
						action: 'shop-upload-delete',
						id: product_id,
						file: filename,
					}
				}).done(function(data) {
					if(IsDebugMode()) console.log('done', data);
					AjaxDone(data);
				}).fail(function(xhr, status, error) {
					if(IsDebugMode()) console.log('fail', xhr, status, error);
					AjaxFail(xhr.responseText, status, error);
				});
			},

			ShopProductsDuplicateBase: function(button, product_id, attach) {
				if($(button).hasClass('in-progress')) return;
				if(FormDataIsChanged()) {
					fave.ShowMsgError('Warning!', 'Something was changed, save changes before duplicate product', true);
					return;
				}
				$(button).addClass('in-progress');
				$.ajax({
					type: "POST",
					url: '/cp/',
					data: {
						action: 'shop-duplicate',
						id: product_id,
						attach: attach,
					}
				}).done(function(data) {
					try {
						eval(data);
					} catch(e) {
						if(e instanceof SyntaxError) {
							console.log(data);
							console.log('Error: JavaScript code eval error', e.message);
						}
					}
				}).fail(function(xhr, status, error) {
					try {
						eval(xhr.responseText);
					} catch(e) {
						if(e instanceof SyntaxError) {
							console.log(xhr.responseText);
							console.log('Error: JavaScript code eval error', e.message);
						}
					}
				}).always(function() {
					$(button).removeClass('in-progress');
				});
			},

			ShopProductsDuplicate: function(button, product_id) {
				fave.ShopProductsDuplicateBase(button, product_id, 0);
			},

			ShopProductsDuplicateWithAttach: function(button, product_id) {
				fave.ShopProductsDuplicateBase(button, product_id, 1);
			},

			ShopProductsRetryImage: function(img, id) {
				var target = $('#' + id);
				var src = target.attr('src');
				target.attr('src', '/assets/cp/img-load.gif');
				setTimeout(function() {
					target.attr('src', src);
				}, 1000);
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
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}).fail(function(xhr, status, error) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
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
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}).fail(function(xhr, status, error) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
						AjaxFail(xhr.responseText, status, error);
					});
				}
			},

			ActionRestoreThemeFile: function(action_name, file_name, message) {
				if(confirm(message)) {
					$.ajax({
						type: "POST",
						url: '/cp/',
						data: {
							action: action_name,
							file: file_name,
						}
					}).done(function(data) {
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}).fail(function(xhr, status, error) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
						AjaxFail(xhr.responseText, status, error);
					});
				}
				return false;
			},

			ShopProductsImageReorder: function(action, orderData) {
				$.ajax({
					type: "POST",
					url: '/cp/',
					data: {
						action: action,
						data: JSON.stringify(orderData),
					},
				}).done(function(data) {
					if(IsDebugMode()) console.log('done', data);
					AjaxDone(data);
				}).fail(function(xhr, status, error) {
					if(IsDebugMode()) console.log('fail', xhr, status, error);
					AjaxFail(xhr.responseText, status, error);
				});
			},

			ShopAttachProduct: function(product_id) {
				var html = '<div class="modal fade" id="sys-modal-shop-product-attach" tabindex="-1" role="dialog" aria-labelledby="sysModalShopProductLabel" aria-hidden="true"> \
					<div class="modal-dialog modal-dialog-centered" role="document"> \
						<div class="modal-content"> \
							<div class="modal-header"> \
								<h5 class="modal-title" id="sysModalShopProductLabel">Attach product</h5> \
								<button type="button" class="close" data-dismiss="modal" aria-label="Close"> \
									<span aria-hidden="true">&times;</span> \
								</button> \
							</div> \
							<div class="modal-body text-left"> \
								<div class="form-group"> \
									<input type="text" class="form-control" name="product-name" value="" placeholder="Type product name here..." readonly autocomplete="off"> \
								</div> \
								<div class="form-group" style="margin-bottom:0px;"> \
									<div class="products-list"></div> \
								</div> \
							</div> \
							<div class="modal-footer"> \
								<button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button> \
							</div> \
						</div> \
					</div> \
				</div>';
				$('#sys-modal-shop-product-attach-placeholder').html(html);
				$("#sys-modal-shop-product-attach").modal({
					backdrop: 'static',
					keyboard: true,
					show: false,
				});
				$('#sys-modal-shop-product-attach').on('hidden.bs.modal', function(e) {
					$('#sys-modal-shop-product-attach-placeholder').html('');
				});
				$("#sys-modal-shop-product-attach").modal('show');
				setTimeout(function() {
					var SearchInput = $('#sys-modal-shop-product-attach input[name="product-name"]');
					SearchInput.keyup(function() {
						if(true || this.value != '') {
							$.ajax({
								type: "POST",
								url: '/cp/',
								data: {
									action: 'shop-attach-product-search',
									words: this.value,
									id: product_id,
								}
							}).done(function(data) {
								if($('#sys-modal-shop-product-attach').length > 0) {
									if(IsDebugMode()) console.log('done', data);
									AjaxDone(data);
								}
							}).fail(function(xhr, status, error) {
								if(false) {
									if($('#sys-modal-shop-product-attach').length > 0) {
										if(IsDebugMode()) console.log('fail', xhr, status, error);
										AjaxFail(xhr.responseText, status, error);
									}
								}
							});
						}
					});
					SearchInput.attr("readonly", false);
					SearchInput.keyup();
					SearchInput.focus();
				}, 500);
			},

			ShopAttachProductTo: function(parent_id, product_id) {
				$.ajax({
					type: "POST",
					url: '/cp/',
					data: {
						action: 'shop-attach-product-to',
						parent_id: parent_id,
						product_id: product_id,
					}
				}).done(function(data) {
					if(IsDebugMode()) console.log('done', data);
					AjaxDone(data);
				}).fail(function(xhr, status, error) {
					if(IsDebugMode()) console.log('fail', xhr, status, error);
					AjaxFail(xhr.responseText, status, error);
				});
			},

			ShopSetOrderStatus: function(object, id, status, message) {
				if(confirm(message)) {
					$.ajax({
						type: "POST",
						url: '/cp/',
						data: {
							action: 'shop-order-set-status',
							id: id,
							status: status,
						}
					}).done(function(data) {
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}).fail(function(xhr, status, error) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
						AjaxFail(xhr.responseText, status, error);
					});
				}
			},

			FilesManagerDialog: function() {
				var html = '<div class="modal fade" id="sys-modal-files-manager" tabindex="-1" role="dialog" aria-labelledby="sysModalFilesManagerLabel" aria-hidden="true"> \
					<div class="modal-dialog modal-dialog-centered" role="document"> \
						<div class="modal-content"> \
							<input type="hidden" name="path" value="/"> \
							<div class="modal-header"> \
								<h5 class="modal-title" id="sysModalFilesManagerLabel">Files manager</h5> \
								<button type="button" class="close" data-dismiss="modal" aria-label="Close"> \
									<span aria-hidden="true">&times;</span> \
								</button> \
							</div> \
							<div class="modal-body text-left"> \
								<div class="dialog-path alert alert-secondary"><span class="text-dotted">/</span></div> \
								<div class="dialog-data"></div> \
							</div> \
							<div class="modal-footer"> \
								<input class="form-control" type="file" id="fmfiles" name="fmfiles" onchange="fave.FilesManagerUploadFile();" style="font-size:12px;background-color:#28a745;border-color:#28a745;color:#fff;cursor:pointer;" multiple=""> \
								<button type="button" class="btn btn-primary folder" onclick="fave.FilesManagerNewFolderClick();" disabled>New folder</button> \
								<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button> \
							</div> \
						</div> \
					</div> \
				</div>';
				$('#sys-modal-files-manager-placeholder').html(html);
				$("#sys-modal-files-manager").modal({
					backdrop: 'static',
					keyboard: true,
					show: false,
				});
				$('#sys-modal-files-manager').on('hidden.bs.modal', function(e) {
					$('#sys-modal-files-manager-placeholder').html('');
				});
				$("#sys-modal-files-manager").modal('show');

				setTimeout(function() {
					fave.FilesManagerLoadData('/');
				}, 500);
			},

			FilesManagerSetPath: function(path) {
				$('#sys-modal-files-manager input[name=path]').val(path);
				$('#sys-modal-files-manager .dialog-path span').html(path);
			},

			FilesManagerGetPath: function() {
				return $('#sys-modal-files-manager input[name=path]').val();
			},

			FilesManagerRemoveFolder: function(filename, msg) {
				if(!confirm(msg)) {
					return;
				}
				$.ajax({
					type: "POST",
					url: '/cp/',
					data: {
						action: 'files-remove-folder',
						file: filename,
					}
				}).done(function(data) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}
				}).fail(function(xhr, status, error) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
						AjaxFail(xhr.responseText, status, error);
					}
				});
			},

			FilesManagerRemoveFile: function(filename, msg) {
				if(!confirm(msg)) {
					return;
				}
				$.ajax({
					type: "POST",
					url: '/cp/',
					data: {
						action: 'files-remove-file',
						file: filename,
					}
				}).done(function(data) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}
				}).fail(function(xhr, status, error) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
						AjaxFail(xhr.responseText, status, error);
					}
				});
			},

			FilesManagerLoadData: function(path) {
				fave.FilesManagerEnableDisableButtons(true);
				$.ajax({
					type: "POST",
					url: '/cp/',
					data: {
						action: 'files-list',
						path: path,
					}
				}).done(function(data) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}
				}).fail(function(xhr, status, error) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
						AjaxFail(xhr.responseText, status, error);
					}
				});
			},

			FilesManagerLoadDataUp: function(path) {
				newPath = path.replace(/\/$/i, '');
				newPath = newPath.replace(/[^\/]+$/i, '');
				fave.FilesManagerLoadData(newPath);
			},

			FilesManagerEnableDisableButtons: function(disabled) {
				$('#sys-modal-files-manager #fmfiles').prop('disabled', disabled);
				$('#sys-modal-files-manager button.folder').prop('disabled', disabled);
			},

			FilesManagerUploadFile: function() {
				var file_el = $('#fmfiles')[0];
				if(!file_el.files) return;
				if(file_el.files.length <= 0) return;

				fave.FilesManagerEnableDisableButtons(true);

				var fd = new FormData();
				fd.append('action', 'files-upload');
				fd.append('count', file_el.files.length);
				fd.append('path', fave.FilesManagerGetPath());
				for(var i = 0; i < file_el.files.length; i++) {
					fd.append('file_' + i, file_el.files[i]);
				}

				$.ajax({
					url: '/cp/',
					method: 'POST',
					type: 'POST',
					data: fd,
					contentType: false,
					processData: false
				}).done(function(data) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('done', data);
						AjaxDone(data);
					}
				}).fail(function(xhr, status, error) {
					if($('#sys-modal-files-manager').length > 0) {
						if(IsDebugMode()) console.log('fail', xhr, status, error);
						AjaxFail(xhr.responseText, status, error);
					}
				}).always(function() {
					file_el.value = '';
					fave.FilesManagerEnableDisableButtons(false);
				});
			},

			FilesManagerNewFolderClick: function() {
				var folderName = prompt('Please enter new folder name', '');
				if(folderName != null) {
					path = fave.FilesManagerGetPath();
					$.ajax({
						type: "POST",
						url: '/cp/',
						data: {
							action: 'files-mkdir',
							path: path,
							name: folderName,
						}
					}).done(function(data) {
						if($('#sys-modal-files-manager').length > 0) {
							if(IsDebugMode()) console.log('done', data);
							AjaxDone(data);
						}
					}).fail(function(xhr, status, error) {
						if($('#sys-modal-files-manager').length > 0) {
							if(IsDebugMode()) console.log('fail', xhr, status, error);
							AjaxFail(xhr.responseText, status, error);
						}
					});
				}
			},
		};
	}(window, $);

	// Make it public
	window.fave = fave;
}(window, jQuery));