package request

import (
	"encoding/json"
	"fmt"
	errTypes "github.com/charmLd/token-generator-api/http/error/types"
	"net/http"
)

// Unpack the request in to the given unpacker struct.
func Unpack(r *http.Request, unpacker UnpackerInterface) error {

	err := json.NewDecoder(r.Body).Decode(unpacker)
	if err != nil {

		verr := errTypes.ValidationError{}
		return verr.New(fmt.Sprintf("Required format: %s", unpacker.RequiredFormat()))
	}
	return nil

}
