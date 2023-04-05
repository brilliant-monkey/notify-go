package notify

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"golang.org/x/exp/slices"
)

type Notifier struct {
	config    *NotifierConfig
	templates map[string]*template.Template
}

type NotifierConfig struct {
	Templates map[string]string `yaml:"templates"`
}

func contains(arr []string, value string) bool {
	return slices.Contains(arr, value)
}

func NewNotifier(config *NotifierConfig) (*Notifier, error) {
	templates, err := buildTemplates(config)
	if err != nil {
		log.Fatalln("Could not build templates:", err)
	}
	err = validateTemplates(templates)
	if err != nil {
		log.Fatalln("Templates are invalid:", err)
	}

	return &Notifier{
		config,
		templates,
	}, nil
}

func buildTemplates(config *NotifierConfig) (templates map[string]*template.Template, err error) {
	funcs := map[string]any{
		"contains": contains,
	}

	for key, value := range config.Templates {
		template, err := template.New(key).Funcs(funcs).Parse(value)
		if err != nil {
			return nil, err
		}

		templates[key] = template
	}
	return
}

func readTemplate(name string) (t *template.Template, err error) {
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}

	t, err = template.New(name).Parse(string(bytes))
	if err != nil {
		return
	}
	return
}

func validateTemplates(templates map[string]*template.Template) (err error) {
	for key, _ := range templates {
		_, err = readTemplate(key)
	}
	return
}

func (notifier *Notifier) testTemplate(name string, object interface{}) bool {
	_, err := runTemplate(notifier.templates[name], object)
	if err != nil {
		return false
	}
	return true
}

func runTemplate(template *template.Template, object interface{}) (result []byte, err error) {
	buf := bytes.Buffer{}
	err = template.Execute(&buf, object)
	if err != nil {
		return
	}
	return buf.Bytes(), nil
}

func (notifier *Notifier) RunTemplate(object interface{}) (output []byte, err error) {
	for key := range notifier.templates {
		isMatch := notifier.testTemplate(key, object)
		if isMatch {
			template, err := readTemplate(key)
			result, err := runTemplate(template, object)
			if err != nil {
				log.Println("Error occured reading template: ")
				break
			}
			return result, nil
		}
	}
	return nil, fmt.Errorf("No matching templates for execution.")
}
