package sqlcom

//Where builder
type Where struct {
	Sqls []string
	Args []interface{}
}

//NewWhere new where
func NewWhere(opts ...func(w *Where)) *Where {
	w := &Where{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}
