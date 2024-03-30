package requests

import (
	"octopus/internal/templates"
)

type responseData struct {
	code int
	body []byte
}

type HttpRequestsManager struct {
	templates    *templates.TemplatesRepository
	threads      int
	requestsChan chan templates.Template
}

func (m *HttpRequestsManager) Execute() {
	m.requestsChan = make(chan templates.Template, m.threads)
	for t := 0; t < m.threads; t++ {
		sender := *NewRequestSender(m.requestsChan)
		sender.Serve()
	}
	template, exists := m.templates.Next()
	for exists {
		m.requestsChan <- *template
		template, exists = m.templates.Next()
	}
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

func (b *HttpRequestsManagerBuilder) Parallelism(threads int) *HttpRequestsManagerBuilder {
	b.manager.threads = threads
	return b
}

func (b *HttpRequestsManagerBuilder) Build() *HttpRequestsManager {
	if b.manager.threads <= 0 {
		b.manager.threads = 1
	}
	b.manager.requestsChan = make(chan templates.Template, b.manager.threads)
	return b.manager
}
