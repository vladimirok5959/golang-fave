package builder

import (
	"html"

	"golang-fave/engine/wrapper"
)

const (
	DFKHidden = iota
	DFKText
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
	CallBack    func(field *DataFormField) string
}

func DataForm(wrap *wrapper.Wrapper, data []DataFormField) string {
	result := `<form class="data-form prev-data-lost" action="/cp/" method="post" autocomplete="off">`
	result += `<div class="hidden">`
	for _, field := range data {
		if field.Kind == DFKHidden {
			if field.CallBack != nil {
				result += field.CallBack(&field)
			} else {
				result += `<input type="hidden" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `">`
			}
		}
	}
	result += `</div>`
	for _, field := range data {
		if field.Kind != DFKHidden && field.Kind != DFKSubmit && field.Kind != DFKMessage {
			if field.CallBack != nil {
				result += field.CallBack(&field)
			} else {
				required := ``
				if field.Required {
					required = ` required`
				}
				result += `<div class="form-group">`
				result += `<div class="row">`
				result += `<div class="col-3">`

				if field.Kind != DFKCheckBox {
					result += `<label for="lbl_` + field.Name + `">` + field.Caption + `</label>`
				} else {
					result += `<label>` + field.Caption + `</label>`
				}

				result += `</div>`
				result += `<div class="col-9">`
				result += `<div>`
				if field.Kind == DFKText {
					result += `<input class="form-control" type="text" id="lbl_` + field.Name + `" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `" placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>`
				} else if field.Kind == DFKEmail {
					result += `<input class="form-control" type="email" id="lbl_` + field.Name + `" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `" placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>`
				} else if field.Kind == DFKPassword {
					result += `<input class="form-control" type="password" id="lbl_` + field.Name + `" name="` + field.Name + `" value="` + html.EscapeString(field.Value) + `" placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>`
				} else if field.Kind == DFKTextArea {
					result += `<textarea class="form-control" id="lbl_` + field.Name + `" name="` + field.Name + `" placeholder="` + field.Placeholder + `" autocomplete="off"` + required + `>` + html.EscapeString(field.Value) + `</textarea>`
				} else if field.Kind == DFKCheckBox {
					checked := ""
					if field.Value != "0" {
						checked = " checked"
					}
					result += `<div class="checkbox-ios"><input class="form-control" type="checkbox" id="lbl_` + field.Name + `" name="` + field.Name + `" value="1"` + `" autocomplete="off"` + required + checked + `><label for="lbl_` + field.Name + `"></label></div>`
				}
				result += `</div>`
				if field.Hint != "" {
					result += `<div><small>` + field.Hint + `</small></div>`
				}
				result += `</div>`
				result += `</div>`
				result += `</div>`
			}
		}
	}

	// TODO: optimize this...

	for _, field := range data {
		if field.Kind == DFKMessage {
			if field.CallBack != nil {
				result += field.CallBack(&field)
			} else {
				result += `<div class="row">`
				result += `<div class="col-3">`
				result += `</div>`
				result += `<div class="col-9">`
				result += `<div class="sys-messages"></div>`
				result += `</div>`
				result += `</div>`
			}
		}
	}

	for _, field := range data {
		if field.Kind == DFKSubmit {
			if field.CallBack != nil {
				result += field.CallBack(&field)
			} else {
				result += `<div class="row hidden">`
				result += `<div class="col-3">`
				result += `&nbsp;`
				result += `</div>`
				result += `<div class="col-9">`
				result += `<button type="submit" class="btn btn-primary" data-target="` + field.Target + `">` + html.EscapeString(field.Value) + `</button>`
				result += `</div>`
				result += `</div>`
			}
		}
	}
	result += `</form>`
	return result
}
