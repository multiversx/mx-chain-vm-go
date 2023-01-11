package vmserver

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

// newDatabase creates a new debugging database (basically, a folder with JSON files)
func newDatabase(rootPath string) *database {
	db := &database{rootPath: rootPath}
	db.initFolders()
	return db
}

func (db *database) initFolders() {
	err := os.MkdirAll(path.Join(db.rootPath, "worlds"), os.ModePerm)
	if err != nil {
		log.Error("database.initFolders", "err", err)
	}

	err = os.MkdirAll(path.Join(db.rootPath, "out"), os.ModePerm)
	if err != nil {
		log.Error("database.initFolders", "err", err)
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

	world, err := newWorld(dataModel)
	if err != nil {
		return nil, err
	}

	return world, nil
}

func (db *database) getWorldFile(worldID string) string {
	return path.Join(db.rootPath, "worlds", fmt.Sprintf("%s.json", worldID))
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

func (db *database) storeWorld(world *world) error {
	filePath := db.getWorldFile(world.id)
	log.Trace("Database.storeWorld()", "file", filePath)

	dataModel := world.toDataModel()
	return db.marshalDataModel(filePath, dataModel)
}

func (db *database) storeOutcome(key string, outcome interface{}) error {
	if len(key) == 0 {
		log.Trace("Database.storeOutcome(), won't store (empty key)")
		return nil
	}

	filePath := db.getOutcomeFile(key)
	log.Trace("Database.storeOutcome()", "file", filePath)
	return db.marshalDataModel(filePath, outcome)
}

func (db *database) getOutcomeFile(uniqueID string) string {
	return path.Join(db.rootPath, "out", fmt.Sprintf("%s.json", uniqueID))
}

func (db *database) unmarshalDataModel(filePath string, dataModel interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dataModel)
}

func (db *database) marshalDataModel(filePath string, dataModel interface{}) error {
	data, err := json.MarshalIndent(dataModel, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}
