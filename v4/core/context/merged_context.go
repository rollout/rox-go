package context

type MergedContext struct {
	globalContext Context
	localContext  Context
}

func NewMergedContext(globalContext, localContext Context) Context {
	return &MergedContext{globalContext: globalContext, localContext: localContext}
}

func (mc *MergedContext) Get(key string) interface{} {
	if mc.localContext != nil {
		if item := mc.localContext.Get(key); item != nil {
			return item
		}
	}

	if mc.globalContext != nil {
		return mc.globalContext.Get(key)
	}

	return nil
}
