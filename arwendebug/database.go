package arwendebug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
)

type database struct {
	rootPath string
}

type sessionRecord struct {
	id        string
	createdOn string
	accounts  mock.AccountsMap
}

func newSessionRecord(sessionID string) *sessionRecord {
	return &sessionRecord{
		id:        sessionID,
		createdOn: "now",
		accounts:  make(mock.AccountsMap),
	}
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

	hook := mock.NewBlockchainHookMock()
	hook.Accounts = record.accounts
	session, err := NewSession(hook)
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
	rawData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	record := &sessionRecord{}
	err = json.Unmarshal(rawData, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (db *database) storeSession() {

}
