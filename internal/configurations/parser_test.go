package configurations

import (
	"encoding/json"
	"os"
	"testing"
)

// TestLoadFile calls configuration load
func TestLoadFile(t *testing.T) {
	// prep
	temp_file := "tmp.json"
	tf, err := os.CreateTemp("", temp_file)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp_file)

	testConfig := TwitchConfigs{
		TwitchIRL: "irc://twitch.irl.com",
		Channel:   "qwe",
	}
	encoding, _ := json.Marshal(testConfig)
	tf.Write(encoding)

	// test
	config, err := Load(tf.Name())
	if err != nil {
		t.Fatal(err)
	}
	if config.TwitchIRL != "twitch.irl.com" {
		t.Fatalf("IRL not parsed, expected: \"twitch.irl.com\": got %s", config.TwitchIRL)
	}
}

func TestReload(t *testing.T) {
	// prep
	temp_file := "tmp.json"
	tf, err := os.CreateTemp("", temp_file)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp_file)

	testConfig := TwitchConfigs{
		TwitchIRL: "irc://twitch.irl.com",
		Channel:   "qwe",
		Commands: []CommandConfig{
			{
				Name:    "",
				Enable:  true,
				Trigger: "",
			},
		},
		Filters: []FilterConfig{
			{
				Name:    "",
				Enable:  true,
				Pattern: "",
			},
		},
	}
	encoding, _ := json.Marshal(testConfig)
	tf.Write(encoding)

	currentConfig, _ := Load(tf.Name())

	// tests
	t.Run("NoDiferences", func(t *testing.T) {
		changes := currentConfig.Reload(tf.Name())
		if changes == true {
			t.Fail()
		}
	})
	currentConfig.TwitchIRL = "test"
	t.Run("TwitchIRL", func(t *testing.T) {
		changes := currentConfig.Reload(tf.Name())
		if changes == false {
			t.Fail()
		}
	})
	currentConfig.Channel = "test"
	t.Run("Channel", func(t *testing.T) {
		changes := currentConfig.Reload(tf.Name())
		if changes == false {
			t.Fail()
		}
	})
	currentConfig.Debug = true
	t.Run("Debug", func(t *testing.T) {
		changes := currentConfig.Reload(tf.Name())
		if changes == false {
			t.Fail()
		}
	})
	t.Run("Commands", func(t *testing.T) {
		currentConfig.Commands = make([]CommandConfig, 0)
		t.Run("NewConfig", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
		})
		currentConfig.Commands[0].Enable = false
		t.Run("ChangeEnable", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
		})
		currentConfig.Commands[0].Trigger = "test"
		t.Run("ChangeTrigger", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
		})

		currentConfig.Commands[0].Name = "test"
		t.Run("NewName", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
			if len(currentConfig.Commands) != 2 {
				t.Fatal()
			}
		})
	})

	t.Run("Filter", func(t *testing.T) {
		currentConfig.Filters = make([]FilterConfig, 0)
		t.Run("NewFilter", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
		})
		currentConfig.Filters[0].Enable = false
		t.Run("ChangeEnable", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
		})
		currentConfig.Filters[0].Pattern = "test"
		t.Run("ChangeTrigger", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
		})

		currentConfig.Filters[0].Name = "test"
		t.Run("NewName", func(t *testing.T) {
			changes := currentConfig.Reload(tf.Name())
			if changes == false {
				t.Fatal()
			}
			if len(currentConfig.Commands) != 2 {
				t.Fatal()
			}
		})
	})
}
