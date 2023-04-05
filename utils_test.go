package notify

import (
	"testing"
	"text/template"
)

func Test_loadTemplate_Valid(t *testing.T) {
	_, err := loadTemplate("test/valid_template.tmpl")
	if err != nil {
		t.Error(err)
	}
}

func Test_loadTemplate_MissingTemplate(t *testing.T) {
	_, err := loadTemplate("missing.tmpl")
	if err == nil {
		t.Error("Template should be missing.")
	}
	t.Log(err)
}

func Test_loadTemplate_InvalidTemplate(t *testing.T) {
	_, err := loadTemplate("test/invalid_template.tmpl")
	if err == nil {
		t.Error("Template should be invalid.")
	}
	t.Log(err)
}

func Test_runTemplate_Success(t *testing.T) {
	const templateString = `{{ eq .Test "test" }}`
	template, _ := template.New("").Parse(templateString)
	result, err := runTemplate(template, struct {
		Test string `json:"test"`
	}{
		Test: "test",
	})
	if err != nil {
		t.Error(err)
	}

	if string(result) != "true" {
		t.Error("Result should be true", string(result))
	}
}

func Test_runTemplate_InvalidObject(t *testing.T) {
	const templateString = `{{ eq .Test "test" }}`
	template, _ := template.New("").Parse(templateString)
	result, err := runTemplate(template, struct {
		Test string `json:"test"`
	}{
		Test: "not_test",
	})

	if err != nil {
		t.Error(err)
	}

	if string(result) != "false" {
		t.Error("Test should be false.", string(result))
	}
	t.Log(string(result))
}

func Test_testTemplate_Success(t *testing.T) {
	const templateString = `{{ eq .Test "test" }}`
	template, _ := template.New("").Parse(templateString)
	success := testTemplate(template, struct {
		Test string `json:"test"`
	}{
		Test: "test",
	})

	if !success {
		t.Error("Should return true")
	}
}

func Test_testTemplate_False(t *testing.T) {
	const templateString = `{{ eq .Test "test" }}`
	template, _ := template.New("").Parse(templateString)
	success := testTemplate(template, struct {
		Test string `json:"test"`
	}{
		Test: "not_test",
	})

	if success {
		t.Error("Should return false")
	}
}

func Test_testTemplate_Fail(t *testing.T) {
	const templateString = `{{ eq .NonExistent "test" }}`
	template, _ := template.New("").Parse(templateString)
	success := testTemplate(template, struct {
		Test string `json:"test"`
	}{
		Test: "test",
	})

	if success {
		t.Error("Should return false")
	}
}
