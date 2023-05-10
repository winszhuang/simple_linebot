package constant

type Directive string

const (
	Pick   Directive = "/pick"
	Add    Directive = "/add"
	List   Directive = "/list"
	Remove Directive = "/rm"
	Near   Directive = "/near"
)

func IsDirective(text string) bool {
	switch text {
	case string(Add), string(List), string(Remove), string(Pick), string(Near):
		return true
	}
	return false
}
