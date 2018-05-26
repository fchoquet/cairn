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

# TODO: Fix this case
(true == false) != false
false

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
