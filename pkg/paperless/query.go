package paperless

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/go-logr/logr"
)

type QueryParams struct {
	TruncateContent bool `param:"truncate_content"`
}

type QueryResults struct {
	Results []Document `json:"results,omitempty"`
}

func (clt *Client) QueryDocuments(ctx context.Context, params QueryParams) ([]Document, error) {
	req, err := clt.makeQueryRequest(ctx, params)
	if err != nil {
		return nil, err
	}

	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("Awaiting response")
	resp, err := clt.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}
	log.V(2).Info("Read response", "body", string(b))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s: %s", resp.Status, string(b))
	}

	result := QueryResults{}
	parseErr := json.Unmarshal(b, &result)
	if parseErr != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", parseErr)
	}
	log.V(1).Info("Parsed response", "result", result)
	return result.Results, nil
}

func (clt *Client) makeQueryRequest(ctx context.Context, params QueryParams) (*http.Request, error) {
	log := logr.FromContextOrDiscard(ctx)

	values := paramsToValues(params)
	values.Set("ordering", "id")

	path := clt.URL + "/api/documents/?" + values.Encode()
	log.V(1).Info("Preparing request", "path", path)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request: %w", err)
	}
	clt.setAuth(req)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func paramsToValues(params QueryParams) url.Values {
	values := url.Values{}
	typ := reflect.TypeOf(params)
	value := reflect.ValueOf(params)
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		tag := structField.Tag.Get("param")
		field := value.Field(i)
		paramValue := ""
		switch field.Kind() {
		case reflect.Bool:
			paramValue = strconv.FormatBool(field.Bool())
		case reflect.String:
			paramValue = field.String()
		default:
			panic(fmt.Errorf("not implemented type: %s", field.Kind()))
		}
		values.Set(tag, paramValue)
	}
	return values
}
