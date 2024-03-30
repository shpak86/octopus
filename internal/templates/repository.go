package templates

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
)

type templatesFileBody struct {
	Templates []Template `json:"templates"`
}

type TemplatesRepository struct {
	mu        sync.Mutex
	templates []Template
	idx       int
}

func (r *TemplatesRepository) Next() (template *Template, exists bool) {
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

func (r *TemplatesRepository) Random() (template Template) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idx = rand.Intn(len(r.templates))
	template = r.templates[r.idx]
	return
}

func (r *TemplatesRepository) inject(variables map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for idx := range r.templates {
		r.templates[idx].Inject(variables)
	}
}

type TemplatesRepositoryBuilder struct {
	repo      *TemplatesRepository
	path      string
	variables map[string]string
}

func NewTemplatesRepositoryBuilder() *TemplatesRepositoryBuilder {
	repo := TemplatesRepositoryBuilder{}
	return &repo
}

func (b *TemplatesRepositoryBuilder) LoadFile(path string) *TemplatesRepositoryBuilder {
	b.path = path
	return b
}

func (b *TemplatesRepositoryBuilder) Inject(variables map[string]string) *TemplatesRepositoryBuilder {
	b.variables = variables
	return b
}

func (b *TemplatesRepositoryBuilder) Build() *TemplatesRepository {
	b.repo = &TemplatesRepository{}
	if file, err := os.ReadFile(b.path); err == nil {
		var body templatesFileBody
		json.Unmarshal(file, &body)
		b.repo.templates = body.Templates
	} else {
		log.Fatalln("Unable to build templates", "error", err)
	}
	b.repo.inject(b.variables)
	log.Println("Loaded templates", "size", len(b.repo.templates), "variables", b.variables)
	return b.repo
}
