# Lambda Calculus Interpreter

Support 
- `let` syntax sugar: `(\x.M)N == let x = N in M`
- Simple module system: `#use filename`
- Standard library `std`
- Pretty printing church number
- Showing each rewrite step by step
- Normal order evaluation, support Y combinator recursion
- Comment syntax `//`

Build:
```sh
git clone <repo> <dir>
cd <dir>
go build

```
Usage:
```bash
# REPL
./lamb

# Evaluate file
./lamb filename
```

REPL examples:
```
> #use std + 1 2
...
R39: (λf.(λx.(f (f (f ((λx.x) x))))))
R40: (λf.(λx.(f (f (f x))))) -> 3

> let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ 0
...
R4: (λf.(λx.(f ((λx.x) x))))
R5: (λf.(λx.(f x))) -> 1
```

File examples:
```bash
./lamb main.lamb
...
R287: (λf.(λx.(f (f ((λx.x) x)))))
R288: (λf.(λx.(f (f x)))) -> 2
```

Where `main.lamb`:
```txt
#use std
let factorial = \r.\n.(if (zero? n) 1 (* n (r (- n 1)))) in
// euqals * 2 1
Y factorial 2
```
