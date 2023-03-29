package constants

type Directive string
type Question string

const (
	Pick   Directive = "/pick"
	Add              = "/add"
	List             = "/list"
	Remove           = "/rm"
)

const (
	Phone Question = "Phone"
	Name           = "Name"
)

func IsDirective(text string) bool {
	return text == string(Pick) || text == Add || text == List || text == Remove
}

func IsQuestion(text string) bool {
	return text == string(Phone) || text == string(Name)
}
