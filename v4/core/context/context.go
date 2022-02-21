package context

type Context interface {
	Get(key string) interface{}
}
