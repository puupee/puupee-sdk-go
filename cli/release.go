package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/puupee/puupee-api-go"
)

type ReleaseOp struct {
	api *puupee.APIClient
}

func NewReleaseOp(api *puupee.APIClient) *ReleaseOp {
	return &ReleaseOp{
		api: api,
	}
}

func (op *ReleaseOp) Create(dto puupee.CreateOrUpdateAppReleaseDto) error {
	resp, _, err := op.api.AppReleaseApi.ApiAppAppReleasePost(context.Background()).Body(dto).Execute()
	if err != nil {
		return err
	}
	PrintObject(resp)
	return nil
}

func (op *ReleaseOp) List(appName string) error {
	appDto, _, err := op.api.AppApi.ApiAppAppByNameGet(context.Background()).Name(appName).Execute()
	if err != nil {
		return err
	}
	dto, _, err := op.api.AppReleaseApi.ApiAppAppReleaseGet(context.Background()).AppId(*appDto.Id).MaxResultCount(100).Execute()
	if err != nil {
		return err
	}
	bts, _ := json.MarshalIndent(dto, "", "  ")
	fmt.Println(string(bts))
	return nil
}
