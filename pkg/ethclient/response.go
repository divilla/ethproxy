package ethclient

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

var RateLimitErr = errors.New("Rate limiting threshold exceeded, please wait before running more queries")

func parseResponse(json []byte, req *jsonRPCRequest) ([]byte, error) {
	if bytes.Compare(json, []byte(RateLimitErr.Error())) == 0 {
		return nil, RateLimitErr
	}

	errorProp := gjson.GetBytes(json, "error")
	if errorProp.Exists() {
		return nil, fmt.Errorf("json RPC response error: %w", errors.New(errorProp.Raw))
	}

	err := req.compareId(gjson.GetBytes(json, "id").String())
	if err != nil {
		return nil, err
	}

	result := gjson.GetBytes(json, "result")
	if result.IsObject() {
		return []byte(result.Raw), nil
	}

	return []byte(result.String()), nil
}
