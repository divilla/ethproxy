package ethclient

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/sjson"
)

type jsonRPCRequest struct {
	json string
	uuid string
}

func request(method string) *jsonRPCRequest {
	uid := uuid.New().String()
	json := `{"jsonrpc":"2.0","params":[]}`

	json, err := sjson.Set(json, "method", "eth_" + method)
	if err != nil {
		panic(err)
	}

	json, err = sjson.Set(json, "id", uid)
	if err != nil {
		panic(err)
	}

	return &jsonRPCRequest{
		json: json,
		uuid: uid,
	}
}

func (r *jsonRPCRequest) param(value interface{}) *jsonRPCRequest {
	json, err := sjson.Set(r.json, "params.-1", value)
	if err != nil {
		panic(err)
	}

	r.json = json

	return r
}

func (r *jsonRPCRequest) compareId(uid string) error {
	if r.uuid != uid {
		return fmt.Errorf("json rpc id's don't match request: '%s', response: '%s'", r.uuid, uid)
	}

	return nil
}

func (r *jsonRPCRequest) String() string {
	return r.json
}
