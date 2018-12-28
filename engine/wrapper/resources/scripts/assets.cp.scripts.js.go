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

			$.ajax({
				type: "POST",
				url: form.attr('action'),
				data: form.serialize()
			}).done(function(data) {
				console.log('done');
				console.log(data);
			}).fail(function() {
				console.log('fail');
			}).always(function() {
				$(form).removeClass('loading');
				console.log('always');
			});

			console.log('1');
			e.preventDefault();
		});
	});
});
`)
