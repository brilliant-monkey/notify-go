package notify

import (
	"fmt"
	"text/template"

	"github.com/brilliant-monkey/notify-go/config"
	"golang.org/x/exp/slices"
)

type Notifier struct {
	config    *config.NotifierConfig
	templates map[string]*template.Template
}

func contains(arr []string, value string) bool {
	return slices.Contains(arr, value)
}

func NewNotifier(config *config.NotifierConfig) (*Notifier, error) {
	templates, err := buildTemplates(config)
	if err != nil {
		return nil, fmt.Errorf("Could not build templates: %s", err)
	}
	err = validateTemplates(templates)
	if err != nil {
		return nil, fmt.Errorf("Template validation failed: %s", err)
	}

	return &Notifier{
		config,
		templates,
	}, nil
}

func buildTemplates(config *config.NotifierConfig) (templates map[string]*template.Template, err error) {
	funcs := map[string]any{
		"contains": contains,
	}

	for key, value := range config.Templates {
		templates = map[string]*template.Template{}
		template, err := template.New(key).Funcs(funcs).Parse(value)
		if err != nil {
			return nil, err
		}
		templates[key] = template
	}
	return
}

func validateTemplates(templates map[string]*template.Template) (err error) {
	for key, _ := range templates {
		_, err = loadTemplate(key)
	}
	return
}

func (notifier *Notifier) RunTemplate(object interface{}) (output []byte, err error) {
	for key, value := range notifier.templates {
		if testTemplate(value, object) {
			template, err := loadTemplate(key)
			if err != nil {
				return nil, err
			}
			return runTemplate(template, object)
		}
	}
	return nil, fmt.Errorf("No matching templates for execution.")
}
