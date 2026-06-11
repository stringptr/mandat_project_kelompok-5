package httputils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type APIResponse[T any] struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Data    T            `json:"data,omitempty"`
	Errors  []*ErrorItem `json:"errors,omitempty"`
	Detail  string       `json:"detail,omitempty"`
	Title   string       `json:"title,omitempty"`
}

type APIResponseOutput[T any] struct {
	Status int            `json:"-"`
	Body   APIResponse[T]
}

func NewOKOutput[T any](data T) *APIResponseOutput[T] {
	body := OK(data)
	return &APIResponseOutput[T]{Status: body.Status, Body: body}
}

func NewCreatedOutput[T any](data T) *APIResponseOutput[T] {
	body := Created(data)
	return &APIResponseOutput[T]{Status: body.Status, Body: body}
}

func NewSuccessOutput[T any](status int, data T, detail string) *APIResponseOutput[T] {
	body := Success(status, data, detail)
	return &APIResponseOutput[T]{Status: body.Status, Body: body}
}

type APIRequestInput[T any] struct {
	Body T
}

func (APIResponse[T]) isUnified() {}

type unifiedMarker interface {
	isUnified()
}

type ErrorItem struct {
	ID       string `json:"id"`
	Location string `json:"location,omitempty"`
	Message  string `json:"message"`
	Value    any    `json:"value,omitempty"`
}

func (e *ErrorItem) Error() string { return e.Message }
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

func Error[T any](status int, detail string, errors []*ErrorItem) APIResponse[T] {
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

func normalizeLocation(location string) string {
	return strings.TrimPrefix(location, "request.")
}

func mapValidationError(location, message string, mappers []ValidationMapper) (string, string) {
	loc := normalizeLocation(location)
	for _, m := range mappers {
		if id, msg, ok := m(loc, message); ok {
			return id, msg
		}
	}
	return "", ""
}

func convertValidationError(e *ValidationError) APIResponse[any] {
	return APIResponse[any]{
		Status:  e.StatusCode,
		Success: false,
		Detail:  e.Detail,
		Title:   http.StatusText(e.StatusCode),
		Errors:  e.Errors,
	}
}

func NewUnifiedTransformer(mappers ...ValidationMapper) huma.Transformer {
	return func(ctx huma.Context, status string, v any) (any, error) {
		switch t := v.(type) {
		case unifiedMarker:
			return v, nil
		case *ValidationError:
			return convertValidationError(t), nil
		case *huma.ErrorModel:
			return convertErrorModel(t, mappers), nil
		case huma.ErrorModel:
			return convertErrorModel(&t, mappers), nil
		default:
			statusCode, _ := strconv.Atoi(status)
			return APIResponse[any]{
				Status:  statusCode,
				Success: statusCode < 400,
				Data:    v,
				Title:   http.StatusText(statusCode),
			}, nil
		}
	}
}

func convertErrorModel(m *huma.ErrorModel, mappers []ValidationMapper) APIResponse[any] {
	var errors []*ErrorItem
	if len(m.Errors) > 0 {
		errors = make([]*ErrorItem, len(m.Errors))
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
					}
				}
				errors[i] = &item
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

type ValidationError struct {
	StatusCode int
	Detail     string
	Errors     []*ErrorItem
}

func (e *ValidationError) Error() string  { return e.Detail }
func (e *ValidationError) GetStatus() int { return e.StatusCode }
func (e *ValidationError) ContentType(ct string) string {
	if ct == "application/json" {
		return "application/json"
	}
	return ct + "+json"
}

func (e *ValidationError) ErrorDetail() *huma.ErrorDetail {
	if len(e.Errors) > 0 {
		return &huma.ErrorDetail{
			Location: e.Errors[0].Location,
			Message:  e.Errors[0].Message,
			Value:    e.Errors[0].Value,
		}
	}
	return &huma.ErrorDetail{Message: e.Detail}
}
