package hooks

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Cleaner interface {
	FindClenear(data *mongo.Cursor, err error) ([]map[string]any, error)
}

type Cleaners struct{}

func (c *Cleaners) FindClenear(data *mongo.Cursor, err error) ([]map[string]any, error) {
	var dataConverted []map[string]any
	ctx := context.TODO()
	data.All(ctx, &dataConverted)
	return dataConverted, err

}
