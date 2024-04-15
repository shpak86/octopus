package templates

import (
	"strings"
)

type KeyValueTemplate struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HttpResponseTemplate struct {
	Log *string `json:"log"`
}

type HttpRequestTemplate struct {
	Target   *string               `json:"target"`
	Method   *string               `json:"method"`
	Delay    *string               `json:"delay"`
	Log      *string               `json:"log"`
	Timeout  *string               `json:"timeout"`
	Response *HttpResponseTemplate `json:"response"`
	Cookies  []KeyValueTemplate    `json:"cookies"`
	Headers  []KeyValueTemplate    `json:"headers"`
}

// Inject variables to the template
func (template *HttpRequestTemplate) Inject(variables map[string]string) {
	for name, value := range variables {
		variable := "${" + name + "}"
		injectString(template.Target, variable, value)
		injectString(template.Method, variable, value)
		injectString(template.Delay, variable, value)
		injectString(template.Log, variable, value)
		injectString(template.Timeout, variable, value)
		if template.Response != nil && template.Response.Log != nil {
			injectString(template.Response.Log, variable, value)
		}
		injectKeyValue(&template.Cookies, variable, value)
		injectKeyValue(&template.Headers, variable, value)
	}
}

// Inject string variable to the template
func injectString(s *string, k string, v string) {
	if s != nil {
		*s = strings.ReplaceAll(*s, k, v)
	}
}

// Inject variables to the array
func injectKeyValue(t *[]KeyValueTemplate, k string, v string) {
	for idx := range *t {
		(*t)[idx].Value = strings.ReplaceAll((*t)[idx].Value, k, v)
	}
}
