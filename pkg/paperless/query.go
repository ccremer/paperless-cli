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
	TruncateContent bool   `param:"truncate_content"`
	Ordering        string `param:"ordering"`
	PageSize        int64  `param:"page_size"`
	page            int64  `param:"page"`
}

type QueryResult struct {
	Results []Document `json:"results,omitempty"`
	Next    string     `json:"next,omitempty"`
}

// NextPage returns the next page number for pagination.
// It returns 1 if QueryResult.Next is empty (first page), or 0 if there's an error parsing QueryResult.Next.
func (r QueryResult) NextPage() int64 {
	if r.Next == "" {
		return 1 // first page
	}
	values, err := url.ParseQuery(r.Next)
	if err != nil {
		return 0
	}
	raw := values.Get("page")
	page, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}
	return page
}

func (clt *Client) QueryDocuments(ctx context.Context, params QueryParams) ([]Document, error) {
	documents := make([]Document, 0)
	params.page = 1
	for i := int64(0); i < params.page; i++ {
		result, err := clt.queryDocumentsInPage(ctx, params)
		if err != nil {
			return nil, err
		}
		params.page = result.NextPage()
		documents = append(documents, result.Results...)
	}
	return documents, nil
}

func (clt *Client) makeQueryRequest(ctx context.Context, params QueryParams) (*http.Request, error) {
	log := logr.FromContextOrDiscard(ctx)

	values := paramsToValues(params)

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

func (clt *Client) queryDocumentsInPage(ctx context.Context, params QueryParams) (*QueryResult, error) {
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

	result := QueryResult{}
	parseErr := json.Unmarshal(b, &result)
	if parseErr != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", parseErr)
	}
	log.V(1).Info("Parsed response", "result", result)
	return &result, nil
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
		case reflect.Int64:
			paramValue = strconv.FormatInt(field.Int(), 10)
		default:
			panic(fmt.Errorf("not implemented type: %s", field.Kind()))
		}
		values.Set(tag, paramValue)
	}
	return values
}
