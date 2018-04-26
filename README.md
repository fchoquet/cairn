# The cairn programming language

Experimentations and exploration around programming languages and compilators.

# Context Free Grammar

## Expression

```
expr : assignment | arithmexpr | strexpr
```

## Assignment

```
assignment : IDENTIFIER ASSIGN expr
```

## Arithmetic expressions

```
arithmexpr : term ((PLUS | MINUS) term)*
term       : factor ((MUL | DIV) factor)*
factor     : (PLUS|MINUS)factor | INTEGER | IDENTIFIER | LPAREN arithmexpr RPAREN 
```

## Strings

```
strexpr : str (CONCAT str)*
str     : STRING | IDENTIFIER
```
