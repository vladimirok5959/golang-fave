package builder

import (
	"html"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

const (
	DFKHidden = iota
	DFKText
	DFKNumber
	DFKEmail
	DFKPassword
	DFKTextArea
	DFKCheckBox
	DFKSubmit
	DFKMessage
)

type DataFormField struct {
	Caption     string
	Kind        int
	Name        string
	Value       string
	Placeholder string
	Hint        string
	Target      string
	Required    bool
	Classes     string
	Min         string
	Max         string
	CallBack    func(field *DataFormField) string
}

func DataForm(wrap *wrapper.Wrapper, data []DataFormField) string {
	var html_hidden string
	var html_element string
	var html_message string
	var html_button string

	for i, field := range data {
		if field.Kind == DFKHidden {
			if field.CallBack != nil {
				html_hidden += field.CallBack(&field)
			} else {
				html_hidden += `<input type="hidden" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `">`
			}
		} else if field.Kind != DFKHidden && field.Kind != DFKSubmit && field.Kind != DFKMessage {
			if field.CallBack != nil {
				html_element += field.CallBack(&field)
			} else {
				required := ``
				if field.Required {
					required = ` required`
				}

				classes := field.Classes
				if classes != "" {
					classes = " " + classes
				}

				html_element += `<div class="form-group n` + utils.IntToStr(i) + `">`
				html_element += `<div class="row">`
				html_element += `<div class="col-md-3">`

				if field.Kind != DFKCheckBox {
					html_element += `<label for="lbl_` + field.Name + `">` + field.Caption + `</label>`
				} else {
					html_element += `<label>` + field.Caption + `</label>`
				}

				html_element += `</div>`
				html_element += `<div class="col-md-9">`
				html_element += `<div>`
				if field.Kind == DFKText {
					html_element += `<input class="form-control` + classes + `" type="text" id="lbl_` + field.Name + `" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `" `
					if field.Min != "" {
						html_element += `minlength="` + field.Min + `" `
					}
					if field.Max != "" {
						html_element += `maxlength="` + field.Max + `" `
					}
					html_element += `placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>`
				} else if field.Kind == DFKNumber {
					html_element += `<input class="form-control` + classes + `" type="number" id="lbl_` + field.Name + `" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `" `
					if field.Min != "" {
						html_element += `min="` + field.Min + `" `
					}
					if field.Max != "" {
						html_element += `max="` + field.Max + `" `
					}
					html_element += `placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>`
				} else if field.Kind == DFKEmail {
					html_element += `<input class="form-control` + classes + `" type="email" id="lbl_` + field.Name + `" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `" placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>`
				} else if field.Kind == DFKPassword {
					html_element += `<input class="form-control` + classes + `" type="password" id="lbl_` + field.Name + `" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `" placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>`
				} else if field.Kind == DFKTextArea {
					html_element += `<textarea class="form-control` + classes + `" id="lbl_` + field.Name + `" name="` + field.Name + `" `
					if field.Min != "" {
						html_element += `minlength="` + field.Min + `" `
					}
					if field.Max != "" {
						html_element += `maxlength="` + field.Max + `" `
					}
					html_element += `placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>` + html.EscapeString(field.Value) + `</textarea>`
				} else if field.Kind == DFKCheckBox {
					checked := ""
					if field.Value != "0" {
						checked = " checked"
					}
					html_element += `<div class="checkbox-ios"><input class="form-control` + classes + `" type="checkbox" id="lbl_` + field.Name + `" name="` + field.Name + `" value="1"` + `" autocomplete="off"` + required + checked + `><label for="lbl_` + field.Name + `"></label></div>`
				}
				html_element += `</div>`
				if field.Hint != "" {
					html_element += `<div><small>` + field.Hint + `</small></div>`
				}
				html_element += `</div>`
				html_element += `</div>`
				html_element += `</div>`
			}
		} else if field.Kind == DFKMessage {
			if field.CallBack != nil {
				html_message += field.CallBack(&field)
			} else {
				html_message += `<div class="row">`
				html_message += `<div class="col-md-3">`
				html_message += `</div>`
				html_message += `<div class="col-md-9">`
				html_message += `<div class="sys-messages"></div>`
				html_message += `</div>`
				html_message += `</div>`
			}
		} else if field.Kind == DFKSubmit {
			if field.CallBack != nil {
				html_button += field.CallBack(&field)
			} else {
				html_button += `<div class="row d-lg-none">`
				html_button += `<div class="col-md-3 d-none d-md-block">`
				html_button += `&nbsp;`
				html_button += `</div>`
				html_button += `<div class="col-md-9">`
				html_button += `<div class="pt-3"><button type="submit" class="btn btn-primary" data-target="` + field.Target + `">` + html.EscapeString(field.Value) + `</button></div>`
				html_button += `</div>`
				html_button += `</div>`
			}
		}
	}

	if html_hidden != "" {
		html_hidden = `<div class="hidden">` + html_hidden + `</div>`
	}

	return `<form class="data-form ` + wrap.CurrModule + `-` + wrap.CurrSubModule + ` prev-data-lost" action="/cp/" method="post" autocomplete="off">` +
		html_hidden + html_element + html_message + html_button + `</form>`
}
