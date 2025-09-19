package base_helpers

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Función para convertir BSON a `map[string]interface{}` y manejar `primitive.ObjectID`
func ConvertBsonToMap(doc bson.M) map[string]any {
	return convertBsonValue(doc).(map[string]any)
}

func convertBsonValue(value any) any {
	switch v := value.(type) {
	case primitive.ObjectID:
		return v.Hex()
	case primitive.DateTime:
		// Convertir a string en formato RFC3339 compatible con proto
		return v.Time().Format(time.RFC3339)
	case primitive.Timestamp:
		// Devolver el timestamp como entero (o combínalo como string si lo prefieres)
		return int64(v.T)
	case primitive.A:
		arr := make([]interface{}, len(v))
		for i, item := range v {
			arr[i] = convertBsonValue(item)
		}
		return arr
	case bson.M:
		mapped := make(map[string]interface{})
		for k, val := range v {
			mapped[k] = convertBsonValue(val)
		}
		return mapped
	case map[string]interface{}:
		mapped := make(map[string]interface{})
		for k, val := range v {
			mapped[k] = convertBsonValue(val)
		}
		return mapped
	default:
		return v
	}
}

func ConvertMapToStrings(input map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string)
	for key, value := range input {
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("el valor para la clave %s no es un string", key)
		}
		result[key] = strValue
	}
	return result, nil
}

func NormalizarFiltros(input map[string]any) map[string][]any {
	normalizados := make(map[string][]any)
	for campo, valor := range input {
		// Si ya viene como arreglo (ej: [">", 3, "<", 10])
		if arr, ok := valor.([]any); ok && len(arr) > 0 {
			normalizados[campo] = arr
			continue
		}
		// Si viene como valor simple, lo convertimos a ["=", valor]
		normalizados[campo] = []any{"=", valor}
	}
	return normalizados
}
