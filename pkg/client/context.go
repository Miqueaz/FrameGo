package client

import (
	"context"
	"errors"
)

// func Init(response http.ResponseWriter, request *http.Request) (*HttpContext, error) {
// 	body, err := io.ReadAll(request.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer request.Body.Close()

// 	return &HttpContext{
// 		Request: HttpRequest{
// 			Method: request.Method,
// 			Url:    request.URL.String(),
// 			Header: request.Header,
// 			Params: request.URL.Query(),
// 			Body:   body,
// 		},
// 		Response: HttpResponse{},
// 		Handler: HttpHandler{
// 			Req: request,
// 			Res: response,
// 		},
// 	}, nil
// }

func GetDTO[T any](ctx context.Context, key string) (T, error) {
	var dtoStruct T

	val := ctx.Value(key)
	if val == nil {
		return dtoStruct, errors.New("no dto in context")
	}

	dto, ok := val.(T)
	if !ok {
		return dtoStruct, errors.New("dto has wrong type")
	}

	return dto, nil
}
