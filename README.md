# la -- A Lambda Calculus Interpreter With Module and Std

Support 
- `let` syntax sugar: `(\x.M)N == let x = N in M`
- Module system: `#use filename` (implemented as a simple text replacement preprocessor)
- Standard library `std`
- Pretty printing church number: `(λf.(λx.(f x))) -> 1`
- Showing each rewrite step by step: `R288: (λf.(λx.(f (f x)))) -> 2` Rewrite step 288
- Normal order evaluation: support Y combinator recursion
- Comment syntax: `// this is a comment`
- Multiple arguments syntax sugar: `\f x.f (f x)` == `\f.\x.f (f x)`

Build for native:
```sh
git clone https://github.com/guiyuanju/lamb
cd lamb
go build -o la ./cmd/native

```
Usage:
```bash
# REPL
./la

# Evaluate file
./la filename.la
```

Syntax:
- Variable: Can be number, letter, special symbol or any combination of them, but cannot start with underscore `_`, which is used inside interpreter for fresh variable generation.
- Abstraction:
    - `\x.x`
    - `\x.\y.y x` function only support one argument
    - Parenthesis is optional: `(\x.x y)` is equal to `\x.x y`
- Application:
    - `a b` is equal to `(a b)`
    - `a b c` is equal to `((a b) c)`, left associative
    - `(\x.x y) y` apply lambda to argument
- Let:
    - `let a = b in body` replace `a` with `b` in body
    - A syntax sugar for `(\a.body) b`
    - Thus cannot define recursively
    - Can be nested: `let a = b in let c = d in a c` => `b d`
- Module:
    - `#use std`, `std` has no quotes, there must be a file `std.la` in current directory
    - A module is a simple nested `let`: `let a = b in c = d in`
    - The content of a module is simply copied and replace the `#use` directive

REPL examples:
```
> (\x.x) y
R1: y

> (\x.x x) (\x.x x)
R1: ((λx.(x x)) (λx.(x x)))

> let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ 0
...
R4: (λf.(λx.(f ((λx.x) x))))
R5: (λf.(λx.(f x))) -> 1

> #use std + 1 2
...
R39: (λf.(λx.(f (f (f ((λx.x) x))))))
R40: (λf.(λx.(f (f (f x))))) -> 3
```

File examples:
```bash
./la main.la
...
R287: (λf.(λx.(f (f ((λx.x) x)))))
R288: (λf.(λx.(f (f x)))) -> 2
```

Where `main.la`:
```txt
#use std
let factorial = \r.\n.(if (zero? n) 1 (* n (r (- n 1)))) in
// euqals * 2 1
Y factorial 2
```

## Wasm

Build as Wasm to run in Browser:
```sh
GOARCH=wasm GOOS=js go build -o ./cmd/js/la.wasm ./cmd/js
# Open a server in ./cmd/js, for example, you can use simplehttpserver
simplehttpserver
```

Then open the link provided by the server in browser, the default of `simplehttpserver` is `http://0.0.0.0:8000/`.

![res](./img/fac3.jpg)

```
```
