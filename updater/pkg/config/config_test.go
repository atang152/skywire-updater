package config

import (
	"testing"

	"github.com/skycoin/services-bk/autoupdater/config"
	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	var expectedServiceMaps = map[string]config.ServiceConfig{
		"maria": {
			OfficialName:         "maria",
			LocalName:            "mariadb",
			UpdateScript:         "$HOME/skycoin/github.com/watercompany/skywire-services/updater/updater/custom_example/custom_script.sh",
			ScriptInterpreter:    "/bin/bash",
			ScriptTimeout:        "20s",
			ScriptExtraArguments: []string{"-a 1"},
			ActiveUpdateChecker:  "dockerhub_fetcher",
			Repository:           "/library/mariadb",
			CheckTag:             "latest",
			Updater:              "custom",
		},
		"top": {
			OfficialName:         "top",
			LocalName:            "skywire",
			PassiveUpdateChecker: "nats",
			Updater:              "swarm",
		},
		"sky-node": {
			OfficialName:         "sky-node",
			LocalName:            "skycoin",
			UpdateScript:         "/Users/ivan/Desktop/skycoin/pkg/github.com/watercompany/skywire-services/updater/pkg/updater/scripts/skycoin.sh",
			ScriptInterpreter:    "/bin/bash",
			ScriptTimeout:        "20m",
			ScriptExtraArguments: []string{"-a 1", "-b 2"},
			ActiveUpdateChecker:  "git_fetcher_1",
			Repository:           "/skycoin/skycoin",
			CheckTag:             "latest",
			Updater:              "custom",
		},
		"skywire": {
			OfficialName:        "skywire",
			LocalName:           "mystack_skywire",
			UpdateScript:        "./updater/custom_example/custom_script.sh",
			ScriptInterpreter:   "/bin/bash",
			ActiveUpdateChecker: "dockerhub_fetch",
			CheckTag:            "latest",
			Updater:             "swarm",
		},
	}

	c := config.NewFromFile("../configuration.example.yml")

	assert.Equal(t, expectedServiceMaps, c.Services)
}

func TestUpdaters(t *testing.T) {
	var expectedUpdaters = map[string]config.UpdaterConfig{
		"custom": {
			Kind:      "custom",
			Notify:    true,
			NotifyUrl: "http://manager.ui:9000",
		},
		"swarm": {
			Kind: "swarm",
		},
	}

	c := config.NewFromFile("../configuration.example.yml")

	assert.Equal(t, expectedUpdaters, c.Updaters)
}

func TestActiveUpdateChekers(t *testing.T) {
	var expectedActiveUpdateChekcers = map[string]config.FetcherConfig{
		"git_fetcher_1": {
			Kind:      "git",
			Interval:  "30s",
			RetryTime: "22s",
			Retries:   3,
		},
		"dockerhub_fetcher": {
			Kind:      "dockerhub",
			Interval:  "30s",
			RetryTime: "22s",
		},
	}

	c := config.NewFromFile("../configuration.example.yml")

	assert.Equal(t, expectedActiveUpdateChekcers, c.ActiveUpdateCheckers)
}

func TestPassiveUpdateCheckers(t *testing.T) {
	var expectedPassiveUpdateCheckers = map[string]config.SubscriberConfig{
		"nats": {
			MessageBroker: "nats",
			Urls:          []string{"http://localhost:4222"},
			Topic:         "top",
		},
	}

	c := config.NewFromFile("../configuration.example.yml")

	assert.Equal(t, expectedPassiveUpdateCheckers, c.PassiveUpdateCheckers)
}
