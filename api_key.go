package puupeesdk

import (
	"context"

	"github.com/puupee/puupee-api-go"
)

type ApiKeyOp struct {
	api *puupee.APIClient
}

func NewApiKeyOp(api *puupee.APIClient) *ApiKeyOp {
	return &ApiKeyOp{
		api: api,
	}
}

func (op *ApiKeyOp) Create(dto puupee.ApiKeyCreateDto) error {
	resp, _, err := op.api.ApiKeysApi.CreateApiKeys(context.Background()).
		Body(dto).
		Execute()
	if err != nil {
		return err
	}
	PrintObject(resp)
	return nil
}
