# Context Free Grammar


## Expressions

```
expr : assignment | arithmexpr | strexpr | boolexpr

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

## Boolean expressions

```
boolexpr : boolterm ((EQ | NEQ ) boolterm)*
boolterm: (NOT)boolterm | bool | LPAREN boolexpr RPAREN
bool: BOOL | IDENTIFIER
```
