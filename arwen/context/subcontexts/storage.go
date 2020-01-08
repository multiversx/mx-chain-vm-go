package subcontexts

type Storage struct {
}

func (storage *Storage) GetStorage(addr []byte, key []byte) []byte {
	panic("not implemented")
}

func (storage *Storage) SetStorage(addr []byte, key []byte, value []byte) int32 {
	panic("not implemented")
}
