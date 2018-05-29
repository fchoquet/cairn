# Grammar

Grammar uses antlr v4 format

Everything is an expression but this grammar ignores it. The reason is that we can't use any expression with binary operators.
For instance `1 + (a := 2)` is not allowed even if `a := 2` is an expression that returns 2
I use the word `statement` for these types of expressions that can't be combined with other expressions

```
eos
    : EOF
    | EOL
    ;

statementList
    : ( statement eos )*
    ;

statement
    : simpleStmt
	;

simpleStmt
    : expression
    | assignment
    ;

expression
    : unaryExpr
    | expression BINARY_OP expression
    ;

unaryExpr
    : primaryExpr
    | UNARY_OP unaryExpr
    ;

primaryExpr
    : operand
    ;    

operand
    : literal
    | operandName
    | LPAREN expression RPAREN
    ;

operandName
    : IDENTIFIER
    ;

literal
    : basicLit
    ;

basicLit
    : INTEGER
    | STRING
    | BOOL
    ;

// assignments are expressions
assignment
    : IDENTIFIER ASSIGN expression
    ;    

block
    : BEGIN statementList END

//////////////
// functions
//////////////
functionDecl
    : FUNC IDENTIFIER ( function | signature )
    ;

function
    : signature block
    ;

signature
    : parameters result
    ;

result
    : type
    ;

parameters
    : LPAREN parameterList? RPAREN
    ;

parameterList
    : parameterDecl ( COMMA parameterDecl )*
    ;

parameterDecl
    : IDENTIFIER type
    ;

type
    : typeName
    | typeLit // not implemented yet
    ;

typeName
    : TYPE IDENTIFIER
    | qualifiedIdent // not implemented yet
    ;

```

## Binary operator precedence and associativity

Operator precedence is managed using the [precedence climbing](https://eli.thegreenplace.net/2012/08/02/parsing-expressions-by-precedence-climbing) algorithm

Implemented so far:

```
'||' | '&&' | '==' | '!=' | '+' | '-' | '^' | '*' | '/'
```

Full list:

```
'||' | '&&' | '==' | '!=' | '<' | '<=' | '>' | '>=' | '+' | '-' | '|' | '^' | '*' | '/' | '%' | '<<' | '>>' | '&' | '&^'
```
