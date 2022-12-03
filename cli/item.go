package cli

type ItemOp struct {
}

func NewItemOp() *ItemOp {
	return &ItemOp{}
}

func (op *ItemOp) Create() error {
	return nil
}

func (op *ItemOp) Update() error {
	return nil
}

func (op *ItemOp) Delete() error {
	return nil
}

func (op *ItemOp) List() error {
	return nil
}
