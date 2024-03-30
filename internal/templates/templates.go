package templates

import (
	"strings"
)

type KeyValueTemplate struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResponseTemplate struct {
	Log *string `json:"log"`
}

type Template struct {
	Target   *string            `json:"target"`
	Url      *string            `json:"url"`
	Method   *string            `json:"method"`
	Version  *string            `json:"version"`
	Delay    *string            `json:"delay"`
	Log      *string            `json:"log"`
	Response *ResponseTemplate  `json:"response"`
	Cookies  []KeyValueTemplate `json:"cookies"`
	Headers  []KeyValueTemplate `json:"headers"`
}

func (template *Template) Inject(variables map[string]string) {
	for name, value := range variables {
		variable := "${" + name + "}"
		injectVariable(template.Target, variable, value)
		injectVariable(template.Url, variable, value)
		injectVariable(template.Method, variable, value)
		injectVariable(template.Version, variable, value)
		injectVariable(template.Delay, variable, value)
		injectVariable(template.Log, variable, value)
		if template.Response != nil && template.Response.Log != nil {
			injectVariable(template.Response.Log, variable, value)
		}
		injectKeyValueVariable(&template.Cookies, variable, value)
		injectKeyValueVariable(&template.Headers, variable, value)
	}
}

func injectVariable(s *string, k string, v string) {
	if s != nil {
		*s = strings.ReplaceAll(*s, k, v)
	}
}

func injectKeyValueVariable(t *[]KeyValueTemplate, k string, v string) {
	for idx := range *t {
		(*t)[idx].Value = strings.ReplaceAll((*t)[idx].Value, k, v)
	}
}
