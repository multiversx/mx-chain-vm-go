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

// NewDatabase -
func NewDatabase(rootPath string) *database {
	_ = os.MkdirAll(rootPath, os.ModePerm)

	return &database{
		rootPath: rootPath,
	}
}

func (db *database) loadSession(sessionID string) (*session, error) {
	filePath := path.Join(db.rootPath, fmt.Sprintf("%s.json", sessionID))
	rawData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	record := &sessionRecord{}
	err = json.Unmarshal(rawData, record)
	if err != nil {
		return nil, err
	}

	hook := mock.NewBlockchainHookMock()
	hook.Accounts = record.accounts
	session := NewSession(hook)
	return session, nil
}

func (db *database) storeSession() {

}
