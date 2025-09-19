package base_controller

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"

	"github.com/miqueaz/GoRestAPI/pkg/client"

	"github.com/gin-gonic/gin"
)

// Métodos CRUD implementados por BaseModelController

// Método Read con soporte de hooks
func (s *Controller[T]) Read(c *gin.Context) {
	// Implementacion del método Read
}

// Insert: Insertar datos en la base de datos
func (s *Controller[T]) Insert(c *gin.Context) {
	// Implementación del método Insert
}

// Update: Actualizar datos en la base de datos
func (s *Controller[T]) Update(c *gin.Context) {
	// Implementación del método Update
}

// Delete: Eliminar datos de la base de datos
func (s *Controller[T]) Delete(c *gin.Context) {
	// Implementación del método Delete
}

// MakeController convierte cualquier función `fn` en un manejador HTTP dinámico para gin.
func MakeController(fn interface{}) gin.HandlerFunc {
	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()
	if fnType.Kind() != reflect.Func {
		panic("MakeController: expected a function")
	}

	if fnType.NumIn() == 1 && fnType.In(0) == reflect.TypeOf((*gin.Context)(nil)) {
		return func(c *gin.Context) {
			fnVal.Call([]reflect.Value{reflect.ValueOf(c)})
		}
	}

	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// ✅ Convertir query params a map[string]interface{} para soportar slices
		queryMap := make(map[string]interface{})
		for k, v := range c.Request.URL.Query() {
			if len(v) == 1 {
				queryMap[k] = v[0]
			} else if len(v) > 1 {
				queryMap[k] = v
			}
		}
		queryBytes, _ := json.Marshal(queryMap)

		argsCount := fnType.NumIn()
		args := make([]reflect.Value, argsCount)
		stringParamIndex := 0

		for i := 0; i < argsCount; i++ {
			paramType := fnType.In(i)

			if paramType == reflect.TypeOf((*gin.Context)(nil)) {
				args[i] = reflect.ValueOf(c)
				continue
			}

			switch paramType.Kind() {
			case reflect.Struct:
				paramPtr := reflect.New(paramType)

				if len(bodyBytes) > 0 {
					if err := json.Unmarshal(bodyBytes, paramPtr.Interface()); err != nil {
						client.Error(c, "invalid request body", err)
						return
					}
				} else if len(queryMap) > 0 {
					if err := json.Unmarshal(queryBytes, paramPtr.Interface()); err != nil {
						client.Error(c, "invalid query params", err)
						return
					}
				}

				args[i] = paramPtr.Elem()

			case reflect.Map:
				paramPtr := reflect.New(paramType).Interface()
				if len(bodyBytes) > 0 {
					if err := json.Unmarshal(bodyBytes, paramPtr); err != nil {
						client.Error(c, "invalid request body (map)", err)
						return
					}
				} else if len(queryMap) > 0 {
					if err := json.Unmarshal(queryBytes, paramPtr); err != nil {
						client.Error(c, "invalid query params (map)", err)
						return
					}
				}
				args[i] = reflect.ValueOf(paramPtr).Elem()

			case reflect.String:
				val := c.Query(paramType.Name())
				if stringParamIndex < len(c.Params) && val == "" {
					val = c.Params[stringParamIndex].Value
				}
				args[i] = reflect.ValueOf(val)
				stringParamIndex++

			case reflect.Int:
				val := c.Query(paramType.Name())
				if stringParamIndex < len(c.Params) && val == "" {
					val = c.Params[stringParamIndex].Value
				}
				if iv, err := strconv.Atoi(val); err == nil {
					args[i] = reflect.ValueOf(iv)
				} else {
					args[i] = reflect.Zero(paramType)
				}
				stringParamIndex++

			case reflect.Float64:
				val := c.Query(paramType.Name())
				if stringParamIndex < len(c.Params) && val == "" {
					val = c.Params[stringParamIndex].Value
				}
				if fv, err := strconv.ParseFloat(val, 64); err == nil {
					args[i] = reflect.ValueOf(fv)
				} else {
					args[i] = reflect.Zero(paramType)
				}
				stringParamIndex++

			default:
				args[i] = reflect.Zero(paramType)
			}
		}

		results := fnVal.Call(args)

		var resp interface{}
		var fnErr error
		switch len(results) {
		case 1:
			if results[0].Type() == reflect.TypeOf((*error)(nil)).Elem() {
				if e, ok := results[0].Interface().(error); ok && e != nil {
					fnErr = e
				}
			}
			resp = results[0].Interface()

		case 2:
			resp = results[0].Interface()
			if e, ok := results[1].Interface().(error); ok && e != nil {
				fnErr = e
			}

		default:
			client.InternalServerError(c, errors.New("MakeController: function must return (T) or (T, error)"))
			return
		}

		if fnErr != nil {
			client.Error(c, "Error in request", fnErr)
			return
		}

		client.Success(c, "OK", resp)
	}
}
