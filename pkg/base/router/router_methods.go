package router

import (
	"errors"

	base_controller "github.com/miqueaz/FrameGo/pkg/base/controller"
	"github.com/miqueaz/FrameGo/pkg/client"

	"github.com/gin-gonic/gin"
)

func Router() *AppRouter {
	// Usamos gin.Default() para crear el enrutador con middleware por defecto (como el logger y recovery
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(func(c *gin.Context) {
		client.NotFound(c, errors.New("Route not found"))
	})
	return &AppRouter{r}
}

func (ar *AppRouter) Use(mw ...gin.HandlerFunc) {
	// Añadimos el middleware a la lista de middlewares
	ar.Engine.Use(mw...)
}

func (ar *AppRouter) GET(path string, fn any, mw ...gin.HandlerFunc) {
	// En gin, las rutas se definen con métodos como GET, POST, etc.
	ar.Engine.GET(path, append(mw, base_controller.MakeController(fn))...)
}

func (ar *AppRouter) POST(path string, fn any, mw ...gin.HandlerFunc) {
	// En gin, las rutas se definen con métodos como GET, POST, etc.
	ar.Engine.POST(path, append(mw, base_controller.MakeController(fn))...)
}

func (ar *AppRouter) PUT(path string, fn any, mw ...gin.HandlerFunc) {
	// En gin, las rutas se definen con métodos como GET, POST, etc.
	ar.Engine.PUT(path, append(mw, base_controller.MakeController(fn))...)
}

func (ar *AppRouter) DELETE(path string, fn any, mw ...gin.HandlerFunc) {
	// En gin, las rutas se definen con métodos como GET, POST, etc.
	ar.Engine.DELETE(path, append(mw, base_controller.MakeController(fn))...)
}

func (ar *AppRouter) Execute(addr string) error {
	// Inicia el servidor en la dirección especificada
	return ar.Engine.Run(addr)
}

func (ar *AppRouter) Group(prefix string) *GroupRouter {
	// Crea un nuevo grupo de rutas con un prefijo y middlewares específicos
	group := ar.Engine.Group(prefix)
	return &GroupRouter{Router: group}
}

func (gr *GroupRouter) GET(path string, fn any, mw ...gin.HandlerFunc) {
	// Si hay middlewares adicionales, los aplicamos solo a esta ruta
	gr.Router.GET(path, append(mw, base_controller.MakeController(fn))...)
}

func (gr *GroupRouter) POST(path string, fn any, mw ...gin.HandlerFunc) {
	// Si hay middlewares adicionales, los aplicamos solo a esta ruta
	gr.Router.POST(path, append(mw, base_controller.MakeController(fn))...)
}
func (gr *GroupRouter) PUT(path string, fn any, mw ...gin.HandlerFunc) {
	// Si hay middlewares adicionales, los aplicamos solo a esta ruta
	gr.Router.PUT(path, append(mw, base_controller.MakeController(fn))...)
}

func (gr *GroupRouter) DELETE(path string, fn any, mw ...gin.HandlerFunc) {
	// Si hay middlewares adicionales, los aplicamos solo a esta ruta
	gr.Router.DELETE(path, append(mw, base_controller.MakeController(fn))...)
}

func (gr *GroupRouter) USE(mw ...gin.HandlerFunc) {
	// Añade middlewares al grupo de rutas
	gr.Router.Use(mw...)
}
