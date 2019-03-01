package builder

type DataTableActionRow struct {
	Icon   string
	Href   string
	Hint   string
	Target string
}

func DataTableAction(data *[]DataTableActionRow) string {
	result := ``
	for _, row := range *data {
		target := ``
		if row.Target != "" {
			target = ` target="` + row.Target + `"`
		}
		result += `<a class="ico" title="` + row.Hint + `" href="` +
			row.Href + `"` + target + `>` + row.Icon + `</a>`
	}
	return result
}
