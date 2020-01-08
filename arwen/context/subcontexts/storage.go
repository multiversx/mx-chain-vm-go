package subcontexts

type Storage struct {
}

func (s *Storage) GetStorage(addr []byte, key []byte) []byte {
	panic("not implemented")
}

func (s *Storage) SetStorage(addr []byte, key []byte, value []byte) int32 {
	panic("not implemented")
}
