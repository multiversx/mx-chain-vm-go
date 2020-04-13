package nodepart

import (
	"os"
	"path"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

func (driver *ArwenDriver) getArwenPath() (string, error) {
	arwenPath, err := driver.getArwenPathInCurrentDirectory()
	if err == nil {
		return arwenPath, nil
	}

	arwenPath, err = driver.getArwenPathFromEnvironment()
	if err == nil {
		return arwenPath, nil
	}

	return "", common.ErrArwenNotFound
}

func (driver *ArwenDriver) getArwenPathInCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	arwenPath := path.Join(cwd, "arwen")
	if fileExists(arwenPath) {
		return arwenPath, nil
	}

	return "", common.ErrArwenNotFound
}

func (driver *ArwenDriver) getArwenPathFromEnvironment() (string, error) {
	arwenPath := os.Getenv(common.EnvVarArwenPath)
	if fileExists(arwenPath) {
		return arwenPath, nil
	}

	return "", common.ErrArwenNotFound
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
