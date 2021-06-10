package domain

type ErrorLoggerInterface interface {
	Log(err error)
}
