# Lambda Calculus Interpreter

Support 
- `let` syntax sugar: `(\x.M)N == let x = N in M`
- Module system: `#use filename` (implemented as a simple text replacement preprocessor)
- Standard library `std`
- Pretty printing church number: `(λf.(λx.(f x))) -> 1`
- Showing each rewrite step by step: `R288: (λf.(λx.(f (f x)))) -> 2` Rewrite step 282
- Normal order evaluation: support Y combinator recursion
- Comment syntax: `// this is a comment`

Build:
```sh
git clone https://github.com/guiyuanju/lamb
cd lamb
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
