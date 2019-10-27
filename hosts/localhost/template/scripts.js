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
			$('#basket-mobile-btn').html(value);
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

		function ShopBasketAjaxGetInfo(success, fail, always) {
			$.ajax({
				type: "GET",
				dataType: 'json',
				url: '/shop/basket/info/'
			}).done(function(data) {
				if(success && data) { success(data); }
			}).fail(function(xhr, status, error) {
				if(fail) { fail(xhr, status, error); }
			}).always(function() {
				if(always) { always(); }
			});
		};

		function ShopBasketAjaxGetCount(success, fail, always) {
			ShopBasketAjaxGetInfo(function(data) {
				if(success && data && data.total_count != undefined) {
					success(data.total_count);
				}
			}, function(xhr, status, error) {
				if(fail) { fail(xhr, status, error); }
			}, function() {
				if(always) { always(); }
			});
		};

		function ShopBasketAjaxUpdateCount() {
			ShopBasketAjaxGetCount(function(count) {
				ShopBasketSetNavBtnProductsCount(count);
			});
		};

		function ShopBasketAjaxProductsHtml(success, fail, always) {
			ShopBasketAjaxGetInfo(function(data) {
				if(data) {
					if(data.total_count != undefined && data.total_count > 0) {
						var table = '';
						table += '<table class="table data-table table-striped table-bordered">';
						table += '<thead><tr><th class="thc-1">&nbsp;</th><th class="thc-2">' + ShopBasketTableProduct + '</th><th class="thc-3">' + ShopBasketTablePrice + '</th><th class="thc-4">' + ShopBasketTableQuantity + '</th><th class="thc-5">' + ShopBasketTableSum + '</th><th class="thc-6">&nbsp;</th></tr></thead>';
						table += '<tbody>';
						for(var i in data.products) {
							table += '<tr>';
							table += '<td class="thc-1"><img src="' + data.products[i].image + '" width="50" height="50" /></td>';
							table += '<td class="thc-2"><a href="' + data.products[i].link + '">' + data.products[i].name + '</a></td>';
							table += '<td class="thc-3">' + data.products[i].price + ' ' + data.currency.code + '</td>';
							table += '<td class="thc-4"><button type="button" class="btn btn-minus" onclick="frontend.ShopBasketProductMinus(this,' + data.products[i].id + ');"><span>-</span></button><input class="form-control" type="text" value="' + data.products[i].quantity + '" readonly><button type="button" class="btn btn-plus" onclick="frontend.ShopBasketProductPlus(this,' + data.products[i].id + ');"><span>+</span></button></td>';
							table += '<td class="thc-5">' + data.products[i].sum + ' ' + data.currency.code + '</td>';
							table += '<td class="thc-6"><a href="" onclick="frontend.ShopBasketProductRemove(this,' + data.products[i].id + ');return false;">&times;</a></td>';
							table += '</tr>';
						}
						table += '</tbody>';
						table += '</table>';
						table += '<div class="total"><span class="caption">' + ShopBasketTotal + '</span><span class="value">' + data.total_sum + ' ' + data.currency.code + '</span></div>';
						if(success) { success(table, data.total_count); }
					} else {
						if(success) { success(ShopBasketEmpty, 0); }
					}
				} else {
					window.location.reload(true);
				}
			}, function(xhr, status, error) {
				if(fail) { fail(xhr, status, error); }
			}, function() {
				if(always) { always(); }
			});
		};

		function ShopBasketEnableDisableOrderBtn(total) {
			$('#sys-modal-shop-basket button.btn-order').prop('disabled', total <= 0);
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
			ShopBasketBtnCollapse: function() {
				if(!$('.navbar-toggler').hasClass('collapsed')) {
					$('.navbar-toggler').click();
				}
				return true;
			},

			ShopBasketOpen: function(object) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);
					var html = '<div class="modal fade" id="sys-modal-shop-basket" tabindex="-1" role="dialog" aria-labelledby="sysModalShopBasketLabel" aria-hidden="true"> \
						<div class="modal-dialog modal-dialog-centered" role="document"> \
							<div class="modal-content"> \
								<input type="hidden" name="action" value="index-user-update-profile"> \
								<div class="modal-header"> \
									<h5 class="modal-title" id="sysModalShopBasketLabel">' + ShopBasketName + '</h5> \
									<button type="button" class="close" data-dismiss="modal" aria-label="Close"> \
										<span aria-hidden="true">&times;</span> \
									</button> \
								</div> \
								<div class="modal-body text-left" style="position:relative;"> \
									<div class="blocker" style="position:absolute;left:0px;top:0px;width:100%;height:100%;background:#fff;opacity:0.5;display:none;"></div> \
									<div class="data"></div> \
								</div> \
								<div class="modal-footer"> \
									<button type="button" class="btn btn-close btn-secondary" data-dismiss="modal">' + ShopBasketBtnContinue + '</button> \
									<button type="button" class="btn btn-order btn-success" disabled>' + ShopBasketBtnOrder + '</button> \
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
						$('#navbar-top').css('margin-right', $('#body').css('padding-right'));
					});

					ShopBasketAjaxProductsHtml(function(html, total) {
						$('#sys-modal-shop-basket .modal-body .data').html(html);
						ShopBasketEnableDisableOrderBtn(total);
						$("#sys-modal-shop-basket").modal('show');
						$('#navbar-top').css('margin-right', $('#body').css('padding-right'));
					}, function(xhr, status, error) {
						window.location.reload(true);
					});

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
						window.location.reload(true);
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
					$('#sys-modal-shop-basket .modal-body .blocker').css('display', 'block');
					ShopBasketAjaxCommand('plus', product_id, function(data) {
						ShopBasketAjaxProductsHtml(function(html, total) {
							$('#sys-modal-shop-basket .modal-body .data').html(html);
							ShopBasketEnableDisableOrderBtn(total);
						}, function(xhr, status, error) {
							window.location.reload(true);
						});
					}, function(xhr, status, error) {
						window.location.reload(true);
					}, function() {
						ShopBasketAjaxUpdateCount();
						ShopBasketUnBlockObject(object);
						$('#sys-modal-shop-basket .modal-body .blocker').css('display', 'none');
					});
				}
				return false;
			},

			ShopBasketProductMinus: function(object, product_id) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);
					$('#sys-modal-shop-basket .modal-body .blocker').css('display', 'block');
					ShopBasketAjaxCommand('minus', product_id, function(data) {
						ShopBasketAjaxProductsHtml(function(html, total) {
							$('#sys-modal-shop-basket .modal-body .data').html(html);
							ShopBasketEnableDisableOrderBtn(total);
						}, function(xhr, status, error) {
							window.location.reload(true);
						});
					}, function(xhr, status, error) {
						window.location.reload(true);
					}, function() {
						ShopBasketAjaxUpdateCount();
						ShopBasketUnBlockObject(object);
						$('#sys-modal-shop-basket .modal-body .blocker').css('display', 'none');
					});
				}
				return false;
			},

			ShopBasketProductRemove: function(object, product_id) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					ShopBasketBlockObject(object);
					$('#sys-modal-shop-basket .modal-body .blocker').css('display', 'block');
					ShopBasketAjaxCommand('remove', product_id, function(data) {
						ShopBasketAjaxProductsHtml(function(html, total) {
							$('#sys-modal-shop-basket .modal-body .data').html(html);
							ShopBasketEnableDisableOrderBtn(total);
						}, function(xhr, status, error) {
							window.location.reload(true);
						});
					}, function(xhr, status, error) {
						window.location.reload(true);
					}, function() {
						ShopBasketAjaxUpdateCount();
						ShopBasketUnBlockObject(object);
						$('#sys-modal-shop-basket .modal-body .blocker').css('display', 'none');
					});
				}
				return false;
			},
		};
	}(window, $);

	window.frontend = frontend;
}(window, jQuery));