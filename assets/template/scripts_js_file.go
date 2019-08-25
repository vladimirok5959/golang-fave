package template

var VarScriptsJsFile = []byte(`$(document).ready(function() {
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
			$('#product_image img').attr('src', $(this).data('hover'));
			$('#product_image img').data('index', $(this).data('index'));
		});
    });
});`)
