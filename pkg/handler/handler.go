package handler

type Middleware func(Interface) Interface

// Chain wrap middleware in order execute
func Chain(middlewares ...Middleware) func(Interface) Interface {
	return func(h Interface) Interface {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}

		return h
	}
}

func WithMiddlewares(handler Interface, middlewares ...Middleware) Interface {
	return Chain(middlewares...)(handler)
}
