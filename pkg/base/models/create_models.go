package base_models

import (
	"database/sql"
	"log"

	helpers "github.com/miqueaz/FrameGo/pkg/base/helpers"
	ORM "github.com/miqueaz/FrameGo/pkg/sql"

	"github.com/jmoiron/sqlx"
)

var DB *sql.DB

// Crear un nuevo modelo y guardarlo
func NewModel[T any](name string, collectionName string, id ...int) *Model[T] {
	db := sqlx.NewDb(DB, "postgres")
	var idInt int
	if len(id) > 0 {
		idInt = id[0]
	}
	model := &Model[T]{
		ID:             idInt,
		Name:           name,
		CollectionName: collectionName,
		Structure:      *new(T),
		QueryBuilder:   ORM.NewQueryBuilder[T](db, collectionName),
	}
	helpers.SaveStructure(model, &models)

	if m, ok := helpers.LoadStructure[Model[T]](&models); ok {
		log.Printf("Modelo '%s' creado y almacenado con éxito.\n", m.Name)
		return model
	}

	return nil
}

func SetGlobalDB(db *sql.DB) {
	DB = db
}

func (m *Model[T]) SetDB(db *sqlx.DB) {
	m.QueryBuilder = ORM.NewQueryBuilder[T](db, m.CollectionName)
}

func SetDB(db *sqlx.DB) {
	models.Range(func(_, value any) bool {
		if model, ok := value.(interface{ SetDB(*sqlx.DB) }); ok {
			model.SetDB(db)
		}
		return true
	})
	log.Println("Base de datos configurada para todos los modelos.")
}
