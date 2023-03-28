package rest

type Session interface {
	IsExpired() bool
}

type NewSession[S Session] func(*Request) (S, error)
