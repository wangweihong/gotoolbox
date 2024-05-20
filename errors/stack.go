package errors

type StackTrace interface {
	Stack() []string
}
