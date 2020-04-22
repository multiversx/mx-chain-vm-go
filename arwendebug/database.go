package arwendebug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type database struct {
	rootPath string
}

// NewDatabase -
func NewDatabase(rootPath string) *database {
	_ = os.MkdirAll(rootPath, os.ModePerm)

	return &database{
		rootPath: rootPath,
	}
}

func (db *database) loadWorld(worldID string) (*world, error) {
	var err error
	dataModel := newWorldDataModel(worldID)
	filePath := db.getWorldFile(worldID)

	if fileExists(filePath) {
		dataModel, err = db.readWorldDataModel(filePath)
		if err != nil {
			return nil, err
		}
	}

	world, err := NewWorld(dataModel)
	if err != nil {
		return nil, err
	}

	return world, nil
}

func (db *database) getWorldFile(worldID string) string {
	return path.Join(db.rootPath, fmt.Sprintf("%s.json", worldID))
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func (db *database) readWorldDataModel(filePath string) (*worldDataModel, error) {
	dataModel := &worldDataModel{}
	err := db.unmarshalDataModel(filePath, dataModel)
	if err != nil {
		return nil, err
	}

	return dataModel, nil
}

func (db *database) storeworld() {

}

func (db *database) unmarshalDataModel(filePath string, dataModel interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dataModel)
}

func (db *database) marshalDataModel(dataModel interface{}) ([]byte, error) {
	return json.MarshalIndent(dataModel, "", "\t")
}
