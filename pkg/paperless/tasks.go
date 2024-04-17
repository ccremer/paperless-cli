package paperless

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

type Task struct {
	ID              int    `json:"id,omitempty"`
	TaskId          string `json:"task_id,omitempty"`
	RelatedDocument string `json:"related_document,omitempty"`
	Status          string `json:"status,omitempty"`
}

func (clt *Client) QueryTasks(ctx context.Context, taskUuid string) ([]Task, error) {
	req, err := clt.makeTaskQueryRequest(ctx, taskUuid)
	if err != nil {
		return nil, err
	}

	result := make([]Task, 0)
	if err := clt.executeQueryAndParse(ctx, req, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (clt *Client) makeTaskQueryRequest(ctx context.Context, taskUuid string) (*http.Request, error) {
	log := logr.FromContextOrDiscard(ctx)

	query := clt.URL + "/api/tasks/"
	if taskUuid != "" {
		query += "?task_id=" + taskUuid
	}
	log.V(1).Info("Preparing request", "query", query)
	req, err := http.NewRequestWithContext(ctx, "GET", query, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request: %w", err)
	}
	clt.setAuth(req)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
