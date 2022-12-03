package cli

import "github.com/puupee/puupee-api-go"

type ProfileOp struct{}

func NewProfileOp() *ProfileOp {
	return &ProfileOp{}
}

func (op *ProfileOp) Update(value puupee.ProfileDto) error {
	return nil
}
