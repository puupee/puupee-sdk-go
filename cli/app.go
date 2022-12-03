package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/puupee/puupee-api-go"
)

type AppOp struct {
	api *puupee.APIClient
}

func NewAppOp(api *puupee.APIClient) *AppOp {
	return &AppOp{
		api: api,
	}
}

func (op *AppOp) Create(dto puupee.CreateOrUpdateAppDto) error {
	resp, _, err := op.api.AppApi.ApiAppAppPost(context.Background()).Body(dto).Execute()
	if err != nil {
		return err
	}
	PrettyPrint(resp)
	return nil
}

func (op *AppOp) List() error {
	dto, _, err := op.api.AppApi.ApiAppAppGet(context.Background()).MaxResultCount(100).Execute()
	if err != nil {
		return err
	}
	bts, _ := json.MarshalIndent(dto, "", "  ")
	fmt.Println(string(bts))
	return nil
}

func (op *AppOp) Get(appName string) error {
	appDto, _, err := op.api.AppApi.ApiAppAppByNameGet(context.Background()).Name(appName).Execute()
	if err != nil {
		return err
	}
	bts, _ := json.MarshalIndent(appDto, "", "  ")
	fmt.Println(string(bts))
	return nil
}
