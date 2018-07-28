package conf

import (
	"os"
	"errors"
)

func GetEnvVariable(name string) (value string, error error) {
	value = os.Getenv(name)
	if len(value) == 0 {
		return "", errors.New("given env variable is not defined")
	}
	return value, nil
}

func SetEnvVriable(name, value string) (err error) {
	err = os.Setenv(name, value)
	if err != nil {
		return err
	}
	return nil
}
