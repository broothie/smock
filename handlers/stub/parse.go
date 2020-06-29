package stub

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

type UnmarshalFunc func([]byte, interface{}) error

var unmarshalFuncs = map[string]UnmarshalFunc{
	".json": json.Unmarshal,
	".toml": toml.Unmarshal,
	".yml":  yaml.Unmarshal,
	".yaml": yaml.Unmarshal,
}

func Parse(fileOrDir string) (Stubs, error) {
	filenames, err := stubFilenames(fileOrDir)
	if err != nil {
		return nil, err
	}

	stubs := make(Stubs)
	for _, filename := range filenames {
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		ext := filepath.Ext(filename)
		unmarshalFunc, ok := unmarshalFuncs[ext]
		if !ok {
			continue
		}

		newStubs := make(Stubs)
		if err := unmarshalFunc(contents, &newStubs); err != nil {
			return nil, err
		}

		mergeStubs(newStubs, stubs)
	}

	return stubs, nil
}

func stubFilenames(fileOrDir string) ([]string, error) {
	handle, err := os.Stat(fileOrDir)
	if err != nil {
		return nil, err
	}

	if handle.IsDir() {
		return filepath.Glob(filepath.Join(fileOrDir, "**.stubs.*"))
	} else {
		return []string{fileOrDir}, err
	}
}

func mergeStubs(merge, into Stubs) {
	for key, value := range merge {
		into[key] = value
	}
}
