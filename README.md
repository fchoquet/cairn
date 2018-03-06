# The cairn programming language

Experimentations and exploration around programming languages and compilators.

# Context Free Grammar

## Arithmetic expressions

```
expr   : term ((PLUS | MINUS) term)*
term   : factor ((MUL | DIV) factor)*
factor : INTEGER | LPAREN expr RPAREN
```
