package appctx

type Option func(*appContext)

func WithPrefix(prefix string) Option {
	return func(ac *appContext) {
		ac.prefix = prefix
	}
}

func WithName(name string) Option {
	return func(ac *appContext) {
		ac.name = name
	}
}

func WithComponent(c Component) Option {
	return func(ac *appContext) {
		if _, ok := ac.store[c.ID()]; !ok {
			ac.components = append(ac.components, c)
			ac.store[c.ID()] = c
		}
	}
}
