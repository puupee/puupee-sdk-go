package puupeesdk

import (
	"context"

	"golang.org/x/xerrors"

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
	resp, _, err := op.api.AppApi.ApiAppAppPost(context.Background()).
		Body(dto).
		Execute()
	if err != nil {
		return err
	}
	PrintObject(resp)
	return nil
}

func (op *AppOp) List() (*puupee.AppDtoPagedResultDto, error) {
	dto, _, err := op.api.AppApi.ApiAppAppGet(context.Background()).
		MaxResultCount(100).
		Execute()
	if err != nil {
		return nil, xerrors.Errorf("api request error: %w", err)
	}
	return dto, nil
}

func (op *AppOp) Get(appName string) (*puupee.AppDto, error) {
	appDto, _, err := op.api.AppApi.
		ApiAppAppByNameGet(context.Background()).
		Name(appName).
		Execute()
	if err != nil {
		return nil, err
	}
	return appDto, nil
}
