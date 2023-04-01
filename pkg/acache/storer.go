package acache

type RouteStorer interface {
	GetRoute(RouteFilter) (Route, error)
	GetRoutes() (Routes, error)
	AddRoute(Route) error
	Clear() error
}
