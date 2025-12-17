package lamb

import "fmt"

// Definition:
//
//	Variable
//	Abstraction
//	application
type (
	Variable    string
	Abstraction struct {
		Param Variable
		Body  Term
	}
)

type Application struct {
	Left  Term
	Right Term
}
type Term interface {
	isTerm()
}

func (v Variable) isTerm()    {}
func (a Abstraction) isTerm() {}
func (a Application) isTerm() {}

func (a Abstraction) String() string {
	return fmt.Sprintf("(λ%s.%s)", a.Param, a.Body)
}

func (a Application) String() string {
	return fmt.Sprintf("(%s %s)", a.Left, a.Right)
}
