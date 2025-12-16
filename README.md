# Lambda Calculus Interpreter

Support `let` syntax sugar: `(\x.M)N == let x = N in M`

Usage:

```
> let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ 0
(λf.(λx.(f x)))

> let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ (succ 0)
(λf.(λx.(f (f x))))

> let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ (succ (succ 0))
(λf.(λx.(f (f (f x)))))

> let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in let 1 = succ 0 in let + = \m.\n.\f.\x .m f (n f x) in (+ 1 (+ 1 (+ 1 1)))
(λf.(λx.(f (f (f (f x))))))
```

Run REPL:

```sh
git clone <repo> <dir>
cd <dir>
go run .
```
