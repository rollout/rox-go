package context

type contextImp struct {
	items map[string]interface{}
}

func NewContext(items map[string]interface{}) Context {
	return &contextImp{items: items}
}

func (c *contextImp) Get(key string) interface{} {
	if key == "" {
		return nil
	}

	if item, ok := c.items[key]; !ok {
		return nil
	} else {
		return item
	}
}
