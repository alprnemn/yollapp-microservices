package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/model"
	"log"
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

func (g *Gateway) RegisterUser(ctx context.Context, user *model.RegisterUserDTO) (*model.RegisterUserResponseDTO, error) {

	URL := fmt.Sprintf("http://127.0.0.1:8081/user/create")

	log.Printf("calling user service, request get %s", URL)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}

		return nil, errors.New(errResp.Error)
	}

	var response model.RegisterUserResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
