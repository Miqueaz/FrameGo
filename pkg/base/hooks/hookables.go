package hooks

// HookFunc define el tipo de función para los hooks
type HookFunc func(data map[string]interface{}) error

// Hookable maneja la ejecución de hooks
type Hookable struct {
	BeforeRead  []HookFunc
	AfterRead   []HookFunc
	BeforeWrite []HookFunc
	AfterWrite  []HookFunc
}

// Ejecutar hooks de la lista proporcionada
func (h *Hookable) ExecuteHooks(hooks []HookFunc, data map[string]interface{}) error {
	for _, hook := range hooks {
		if err := hook(data); err != nil {
			return err
		}
	}
	return nil
}
