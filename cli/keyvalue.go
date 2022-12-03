package cli

type KeyValueOp struct {
}

func NewKeyValueOp() *KeyValueOp {
	return &KeyValueOp{}
}

func (op *KeyValueOp) Set(key string, value interface{}) error {
	return nil
}

func (op *KeyValueOp) Get(key string) error {
	return nil
}
