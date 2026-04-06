package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/model"
	"github.com/alprnemn/yollapp-microservices/shared/errs"
	"io"
	"net/http"
	"time"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

func (g *Gateway) MakeRequest(ctx context.Context, url string, data any, method string, headers map[string]string) (*http.Response, error) {

	var body io.Reader

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshal json gateway -> makeRequest: %s", err.Error())
		}
		body = bytes.NewReader(jsonData)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request gateway -> makeRequest: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error do request gateway -> makeRequest: %s", err.Error())
	}

	return resp, nil
}

func (g *Gateway) RegisterUser(ctx context.Context, user *model.RegisterUserDTO) (*model.RegisterUserResponseDTO, error) {

	resp, err := g.MakeRequest(ctx, "http://127.0.0.1:8081/user/create", user, http.MethodPost, nil)
	if err != nil {
		return nil, err
	}

	return decodeResponse[model.RegisterUserResponseDTO](resp)
}

func (g *Gateway) ActivateUser(ctx context.Context, user *model.ActivateUserDTO) (*model.ActivateResponse, error) {

	resp, err := g.MakeRequest(ctx, "http://127.0.0.1:8081/user/activate", user, http.MethodPatch, nil)
	if err != nil {
		return nil, err
	}

	return decodeResponse[model.ActivateResponse](resp)
}

func decodeResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {

		var errResp errs.ErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error decoding json gateway -> decoderesponse: %s", err.Error())
		}

		return nil, &errResp
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding json gateway result  -> decoderesponse: %s", err.Error())
	}

	return &result, nil
}
