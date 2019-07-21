package template

var AllData = map[string][]byte{
	"blog-category.html": VarBlogCategoryHtmlFile,
	"footer.html":        VarFooterHtmlFile,
	"styles.css":         VarStylesCssFile,
	"header.html":        VarHeaderHtmlFile,
	"blog.html":          VarBlogHtmlFile,
	"index.html":         VarIndexHtmlFile,
	"robots.txt":         VarRobotsTxtFile,
	"page.html":          VarPageHtmlFile,
	"404.html":           Var404HtmlFile,
	"blog-post.html":     VarBlogPostHtmlFile,
	"scripts.js":         VarScriptsJsFile,
	"sidebar-left.html":  VarSidebarLeftHtmlFile,
	"sidebar-right.html": VarSidebarRightHtmlFile,
}
