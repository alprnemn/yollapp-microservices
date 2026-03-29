package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

func (g *Gateway) GetUserByUsername() {
	URL := "http://127.0.0.1:8081/users/"
}

func (g *Gateway) Get(ctx context.Context, ID uint32) (*metadataModel.Metadata, error) {

	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	URL := fmt.Sprintf("http://%s/metadata/%d", addrs[rand.Intn(len(addrs))], ID)

	log.Printf("calling metadata service, request get %s", URL)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrorNotFound
	} else if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var v *metadataModel.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil

}
