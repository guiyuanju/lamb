package lamb

import "fmt"

// Capture-avoiding substitution M[x := N]
// - x[x := N] = N; y[x := N] = y if y != x
// - (M1 M2)[x := N] = (M1[x := N])(M2[x := N])
// - (λy.M)[x := N]
//   - If y == x, then λx.M (no change)
//   - If y ∉ FV(N), then λy.(M[x := N]) (FV(N): free variables of N)
//   - If y ∈ FV(N), then λy'.M[y := y'], where y' is fresh, then continue as above
func substitute(m Term, x Variable, n Term) Term {
	switch m := m.(type) {
	case Variable:
		if m != x {
			return m
		}
		return n
	case Application:
		return Application{
			Left:  substitute(m.Left, x, n),
			Right: substitute(m.Right, x, n),
		}
	case Abstraction:
		if m.Param == x {
			return m
		}
		if !isFreeVariable(m.Param, n) {
			return Abstraction{
				Param: m.Param,
				Body:  substitute(m.Body, x, n),
			}
		}
		freshVar := getFreshVariable()
		newAbs := Abstraction{
			Param: freshVar,
			Body:  substitute(m.Body, m.Param, freshVar),
		}
		return substitute(newAbs, x, n)
	default:
		panic("unrecognized term")
	}
}

var freshVarCounter int = 0

func getFreshVariable() Variable {
	res := fmt.Sprintf("_%d", freshVarCounter)
	freshVarCounter++
	return Variable(res)
}

func isFreeVariable(x Variable, t Term) bool {
	switch t := t.(type) {
	case Variable:
		return x == t
	case Application:
		return isFreeVariable(x, t.Left) || isFreeVariable(x, t.Right)
	case Abstraction:
		if t.Param == x {
			return false
		}
		return isFreeVariable(x, t.Body)
	default:
		panic("unrecognized term")
	}
}

// Reduce rewrite one piece of term at a time.
// Need an outer loop to proceed steps,
// and termination condition should be new term returned is the same as the input term.
// The reason for this approach is because simple recursion won't terminate
// when evaluate e.g. Y combinator, step-wise reduce enables comparison
// of previous and current terms, and stop if no changes (not reducible).
func Reduce(t Term) Term {
	switch t := t.(type) {
	case Variable:
		return t
	case Application:
		switch lhs := t.Left.(type) {
		case Variable:
			return Application{lhs, Reduce(t.Right)}
		case Application:
			return Application{Reduce(lhs), t.Right}
		case Abstraction:
			return substitute(lhs.Body, lhs.Param, t.Right)
		default:
			panic("unrecognized term")
		}
	case Abstraction:
		return Abstraction{
			Param: t.Param,
			Body:  Reduce(t.Body),
		}
	default:
		panic("unrecognized term")
	}
}
