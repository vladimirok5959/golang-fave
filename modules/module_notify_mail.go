package modules

import (
	"html"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterModule_NotifyMail() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "notify-mail",
		Name:   "Mail notifier",
		Order:  803,
		System: true,
		Icon:   assets.SysSvgIconEmail,
		Sub: &[]MISub{
			{Mount: "default", Name: "All", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "success", Name: "Success", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "in-progress", Name: "In progress", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "error", Name: "Error", Show: true, Icon: assets.SysSvgIconList},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" || wrap.CurrSubModule == "success" || wrap.CurrSubModule == "in-progress" || wrap.CurrSubModule == "error" {
			ModuleName := "All"
			ModulePagination := "/cp/" + wrap.CurrModule + "/"
			ModuleSqlWhere := ""

			if wrap.CurrSubModule == "success" {
				ModuleName = "Success"
				ModulePagination = "/cp/" + wrap.CurrModule + "/success/"
				ModuleSqlWhere = " WHERE notify_mail.status = 1"
			} else if wrap.CurrSubModule == "in-progress" {
				ModuleName = "In progress"
				ModulePagination = "/cp/" + wrap.CurrModule + "/in-progress/"
				ModuleSqlWhere = " WHERE notify_mail.status = 2"
			} else if wrap.CurrSubModule == "error" {
				ModuleName = "Error"
				ModulePagination = "/cp/" + wrap.CurrModule + "/error/"
				ModuleSqlWhere = " WHERE notify_mail.status = 0"
			}

			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: ModuleName},
			})
			content += builder.DataTable(
				wrap,
				"notify_mail",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField:     "email",
						NameInTable: "Email / Subject",
						CallBack: func(values *[]string) string {
							email := `<a href="/cp/` + wrap.CurrModule + `/details/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
							subject := html.EscapeString((*values)[2])
							if subject != "" {
								subject = `<div><small>` + subject + `</small></div>`
							}
							return `<div>` + email + `</div>` + subject
						},
					},
					{
						DBField: "subject",
					},
					{
						DBField:     "datetime",
						DBExp:       "UNIX_TIMESTAMP(`datetime`)",
						NameInTable: "Date / Time",
						Classes:     "d-none d-md-table-cell",
						CallBack: func(values *[]string) string {
							t := int64(utils.StrToInt((*values)[3]))
							return `<div>` + utils.UnixTimestampToFormat(t, "02.01.2006") + `</div>` +
								`<div><small>` + utils.UnixTimestampToFormat(t, "15:04:05") + `</small></div>`
						},
					},
					{
						DBField:     "status",
						NameInTable: "Status",
						Classes:     "d-none d-sm-table-cell",
						CallBack: func(values *[]string) string {
							return builder.CheckBox(utils.StrToInt((*values)[4]))
						},
					},
				},
				nil,
				ModulePagination,
				func() (int, error) {
					var count int
					return count, wrap.DB.QueryRow(
						"SELECT COUNT(*) FROM `notify_mail`" + ModuleSqlWhere + ";",
					).Scan(&count)
				},
				func(limit_offset int, pear_page int) (*sqlw.Rows, error) {
					return wrap.DB.Query(
						`SELECT
							notify_mail.id,
							notify_mail.email,
							notify_mail.subject,
							UNIX_TIMESTAMP(`+"`notify_mail`.`datetime`"+`) AS datetime,
							notify_mail.status
						FROM
							notify_mail
						`+ModuleSqlWhere+`
						ORDER BY
							notify_mail.id DESC
						LIMIT ?, ?;`,
						limit_offset,
						pear_page,
					)
				},
				true,
			)
		} else if wrap.CurrSubModule == "details" {
			//
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
