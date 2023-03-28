package rest

type Handler[S Session] func(S, *Request) (any, error)

func (h Handler[S]) Do(t S, r *Request) (any, error) { return h(t, r) }
