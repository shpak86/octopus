package templates

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
)

type templatesFileBody struct {
	Defaults  HttpRequestTemplate   `json:"defaults"`
	Templates []HttpRequestTemplate `json:"templates"`
}

type TemplatesRepository struct {
	mu        sync.Mutex
	defaults  HttpRequestTemplate
	templates []HttpRequestTemplate
	idx       int
}

// Return next template
func (r *TemplatesRepository) Next() (template *HttpRequestTemplate, exists bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	exists = r.idx != len(r.templates)
	if !exists {
		return
	}
	template = &r.templates[r.idx]
	r.idx++
	return
}

// Return random template
func (r *TemplatesRepository) Random() (template HttpRequestTemplate) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idx = rand.Intn(len(r.templates))
	template = r.templates[r.idx]
	return
}

// Inject variables to the all templates
func (r *TemplatesRepository) inject(variables map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for idx := range r.templates {
		r.templates[idx].Inject(variables)
	}
}

// Templates repository builder helper
type TemplatesRepositoryBuilder struct {
	repo      *TemplatesRepository
	path      string
	variables map[string]string
}

func NewTemplatesRepositoryBuilder() *TemplatesRepositoryBuilder {
	repo := TemplatesRepositoryBuilder{}
	return &repo
}

// Load templates file
func (b *TemplatesRepositoryBuilder) LoadFile(path string) *TemplatesRepositoryBuilder {
	b.path = path
	return b
}

// Inject variables
func (b *TemplatesRepositoryBuilder) Inject(variables map[string]string) *TemplatesRepositoryBuilder {
	b.variables = variables
	return b
}

// Build repository
func (b *TemplatesRepositoryBuilder) Build() *TemplatesRepository {
	b.repo = &TemplatesRepository{}
	if file, err := os.ReadFile(b.path); err == nil {
		var body templatesFileBody
		json.Unmarshal(file, &body)
		b.repo.defaults = body.Defaults
		b.repo.templates = body.Templates
		for idx := 0; idx < len(b.repo.templates); idx++ {
			mergeDefaults(&b.repo.templates[idx], &b.repo.defaults)
		}
	} else {
		log.Fatalln("Unable to build templates", "error", err)
	}
	b.repo.inject(b.variables)
	log.Println("Loaded templates", "size", len(b.repo.templates), "variables", b.variables)
	return b.repo
}

func mergeDefaults(template *HttpRequestTemplate, defaults *HttpRequestTemplate) {
	if defaults.Target != nil && template.Target == nil {
		template.Target = defaults.Target
	}
	if defaults.Cookies != nil && template.Cookies == nil {
		template.Cookies = defaults.Cookies
	}
	if defaults.Headers != nil && template.Headers == nil {
		template.Headers = defaults.Headers
	}
	if defaults.Delay != nil && template.Delay == nil {
		template.Delay = defaults.Delay
	}
	if defaults.Log != nil && template.Log == nil {
		template.Log = defaults.Log
	}
	if defaults.Timeout != nil && template.Timeout == nil {
		template.Timeout = defaults.Timeout
	}
	if defaults.Method != nil && template.Method == nil {
		template.Method = defaults.Method
	}
	if defaults.Response != nil && template.Response == nil {
		template.Response = defaults.Response
	}
}
