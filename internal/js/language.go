package js

import (
	"fmt"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx", "d.ts", "mjs", "cjs",
}

type Language struct {
	Workspaces *Workspaces
	Cfg        *Config
}

var _ language.Language[js_grammar.File] = &Language{}

// findPackageJson starts from a search path and goes up dir by dir
// until a package.json file is found. If one is found, it returns the
// dir where it was found and a parsed TsConfig object in case that there
// was also a tsconfig.json file.
func _findPackageJson(searchPath string) (TsConfig, string, error) {
	packageJsonPath := filepath.Join(searchPath, "package.json")
	if utils.FileExists(packageJsonPath) {
		tsConfigPath := filepath.Join(searchPath, "tsconfig.json")
		var tsConfig TsConfig
		var err error
		if utils.FileExists(tsConfigPath) {
			tsConfig, err = ParseTsConfig(tsConfigPath)
			if err != nil {
				err = fmt.Errorf("found TypeScript config file in %s but there was an error reading it: %w", tsConfigPath, err)
			}
		}
		return tsConfig, searchPath, err
	}
	nextSearchPath := filepath.Dir(searchPath)
	if nextSearchPath != searchPath {
		return _findPackageJson(nextSearchPath)
	} else {
		return TsConfig{}, "", nil
	}
}

var findPackageJson = utils.Cached1In2OutErr(_findPackageJson)

func MakeJsLanguage(entrypoint string, cfg *Config) (language.Language[js_grammar.File], error) {
	if !utils.FileExists(entrypoint) {
		return nil, fmt.Errorf("file %s does not exist", entrypoint)
	}
	workspaces, err := NewWorkspaces(entrypoint)
	if err != nil {
		return nil, err
	}

	return &Language{
		Cfg:        cfg,
		Workspaces: workspaces,
	}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
