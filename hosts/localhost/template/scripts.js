(function(window, $) {
	var frontend = function(window, $) {
		var ShopProductsLightGalleryLastImage = '';

		// Private
		function ShopProductsInitLightGallery() {
			$('#product_thumbnails').lightGallery({
				pager: false,
				showThumbByDefault: false,
				toogleThumb: false,
				controls: false,
				download: false
			});
			$('#product_image img').css('cursor', 'pointer').click(function() {
				$($('#product_thumbnails a.thumbnail').get($(this).data('index'))).click();
			});
			$('#product_thumbnails a.thumbnail').each(function() {
				$(this).mouseover(function() {
					if($(this).data('hover') != ShopProductsLightGalleryLastImage) {
						ShopProductsLightGalleryLastImage = $(this).data('hover');
						$('#product_image img').attr('src', $(this).data('hover'));
						$('#product_image img').data('index', $(this).data('index'));
					}
				});
			});
		};

		function ShopSetBasketNavBtnProductsCount(value) {
			$('#basket-nav-btn .badge').html(value);
		};

		function ShopBasketAjaxCommand(cmd, product_id, success, fail, always) {
			$.ajax({
				type: "GET",
				dataType: 'json',
				url: '/shop/basket/' + cmd + '/' + product_id + '/'
			}).done(function(data) {
				if(success) { success(data); }
			}).fail(function(xhr, status, error) {
				if(fail) { fail(xhr, status, error); }
			}).always(function() {
				if(always) { always(); }
			});
		};

		function ShopBasketAjaxGetCount(success, fail, always) {
			$.ajax({
				type: "GET",
				dataType: 'json',
				url: '/shop/basket/info/'
			}).done(function(data) {
				if(success && data && data.total_count) { success(data.total_count); }
			}).fail(function(xhr, status, error) {
				if(fail) { fail(xhr, status, error); }
			}).always(function() {
				if(always) { always(); }
			});
		};

		function ShopBasketAjaxUpdateCount() {
			ShopBasketAjaxGetCount(function(count) {
				ShopSetBasketNavBtnProductsCount(count);
			});
		};

		function Initialize() {
			// Check if jQuery was loaded
			if(typeof $ == 'function') {
				ShopProductsInitLightGallery();
				ShopBasketAjaxUpdateCount();
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
			ShopBasketOpen: function(object) {
				if(!$(object).hasClass('click-blocked')) {
					$(object).addClass('click-blocked');

					// ShopSetBasketNavBtnProductsCount(0);
					// console.log('ShopOpenBasket', object);

					$(object).removeClass('click-blocked');
				}
				return false;
			},

			ShopBasketProductPlus: function(object, product_id) {
				if(!$(object).hasClass('click-blocked')) {
					$(object).addClass('click-blocked');
					ShopBasketAjaxCommand('plus', product_id, function(data) {
						// console.log('success', data, product_id);
					}, function(xhr, status, error) {
						// console.log('fail', xhr, status, error, product_id);
					}, function() {
						ShopBasketAjaxUpdateCount();
						$(object).removeClass('click-blocked');
					});
				}
				return false;
			},

			ShopBasketProductMinus: function(object, product_id) {
				if(!$(object).hasClass('click-blocked')) {
					$(object).addClass('click-blocked');
					ShopBasketAjaxCommand('minus', product_id, function(data) {
						// console.log('success', data, product_id);
					}, function(xhr, status, error) {
						// console.log('fail', xhr, status, error, product_id);
					}, function() {
						ShopBasketAjaxUpdateCount();
						$(object).removeClass('click-blocked');
					});
				}
				return false;
			},

			ShopBasketProductRemove: function(object, product_id) {
				if(!$(object).hasClass('click-blocked')) {
					$(object).addClass('click-blocked');
					ShopBasketAjaxCommand('remove', product_id, function(data) {
						// console.log('success', data, product_id);
					}, function(xhr, status, error) {
						// console.log('fail', xhr, status, error, product_id);
					}, function() {
						ShopBasketAjaxUpdateCount();
						$(object).removeClass('click-blocked');
					});
				}
				return false;
			},
		};
	}(window, $);

	window.frontend = frontend;
}(window, jQuery));