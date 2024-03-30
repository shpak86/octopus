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

type responseData struct {
	code int
	body []byte
}

type HttpRequestsManager struct {
	templates    *templates.TemplatesRepository
	defaultDelay *time.Duration
}

type HttpRequestsManagerBuilder struct {
	manager *HttpRequestsManager
}

func NewHttpRequestsManagerBuilder() *HttpRequestsManagerBuilder {
	return &HttpRequestsManagerBuilder{
		manager: &HttpRequestsManager{},
	}
}

func (b *HttpRequestsManagerBuilder) Templates(repository *templates.TemplatesRepository) *HttpRequestsManagerBuilder {
	b.manager.templates = repository
	return b
}

func (b *HttpRequestsManagerBuilder) DefaultDelay(duration time.Duration) *HttpRequestsManagerBuilder {
	b.manager.defaultDelay = &duration
	return b
}

func (b *HttpRequestsManagerBuilder) Build() *HttpRequestsManager {
	if b.manager.defaultDelay == nil {
		delay := time.Millisecond
		b.manager.defaultDelay = &delay
	}
	return b.manager
}

func (m *HttpRequestsManager) Execute() {
	for {
		if template, exists := m.templates.Next(); exists {
			if template.Delay != nil {
				if ms, err := strconv.Atoi(*template.Delay); err == nil {
					time.Sleep(time.Duration(ms) * time.Millisecond)
				} else {
					log.Println("Unable to parse delay.", "delay", *template.Delay)
				}
			} else {
				time.Sleep(*m.defaultDelay)
			}
			if template.Log != nil {
				log.Println(*template.Log)
			}
			if response, err := m.send(template); err == nil {
				processResponse(template, &response)
			} else {
				log.Println(err)
			}
		} else {
			break
		}
	}
}

func processResponse(t *templates.Template, r *responseData) {
	if t.Response != nil && t.Response.Log != nil {
		respCode := strconv.Itoa(r.code)
		message := strings.ReplaceAll(*t.Response.Log, "${respCode}", respCode)
		message = strings.ReplaceAll(message, "${respBody}", string(string(r.body)))
		log.Println(message)
	}
}

func (m *HttpRequestsManager) send(t *templates.Template) (respData responseData, err error) {
	httpClient := &http.Client{}
	var request *http.Request
	method := "get"
	if t.Method != nil {
		method = *t.Method
	}
	request, _ = http.NewRequest(method, *t.Target, nil)
	for _, header := range t.Headers {
		request.Header.Add(header.Name, header.Value)
	}
	for _, cookie := range t.Cookies {
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
