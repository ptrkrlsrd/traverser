package acache

type RouteStorer interface {
	GetRoutes() (routes Routes, err error)
	AddRoute(Route) error
}
