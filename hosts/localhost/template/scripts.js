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
							table += '<td class="thc-1 d-none d-md-table-cell"><img src="' + data.products[i].image + '" width="50" height="50" /></td>';
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

		function ClearOrderFormErrorMessage() {
			$('#sys-modal-shop-basket .modal-body .order-form .sys-messages').html('');
			$('#sys-modal-shop-basket .modal-body .order-form .input-error-msg').css('display', 'none');
		};

		function ShowOrderFormErrorMessage(title, message) {
			$('#sys-modal-shop-basket .modal-body .order-form .sys-messages').html('<div class="alert alert-danger alert-dismissible fade show" role="alert"><strong>' + title + '</strong> ' + message + '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>');
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
									<div class="order-form mt-3" style="display:none;"> \
										<hr class="mb-4"> \
										<form class="data-form" action="/" method="post" autocomplete="off"> \
											<div class="hidden"><input type="hidden" name="action" value="shop-order"></div> \
											<div class="row"> \
												<div class="col-md-6"> \
													<div class="form-group"> \
														<div class="row"> \
															<div class="col-md-4"> \
																<label for="lbl_client_last_name">' + ShopOrderLabelLastName + '</label> \
															</div> \
															<div class="col-md-8"> \
																<div> \
																	<input class="form-control" type="text" id="lbl_client_last_name" name="client_last_name" minlength="1" maxlength="64" autocomplete="off" onkeydown="$(this).parent().parent().find(\'.input-error-msg\').css(\'display\', \'none\');"> \
																</div> \
																<div class="input-error-msg" style="display:none;"> \
																	<small>' + ShopOrderEmptyLastName + '</small> \
																</div> \
															</div> \
														</div> \
													</div> \
													<div class="form-group"> \
														<div class="row"> \
															<div class="col-md-4"> \
																<label for="lbl_client_first_name">' + ShopOrderLabelFirstName + '</label> \
															</div> \
															<div class="col-md-8"> \
																<div> \
																	<input class="form-control" type="text" id="lbl_client_first_name" name="client_first_name" minlength="1" maxlength="64" autocomplete="off" onkeydown="$(this).parent().parent().find(\'.input-error-msg\').css(\'display\', \'none\');"> \
																</div> \
																<div class="input-error-msg" style="display:none;"> \
																	<small>' + ShopOrderEmptyFirstName + '</small> \
																</div> \
															</div> \
														</div> \
													</div> \
													<div class="form-group"> \
														<div class="row"> \
															<div class="col-md-4"> \
																<label for="lbl_client_middle_name">' + ShopOrderLabelMiddleName + '</label> \
															</div> \
															<div class="col-md-8"> \
																<div> \
																	<input class="form-control" type="text" id="lbl_client_middle_name" name="client_middle_name" minlength="1" maxlength="64" autocomplete="off" onkeydown="$(this).parent().parent().find(\'.input-error-msg\').css(\'display\', \'none\');"> \
																</div> \
																<div class="input-error-msg" style="display:none;"> \
																	<small>' + ShopOrderEmptyMiddleName + '</small> \
																</div> \
															</div> \
														</div> \
													</div> \
												</div> \
												<div class="col-md-6"> \
													<div class="form-group"> \
														<div class="row"> \
															<div class="col-md-4"> \
																<label for="lbl_client_phone">' + ShopOrderLabelMobilePhone + '</label> \
															</div> \
															<div class="col-md-8"> \
																<div> \
																	<input class="form-control" type="text" id="lbl_client_phone" name="client_phone" minlength="1" maxlength="20" autocomplete="off" onkeydown="$(this).parent().parent().find(\'.input-error-msg\').css(\'display\', \'none\');"> \
																</div> \
																<div class="input-error-msg" style="display:none;"> \
																	<small>' + ShopOrderEmptyMobilePhone + '</small> \
																</div> \
															</div> \
														</div> \
													</div> \
													<div class="form-group"> \
														<div class="row"> \
															<div class="col-md-4"> \
																<label for="lbl_client_email">' + ShopOrderLabelEmailAddress + '</label> \
															</div> \
															<div class="col-md-8"> \
																<div> \
																	<input class="form-control" type="text" id="lbl_client_email" name="client_email" minlength="1" maxlength="64" autocomplete="off" onkeydown="$(this).parent().parent().find(\'.input-error-msg\').css(\'display\', \'none\');"> \
																</div> \
																<div class="input-error-msg" style="display:none;"> \
																	<small>' + ShopOrderEmptyEmailAddress + '</small> \
																</div> \
															</div> \
														</div> \
													</div> \
													<div class="form-group"> \
														<div class="row"> \
															<div class="col-md-4"> \
																<label for="lbl_client_delivery_comment">' + ShopOrderLabelDelivery + '</label> \
															</div> \
															<div class="col-md-8"> \
																<div> \
																	<input class="form-control" type="text" id="lbl_client_delivery_comment" name="client_delivery_comment" minlength="1" maxlength="255" autocomplete="off" onkeydown="$(this).parent().parent().find(\'.input-error-msg\').css(\'display\', \'none\');"> \
																</div> \
																<div class="input-error-msg" style="display:none;"> \
																	<small>' + ShopOrderEmptyDelivery + '</small> \
																</div> \
															</div> \
														</div> \
													</div> \
												</div> \
											</div> \
											<div class="form-group"> \
												<div class="row"> \
													<div class="col-md-12 d-table-cell d-md-none"> \
														<label for="lbl_client_order_comment">' + ShopOrderLabelComment + '</label> \
													</div> \
													<div class="col-md-12"> \
														<div> \
															<textarea class="form-control" id="lbl_client_order_comment" name="client_order_comment" autocomplete="off" onkeydown="$(this).parent().parent().find(\'.input-error-msg\').css(\'display\', \'none\');"></textarea> \
														</div> \
														<div class="input-error-msg" style="display:none;"> \
															<small>' + ShopOrderEmptyComment + '</small> \
														</div> \
													</div> \
												</div> \
											</div> \
											<button type="submit" style="display:none;"></button> \ \
										</form> \
										<div class="sys-messages"></div>\
									</div> \
								</div> \
								<div class="modal-footer"> \
									<button type="button" class="btn btn-close btn-secondary" data-dismiss="modal">' + ShopBasketBtnContinue + '</button> \
									<button type="button" class="btn btn-order btn-success" onclick="frontend.ShopBasketMakeOrder(this);return false;" disabled>' + ShopBasketBtnOrder + '</button> \
								</div> \
							</div> \
						</div> \
					</div>';
					$('#sys-modal-shop-basket-placeholder').html(html);
					$('#sys-modal-shop-basket .order-form form').submit(function(e) {
						$('#sys-modal-shop-basket .modal-footer button.btn-success').click();
						e.preventDefault();
					});
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
			ShopBasketMakeOrder: function(object) {
				if(ShopBasketObjectIsNotBlocked(object)) {
					var OrderFormBlock = $('#sys-modal-shop-basket .modal-body .order-form');
					if(OrderFormBlock.css('display') == 'none') {
						OrderFormBlock.css('display', 'block');
						setTimeout(function() { OrderFormBlock.find('input.form-control').first().focus(); }, 200);
						return;
					}
					ClearOrderFormErrorMessage();
					var OrderForm = $('#sys-modal-shop-basket .modal-body .order-form form');
					// Validate
					var ValidateError = false;
					if(ShopOrderRequiredLastName && $.trim(OrderForm.find('input[name=client_last_name]').val()) == '') {
						OrderForm.find('input[name=client_last_name]').parent().parent().find('.input-error-msg').css('display', 'inline');
						if(!ValidateError) {
							setTimeout(function() { OrderForm.find('input[name=client_last_name]').first().focus(); }, 200);
						}
						ValidateError = true;
					}
					if(ShopOrderRequiredFirstName && $.trim(OrderForm.find('input[name=client_first_name]').val()) == '') {
						OrderForm.find('input[name=client_first_name]').parent().parent().find('.input-error-msg').css('display', 'inline');
						if(!ValidateError) {
							setTimeout(function() { OrderForm.find('input[name=client_first_name]').first().focus(); }, 200);
						}
						ValidateError = true;
					}
					if(ShopOrderRequiredMiddleName && $.trim(OrderForm.find('input[name=client_middle_name]').val()) == '') {
						OrderForm.find('input[name=client_middle_name]').parent().parent().find('.input-error-msg').css('display', 'inline');
						if(!ValidateError) {
							setTimeout(function() { OrderForm.find('input[name=client_middle_name]').first().focus(); }, 200);
						}
						ValidateError = true;
					}
					if(ShopOrderRequiredMobilePhone && $.trim(OrderForm.find('input[name=client_phone]').val()) == '') {
						OrderForm.find('input[name=client_phone]').parent().parent().find('.input-error-msg').css('display', 'inline');
						if(!ValidateError) {
							setTimeout(function() { OrderForm.find('input[name=client_phone]').first().focus(); }, 200);
						}
						ValidateError = true;
					}
					if(ShopOrderRequiredEmailAddress && $.trim(OrderForm.find('input[name=client_email]').val()) == '') {
						OrderForm.find('input[name=client_email]').parent().parent().find('.input-error-msg').css('display', 'inline');
						if(!ValidateError) {
							setTimeout(function() { OrderForm.find('input[name=client_email]').first().focus(); }, 200);
						}
						ValidateError = true;
					}
					if(ShopOrderRequiredDelivery && $.trim(OrderForm.find('input[name=client_delivery_comment]').val()) == '') {
						OrderForm.find('input[name=client_delivery_comment]').parent().parent().find('.input-error-msg').css('display', 'inline');
						if(!ValidateError) {
							setTimeout(function() { OrderForm.find('input[name=client_delivery_comment]').first().focus(); }, 200);
						}
						ValidateError = true;
					}
					if(ShopOrderRequiredComment && $.trim(OrderForm.find('textarea[name=client_order_comment]').val()) == '') {
						OrderForm.find('textarea[name=client_order_comment]').parent().parent().find('.input-error-msg').css('display', 'inline');
						if(!ValidateError) {
							setTimeout(function() { OrderForm.find('textarea[name=client_order_comment]').first().focus(); }, 200);
						}
						ValidateError = true;
					}
					if(ValidateError) {
						return;
					}
					// Send form
					ShopBasketBlockObject(object);
					$.ajax({
						type: "POST",
						url: OrderForm.attr('action'),
						data: OrderForm.serialize()
					}).done(function(data) {
						try {
							jdata = JSON.parse(data);
							if(jdata.error === true) {
								ShowOrderFormErrorMessage(ShopOrderError, window[jdata.variable]);
								if(jdata.field) {
									setTimeout(function() { OrderForm.find('[name=' + jdata.field + ']').first().focus(); }, 200);
									OrderForm.find('[name=' + jdata.field + ']').parent().parent().find('.input-error-msg').css('display', 'inline');
								}
							} else {
								ShopBasketSetNavBtnProductsCount(0);
								OrderFormBlock.css('display', 'none');
								$('#sys-modal-shop-basket .modal-body .data').html(window[jdata.variable]);
								$('#sys-modal-shop-basket button.btn-order').prop('disabled', true);
							}
						} catch(e) {
							ShowOrderFormErrorMessage(ShopOrderError, e.message);
						}
					}).fail(function(xhr, status, error) {
						ShowOrderFormErrorMessage(ShopOrderError, error);
					}).always(function() {
						ShopBasketUnBlockObject(object);
					});
				}
			},
		};
	}(window, $);

	window.frontend = frontend;
}(window, jQuery));