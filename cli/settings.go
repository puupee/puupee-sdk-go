package cli

import "github.com/puupee/puupee-api-go"

type SettingsOp struct{}

func NewSettingsOp() *SettingsOp {
	return &SettingsOp{}
}

func (op *SettingsOp) Set(value puupee.SettingsDto) error {
	return nil
}

func (op *SettingsOp) Get() (puupee.SettingsDto, error) {
	return puupee.SettingsDto{}, nil
}
