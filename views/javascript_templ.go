// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"

func goToUrl(id string) templ.ComponentScript {
	return templ.ComponentScript{
		Name:       `__templ_goToUrl_e6a8`,
		Function:   `function __templ_goToUrl_e6a8(id){window.location.href = "wallet/" + id}`,
		Call:       templ.SafeScript(`__templ_goToUrl_e6a8`, id),
		CallInline: templ.SafeScriptInline(`__templ_goToUrl_e6a8`, id),
	}
}
