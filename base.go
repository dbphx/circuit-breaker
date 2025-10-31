package ecocircuitbreaker

type IBreaker[T any] interface {
	Execute(fn func() (T, error)) (T, error)
	CurrentState() State
}
