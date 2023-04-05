package notify

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"
	t "text/template"
)

func loadTemplate(path string) (template *t.Template, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s template: %s.", path, err)
	}

	template, err = t.New(path).Parse(string(bytes))
	if err != nil {
		return nil, fmt.Errorf("Failed to pass %s Go Template: %s", path, err)
	}
	return
}

func runTemplate(template *template.Template, object interface{}) (result []byte, err error) {
	buf := bytes.Buffer{}
	err = template.Execute(&buf, object)
	if err != nil {
		return
	}
	return buf.Bytes(), nil
}

func testTemplate(template *t.Template, object interface{}) bool {
	result, err := runTemplate(template, object)
	if err != nil {
		log.Print("Failed to run template:", err)
		return false
	}
	return string(result) == "true"
}

func loadAndTestTemplate(path string, object interface{}) bool {
	template, err := loadTemplate(path)
	if err != nil {
		return false
	}
	return testTemplate(template, object)
}
