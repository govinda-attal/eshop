package platform

type CtxKey string

const (
	TraceID CtxKey = "trace-id"
)

func (k CtxKey) String() string {
	return string(k)
}
