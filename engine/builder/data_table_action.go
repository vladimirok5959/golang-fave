package builder

type DataTableActionRow struct {
	Icon    string
	Href    string
	Hint    string
	Target  string
	Classes string
}

func DataTableAction(data *[]DataTableActionRow) string {
	result := ``
	for _, row := range *data {
		target := ``
		if row.Target != "" {
			target = ` target="` + row.Target + `"`
		}

		classes := row.Classes
		if classes != "" {
			classes = " " + classes
		}

		result += `<a class="ico` + classes + `" title="` + row.Hint + `" href="` +
			row.Href + `"` + target + `>` + row.Icon + `</a>`
	}
	return result
}
