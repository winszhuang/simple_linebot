package constants

type Directive string
type Question string

const (
	Pick   Directive = "/pick"
	Add              = "/add"
	List             = "/list"
	Remove           = "/rm"
	Near             = "/near"
)

const (
	Phone Question = "Phone"
	Name           = "Name"
)

func IsDirective(text string) bool {
	switch text {
	case Add, List, Remove, string(Pick), Near:
		return true
	}
	return false
}

func IsQuestion(text string) bool {
	return text == string(Phone) || text == string(Name)
}
