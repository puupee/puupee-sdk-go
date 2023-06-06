package puupeesdk

import "github.com/puupee/puupee-api-go"

type TodoOp struct {
	api *puupee.APIClient
}

func NewTodoOp(api *puupee.APIClient) *TodoOp {
	return &TodoOp{
		api: api,
	}
}

func (op *TodoOp) Create() error {
	return nil
}

func (op *TodoOp) Update() error {
	return nil
}

func (op *TodoOp) Delete() error {
	return nil
}
