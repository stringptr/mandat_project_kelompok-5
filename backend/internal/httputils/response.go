package httputils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type APIResponse[T any] struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Data    T           `json:"data,omitempty"`
	Errors  []ErrorItem `json:"errors,omitempty"`
	Detail  string      `json:"detail,omitempty"`
	Title   string      `json:"title,omitempty"`
}

type APIResponseOutput[T any] struct {
	Body APIResponse[T]
}

type APIRequestInput[T any] struct {
	Body T
}

func (APIResponse[T]) isUnified() {}

type unifiedMarker interface {
	isUnified()
}

type ErrorItem struct {
	ID       string `json:"id,omitempty"`
	Location string `json:"location,omitempty"`
	Message  string `json:"message"`
	Value    any    `json:"value,omitempty"`
}

func OK[T any](data T) APIResponse[T] {
	return APIResponse[T]{
		Status:  http.StatusOK,
		Success: true,
		Data:    data,
		Detail:  "OK",
		Title:   "Success",
	}
}

type ValidationMapper func(location, message string) (errorID, customMessage string, matched bool)

func Created[T any](data T) APIResponse[T] {
	return APIResponse[T]{
		Status:  http.StatusCreated,
		Success: true,
		Data:    data,
		Detail:  "Created",
		Title:   "Created",
	}
}

func Success[T any](status int, data T, detail string) APIResponse[T] {
	return APIResponse[T]{
		Status:  status,
		Success: true,
		Data:    data,
		Detail:  detail,
		Title:   http.StatusText(status),
	}
}

func Error[T any](status int, detail string, errors []ErrorItem) APIResponse[T] {
	var zero T
	return APIResponse[T]{
		Status:  status,
		Success: false,
		Data:    zero,
		Detail:  detail,
		Title:   http.StatusText(status),
		Errors:  errors,
	}
}

func mapValidationError(location, message string, mappers []ValidationMapper) (string, string) {
	for _, m := range mappers {
		if id, msg, ok := m(location, message); ok {
			return id, msg
		}
	}
	return "", ""
}

func NewUnifiedTransformer(mappers ...ValidationMapper) huma.Transformer {
	return func(ctx huma.Context, status string, v any) (any, error) {
		if _, ok := v.(unifiedMarker); ok {
			return v, nil
		}
		if errModel, ok := v.(*huma.ErrorModel); ok {
			return convertErrorModel(errModel, mappers), nil
		}
		if errModel, ok := v.(huma.ErrorModel); ok {
			return convertErrorModel(&errModel, mappers), nil
		}
		statusCode, _ := strconv.Atoi(status)
		return APIResponse[any]{
			Status:  statusCode,
			Success: statusCode < 400,
			Data:    v,
			Title:   http.StatusText(statusCode),
		}, nil
	}
}

func convertErrorModel(m *huma.ErrorModel, mappers []ValidationMapper) APIResponse[any] {
	var errors []ErrorItem
	if len(m.Errors) > 0 {
		errors = make([]ErrorItem, len(m.Errors))
		for i, e := range m.Errors {
			if e != nil {
				item := ErrorItem{
					Location: e.Location,
					Message:  e.Message,
					Value:    e.Value,
				}
				if m.Status == 422 {
					if errorID, customMsg := mapValidationError(e.Location, e.Message, mappers); errorID != "" {
						item.ID = errorID
						item.Message = customMsg
						item.Location = strings.Join([]string{"request.", e.Location}, "")
					}
				}
				errors[i] = item
			}
		}
	}
	return APIResponse[any]{
		Status:  m.Status,
		Success: false,
		Detail:  m.Detail,
		Title:   m.Title,
		Errors:  errors,
	}
}
