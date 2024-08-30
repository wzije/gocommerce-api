package router

type AppRouterInterface interface {
	GuestRouter()
	AuthRouter()
}
