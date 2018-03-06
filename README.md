# The cairn programming language

Experimentations and exploration around programming languages and compilators.

# Context Free Grammar

## Expression

```
expr : arithmexpr | strexpr
```

## Arithmetic expressions

```
arithmexpr : term ((PLUS | MINUS) term)*
term       : factor ((MUL | DIV) factor)*
factor     : (PLUS|MINUS)factor | INTEGER | LPAREN arithmexpr RPAREN
```

## Strings

```
strexpr : str (CONCAT str)*
str     : STRING
```
