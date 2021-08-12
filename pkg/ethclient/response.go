package ethclient

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type response struct {
	json []byte
	err  error
}

func parseResponse(json []byte, req *jsonRPCRequest) ([]byte, error) {
	errorProp := gjson.GetBytes(json,"error")
	if errorProp.Exists() {
		return nil, fmt.Errorf("json RPC response error: %w", errors.New(errorProp.Raw))
	}
	
	err := req.compareId(gjson.GetBytes(json, "id").String())
	if err != nil {
		return nil, err
	}

	result := gjson.GetBytes(json,"result")
	if result.IsObject() {
		return []byte(result.Raw), nil
	}

	return []byte(result.String()), nil
}
