# The cairn programming language

Experimentations and exploration around programming languages and compilators.

# Arithmetic

```
12
> 12

12 + 34
> 46

(12 + 34) * 56 / 16
> 161
```

## Strings

```
"hello"
> hello

"hello \"Fred\""
> hello "Fred"

"hello" ++ "world"
> hello world
```

## Booleans

```
true
> true

false
> false

true == false
false

!true != !false
true

(true == false) != false
false

true && true
true

true || !false
true

```


## Variables

```
foo := "hello"
> hello

bar := " world"
> world

foo ++ bar
> hello world
```

## functions

```
func add(a:int, b:int) :int
    a + b
```
