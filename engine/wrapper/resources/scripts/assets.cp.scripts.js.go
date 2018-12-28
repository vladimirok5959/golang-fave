package scripts

var File_assets_cp_scripts_js = []byte(`
$(document).ready(function() {
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
				try {
					eval(data);
				} catch(e) {
					if(e instanceof SyntaxError) {
						console.log('JavaScript eval error:', e.message);
						console.log(data);
					}
				}
			}).fail(function() {
				console.log('Form send fail, page will be reloaded');
				window.location.reload(false);
			}).always(function() {
				$(form).removeClass('loading');
				$(button).removeClass('progress-bar-striped').removeClass('progress-bar-animated');
			});

			e.preventDefault();
		});
	});
});
`)
