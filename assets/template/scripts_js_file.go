package template

var VarScriptsJsFile = []byte(`(function(window, $) {
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

		function ShopBasketBlockObject(object) {
			if(object && !$(object).hasClass('click-blocked')) {
				$(object).addClass('click-blocked');
			}
		};

		function ShopBasketUnBlockObject(object) {
			if(object && $(object).hasClass('click-blocked')) {
				$(object).removeClass('click-blocked');
			}
		};

		function ShopBasketObjectIsNotBlocked(object) {
			if(!object) {
				return true;
			}
			return !$(object).hasClass('click-blocked');
		};

		function ShopBasketSetNavBtnProductsCount(value) {
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
				ShopBasketSetNavBtnProductsCount(count);
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
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);

					// ShopBasketSetNavBtnProductsCount(0);
					// console.log('ShopOpenBasket', object);
					// --------------------------------------------------
					var html = '<div class="modal fade" id="sys-modal-shop-basket" tabindex="-1" role="dialog" aria-labelledby="sysModalShopBasketLabel" aria-hidden="true"> \
						<div class="modal-dialog modal-dialog-centered" role="document"> \
							<div class="modal-content"> \
								<input type="hidden" name="action" value="index-user-update-profile"> \
								<div class="modal-header"> \
									<h5 class="modal-title" id="sysModalShopBasketLabel">Basket</h5> \
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
									<button type="button" class="btn btn-secondary" data-dismiss="modal">Continue shopping</button> \
								</div> \
							</div> \
						</div> \
					</div>';
					$('#sys-modal-shop-basket-placeholder').html(html);
					$("#sys-modal-shop-basket").modal({
						backdrop: 'static',
						keyboard: true,
						show: false,
					});
					$('#sys-modal-shop-basket').on('hidden.bs.modal', function(e) {
						$('#sys-modal-shop-basket-placeholder').html('');
					});
					$("#sys-modal-shop-basket").modal('show');
					// --------------------------------------------------

					ShopBasketUnBlockObject(object);
				}
				return false;
			},

			ShopBasketProductAdd: function(object, product_id) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);
					ShopBasketAjaxCommand('plus', product_id, function(data) {
						frontend.ShopBasketOpen();
					}, function(xhr, status, error) {
						// console.log('fail', xhr, status, error, product_id);
						// Page reload
					}, function() {
						ShopBasketAjaxUpdateCount();
						ShopBasketUnBlockObject(object);
					});
				}
				return false;
			},

			ShopBasketProductPlus: function(object, product_id) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);
					ShopBasketAjaxCommand('plus', product_id, function(data) {
						// console.log('success', data, product_id);
						// Update popup content
					}, function(xhr, status, error) {
						// console.log('fail', xhr, status, error, product_id);
						// Page reload
					}, function() {
						ShopBasketAjaxUpdateCount();
						ShopBasketUnBlockObject(object);
					});
				}
				return false;
			},

			ShopBasketProductMinus: function(object, product_id) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);
					ShopBasketAjaxCommand('minus', product_id, function(data) {
						// console.log('success', data, product_id);
						// Update popup content
					}, function(xhr, status, error) {
						// console.log('fail', xhr, status, error, product_id);
						// Page reload
					}, function() {
						ShopBasketAjaxUpdateCount();
						ShopBasketUnBlockObject(object);
					});
				}
				return false;
			},

			ShopBasketProductRemove: function(object, product_id) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);
					ShopBasketAjaxCommand('remove', product_id, function(data) {
						// console.log('success', data, product_id);
						// Update popup content
					}, function(xhr, status, error) {
						// console.log('fail', xhr, status, error, product_id);
						// Page reload
					}, function() {
						ShopBasketAjaxUpdateCount();
						ShopBasketUnBlockObject(object);
					});
				}
				return false;
			},
		};
	}(window, $);

	window.frontend = frontend;
}(window, jQuery));`)
