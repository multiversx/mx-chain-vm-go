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

func (db *database) loadSession(sessionID string) (*session, error) {
	var err error
	record := newSessionRecord(sessionID)
	filePath := db.getSessionFile(sessionID)

	if fileExists(filePath) {
		record, err = db.readSessionRecord(filePath)
		if err != nil {
			return nil, err
		}
	}

	session, err := NewSession(record)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (db *database) getSessionFile(sessionID string) string {
	return path.Join(db.rootPath, fmt.Sprintf("%s.json", sessionID))
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func (db *database) readSessionRecord(filePath string) (*sessionRecord, error) {
	record := &sessionRecord{}
	err := db.unmarshalRecord(filePath, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (db *database) storeSession() {

}

func (db *database) unmarshalRecord(filePath string, record interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, record)
}

func (db *database) marshalRecord(record interface{}) ([]byte, error) {
	return json.MarshalIndent(record, "", "\t")
}
