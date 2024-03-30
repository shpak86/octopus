package requests

import (
	"io"
	"log"
	"net/http"
	"octopus/internal/templates"
	"strconv"
	"strings"
	"time"
)

type RequestSender struct {
	templatesChannel chan templates.Template
}

func NewRequestSender(templatesChannel chan templates.Template) *RequestSender {
	return &RequestSender{templatesChannel: templatesChannel}
}

func (rs *RequestSender) Serve() {
	go func() {
		for template := range rs.templatesChannel {
			rs.prepareRequest(&template)
			respData, err := rs.send(&template)
			rs.report(&template, &respData, err)
		}
	}()
}

func (rs *RequestSender) send(template *templates.Template) (respData responseData, err error) {
	httpClient := &http.Client{}
	var request *http.Request
	if template.Timeout != nil {
		if ms, err := strconv.Atoi(*template.Timeout); err != nil {
			httpClient.Timeout = time.Duration(ms) * time.Millisecond
		}
	}
	method := "get"
	if template.Method != nil {
		method = *template.Method
	}
	request, _ = http.NewRequest(method, *template.Target, nil)
	for _, header := range template.Headers {
		request.Header.Add(header.Name, header.Value)
	}
	for _, cookie := range template.Cookies {
		request.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}
	resp, err := httpClient.Do(request)
	if err == nil {
		respData.code = resp.StatusCode
		respData.body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
	return
}

func (rs *RequestSender) prepareRequest(template *templates.Template) {
	if template.Delay != nil {
		if ms, err := strconv.Atoi(*template.Delay); err == nil {
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}
	}
	if template.Log != nil {
		log.Println(*template.Log)
	}
}

func (rs *RequestSender) report(template *templates.Template, respData *responseData, err error) {
	if err != nil {
		log.Println(err)
	} else if template.Response != nil && template.Response.Log != nil {
		respCode := strconv.Itoa(respData.code)
		message := strings.ReplaceAll(*template.Response.Log, "${respCode}", respCode)
		message = strings.ReplaceAll(message, "${respBody}", string(string(respData.body)))
		log.Println(message)
	}
}
