package api

import (
	"fmt"
	"net/http"

	"github.com/kolo/xmlrpc"
)

type XMLRPCEndpointer interface {
	ApiMethod() string
	ApiArgs(string) []interface{}
	Unmarshal([]byte) error
}

type Rpc struct {
	Client     *xmlrpc.Client
	AuthString string
}

func NewClient(url, user, password string, transport http.RoundTripper) (*Rpc, error) {
	client, err := xmlrpc.NewClient(url, transport)
	if err != nil {
		return nil, err
	}
	return &Rpc{
		Client:     client,
		AuthString: fmt.Sprintf("%s:%s", user, password),
	}, nil
}

func (c *Rpc) Call(endpoint XMLRPCEndpointer) error {
	args := endpoint.ApiArgs(c.AuthString)
	method := endpoint.ApiMethod()
	result := []interface{}{}

	if err := c.Client.Call(method, args, &result); err != nil {
		return err
	}

	apiCallSucceeded, ok := result[0].(bool)
	if !ok {
		return fmt.Errorf("malformed XMLRPC response")
	}
	if !apiCallSucceeded {
		return fmt.Errorf(result[0].(string))
	}
	if w, ok := result[1].(string); ok {
		endpoint.Unmarshal([]byte(w))
		return nil
	}
	return fmt.Errorf("no know result type received from RPC call")
}