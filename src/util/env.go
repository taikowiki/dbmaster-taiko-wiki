package util

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

func getEnvDir() string {
	var envDir string
	if *Flag.UseCwd {
		cwd, err := os.Getwd()
		if err != nil {
			return ""
		}

		envDir = cwd
	} else {
		execPath, err := os.Executable()
		if err != nil {
			return ""
		}

		execDir := filepath.Dir(execPath)
		envDir = execDir
	}
	return envDir
}

func LoadEnv() (map[string]string, error) {
	env := make(map[string]string)

	envPath := filepath.Join(getEnvDir(), ".env.json")
	envJson, err := os.ReadFile(envPath)
	if err != nil {
		return env, err
	}

	err = json.Unmarshal(envJson, &env)
	if err != nil {
		return env, err
	}

	return env, nil
}

func LoadConnDatas() ([]types.DBConnectionData, error) {
	connDatasPath := filepath.Join(getEnvDir(), "connDatas.env.json")
	jsonContent, err := os.ReadFile(connDatasPath)
	if err != nil {
		return nil, err
	}

	var connDatas []types.DBConnectionData
	err = json.Unmarshal(jsonContent, &connDatas)
	if err != nil {
		return nil, err
	}
	return connDatas, nil
}
