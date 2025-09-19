package modules

var modules = make([]func(), 0)

func NewModule(fn func()) {

	modules = append(modules, fn)
}

func Execute() {
	for _, fn := range modules {
		fn()
	}
}
