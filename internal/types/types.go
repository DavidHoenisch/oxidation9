package types

type Tool interface {
	Handler(args []string) error
	Parse(args []string) error
	Exec() error
	Doc()
}
