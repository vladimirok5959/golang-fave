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

		function Initialize() {
			// Check if jQuery was loaded
			if(typeof $ == 'function') {
				ShopProductsInitLightGallery();
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
			ShopOpenBasket: function(object) {
				if(!$(object).hasClass('click-blocked')) {
					$(object).addClass('click-blocked');

					ShopSetBasketNavBtnProductsCount(0);
					console.log('ShopOpenBasket', object);

					$(object).removeClass('click-blocked');
				}
				return false;
			},

			// TODO: add product to basket or count++ if already in basket
			// Update products counter in header nav bar button
			// Automatically open basket popup
			ShopAddProductToBasket: function(object, product_id) {
				if(!$(object).hasClass('click-blocked')) {
					$(object).addClass('click-blocked');

					ShopSetBasketNavBtnProductsCount(product_id);
					console.log('ShopAddProductToBasket', object, product_id);

					$(object).removeClass('click-blocked');
				}
				return false;
			},
		};
	}(window, $);

	window.frontend = frontend;
}(window, jQuery));