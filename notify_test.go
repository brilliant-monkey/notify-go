package notify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/brilliant-monkey/notify-go/config"
	test "github.com/brilliant-monkey/notify-go/test/types"
	"gopkg.in/yaml.v3"
)

func getTestConfig(name string) (config config.NotifierConfig, err error) {
	file, _ := ioutil.ReadFile(fmt.Sprintf("test/%s", name))
	err = yaml.Unmarshal(file, &config)
	return
}

func Test_NewNotifier_Valid(t *testing.T) {
	config, err := getTestConfig("config_valid.yml")
	_, err = NewNotifier(&config)
	if err != nil {
		t.Error(err)
	}
}

func Test_NewNotifier_MissingTemplateConfig(t *testing.T) {
	config, err := getTestConfig("config_missing_template.yml")
	_, err = NewNotifier(&config)
	if err == nil {
		t.Error("Test must be invalid.")
		return
	}
	t.Log("TestNewNotifier_MissingTemplateConfig:", err)
}

func Test_NewNotifier_InvalidRuleConfig(t *testing.T) {
	config, err := getTestConfig("config_invalid_rule.yml")
	_, err = NewNotifier(&config)
	if err == nil {
		t.Error("Test must be invalid.")
		return
	}
	t.Log("TestNewNotifier_InvalidRuleConfig:", err)
}

func Test_NewNotifier_InvalidTemplate(t *testing.T) {
	config, err := getTestConfig("config_invalid_template.yml")
	_, err = NewNotifier(&config)
	if err == nil {
		t.Error("Test must be invalid.")
		return
	}
	t.Log("TestNewNotifier_InvalidTemplateConfig:", err)
}

func Test_RunTemplate_Success(t *testing.T) {
	config, err := getTestConfig("config_valid.yml")
	n, err := NewNotifier(&config)
	if err != nil {
		t.Error(err)
	}

	const jsonString = `{
		"type": "end",
		"after": {
			"label": "person",
			"camera": "northwest",
			"entered_zones": [
				"property"
			]
		}
	}`

	var jsonData test.FrigatePayload
	err = json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		t.Error(err)
	}
	output, err := n.RunTemplate(jsonData)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(output))
	}
}
