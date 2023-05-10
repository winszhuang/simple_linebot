package constant

type Question string

const (
	Phone Question = "Phone"
	Name  Question = "Name"
)

func IsQuestion(text Question) bool {
	switch text {
	case Phone, Name:
		return true
	}
	return false
}
