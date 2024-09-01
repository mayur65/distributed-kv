package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func createConfig(t *testing.T, contents string) Config {

	t.Helper()

	f, err := ioutil.TempFile(os.TempDir(), "sharding.toml")

	if err != nil {
		t.Fatalf("File creation failed: %v", err)
	}

	defer f.Close()

	name := f.Name()
	defer os.Remove(name)

	if _, err := f.WriteString(contents); err != nil {
		t.Fatalf("File write failed: %v", err)
	}

	c, err := ParseConfigFile(name)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	return *c
}

func TestConfigParse(t *testing.T) {

	actual := createConfig(t, ` [[shards]]
	name = "Mayur"
	id = 0
	addr = "localhost:8080"`)

	expected := Config{
		Shards: []Shard{
			{
				Name: "Mayur",
				Id:   0,
				Addr: "localhost:8080",
			},
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Config does not match: %#v %#v", expected, actual)
	}
}
