/*
 * Moxie Programming Language Grammar
 * ANTLR 4 Grammar Specification
 *
 * Based on Go language specification with modifications from go-language-revision.md:
 * - Explicit pointer types for slices, maps, channels
 * - Mutable strings (string = *[]byte)
 * - Const enforcement with MMU protection
 * - Native FFI (dlopen/dlsym)
 * - Zero-copy type coercion
 * - No make() or append() functions
 * - Explicit int sizes (no platform-dependent int/uint)
 * - Concatenation operator: | (vertical bar) for strings and slices
 *   Following cryptographic notation where a | b means concatenation
 */

grammar Moxie;

// ===== Source File Structure =====

sourceFile
    : packageClause eos importDecl* topLevelDecl* EOF
    ;

packageClause
    : 'package' IDENTIFIER
    ;

importDecl
    : 'import' ( importSpec | '(' ( importSpec eos )* ')' )
    ;

importSpec
    : ( '.' | IDENTIFIER )? string_
    ;

// ===== Top Level Declarations =====

topLevelDecl
    : declaration
    | functionDecl
    | methodDecl
    ;

declaration
    : constDecl
    | typeDecl
    | varDecl
    ;

// ===== Const Declarations (Enhanced for Moxie) =====

constDecl
    : 'const' ( constSpec | '(' ( constSpec eos )* ')' )
    ;

constSpec
    : identifierList ( type_? '=' expressionList )?
    ;

// ===== Type Declarations =====

typeDecl
    : 'type' ( typeSpec | '(' ( typeSpec eos )* ')' )
    ;

typeSpec
    : IDENTIFIER typeParameters? '=' type_     # TypeAlias
    | IDENTIFIER typeParameters? type_         # TypeDef
    ;

typeParameters
    : '[' typeParameterDecl ( ',' typeParameterDecl )* ','? ']'
    ;

typeParameterDecl
    : identifierList typeConstraint
    ;

typeConstraint
    : type_
    ;

// ===== Variable Declarations =====

varDecl
    : 'var' ( varSpec | '(' ( varSpec eos )* ')' )
    ;

varSpec
    : identifierList ( type_ ( '=' expressionList )? | '=' expressionList )
    ;

// ===== Function Declarations =====

functionDecl
    : 'func' IDENTIFIER typeParameters? signature block?
    ;

methodDecl
    : 'func' receiver IDENTIFIER signature block?
    ;

receiver
    : '(' ( IDENTIFIER? type_ ) ')'
    ;

// ===== Types =====

type_
    : typeName typeArgs?              # NamedType
    | typeLit                          # TypeLiteral
    | '(' type_ ')'                    # ParenType
    | 'const' type_                    # ConstType
    ;

typeName
    : IDENTIFIER ( '.' IDENTIFIER )?
    ;

typeArgs
    : '[' typeList ','? ']'
    ;

typeLit
    : arrayType
    | structType
    | pointerType
    | functionType
    | interfaceType
    | sliceType
    | mapType
    | channelType
    ;

// Array: [N]T
arrayType
    : '[' arrayLength ']' elementType
    ;

arrayLength
    : expression
    ;

elementType
    : type_
    ;

// Slice: *[]T (explicit pointer in Moxie)
sliceType
    : '*' '[' ']' elementType
    | '[' ']' elementType  // For parsing compatibility
    ;

// Struct
structType
    : 'struct' '{' ( fieldDecl eos )* '}'
    ;

fieldDecl
    : ( identifierList type_ | embeddedField ) tag_?
    ;

embeddedField
    : '*'? typeName
    ;

tag_
    : RAW_STRING_LIT
    | INTERPRETED_STRING_LIT
    ;

// Pointer: *T
pointerType
    : '*' type_
    ;

// Function: func(params) result
functionType
    : 'func' signature
    ;

signature
    : parameters result?
    ;

result
    : parameters
    | type_
    ;

parameters
    : '(' ( parameterDecl ( ',' parameterDecl )* ','? )? ')'
    ;

parameterDecl
    : identifierList? '...'? type_
    ;

// Interface
interfaceType
    : 'interface' '{' ( interfaceElem eos )* '}'
    ;

interfaceElem
    : methodElem
    | typeElem
    ;

methodElem
    : IDENTIFIER signature
    ;

typeElem
    : typeTerm ( '|' typeTerm )*
    ;

typeTerm
    : type_
    | '~' type_
    ;

// Map: *map[K]V (explicit pointer in Moxie)
mapType
    : '*' 'map' '[' type_ ']' elementType
    | 'map' '[' type_ ']' elementType  // For parsing compatibility
    ;

// Channel: *chan T (explicit pointer in Moxie)
channelType
    : '*' 'chan' '<-'? elementType     # SendRecvChan
    | '*' '<-' 'chan' elementType      # RecvOnlyChan
    | 'chan' '<-'? elementType         # SendRecvChanCompat
    | '<-' 'chan' elementType          # RecvOnlyChanCompat
    ;

// ===== Statements =====

block
    : '{' statementList? '}'
    ;

statementList
    : ( statement eos )+
    ;

statement
    : declaration                      # DeclStmt
    | simpleStmt                       # SimpleStatement
    | labeledStmt                      # LabeledStatement
    | goStmt                          # GoStatement
    | returnStmt                       # ReturnStatement
    | breakStmt                        # BreakStatement
    | continueStmt                     # ContinueStatement
    | gotoStmt                        # GotoStatement
    | fallthroughStmt                  # FallthroughStatement
    | block                           # BlockStatement
    | ifStmt                          # IfStatement
    | switchStmt                       # SwitchStatement
    | selectStmt                       # SelectStatement
    | forStmt                         # ForStatement
    | deferStmt                       # DeferStatement
    ;

simpleStmt
    : expressionStmt
    | sendStmt
    | incDecStmt
    | assignment
    | shortVarDecl
    ;

expressionStmt
    : expression
    ;

sendStmt
    : expression '<-' expression
    ;

incDecStmt
    : expression ( '++' | '--' )
    ;

assignment
    : expressionList assign_op expressionList
    ;

assign_op
    : ( '+' | '-' | '|' | '^' | '*' | '/' | '%' | '<<' | '>>' | '&' | '&^' )? '='
    ;

shortVarDecl
    : identifierList ':=' expressionList
    ;

labeledStmt
    : IDENTIFIER ':' statement
    ;

returnStmt
    : 'return' expressionList?
    ;

breakStmt
    : 'break' IDENTIFIER?
    ;

continueStmt
    : 'continue' IDENTIFIER?
    ;

gotoStmt
    : 'goto' IDENTIFIER
    ;

fallthroughStmt
    : 'fallthrough'
    ;

deferStmt
    : 'defer' expression
    ;

ifStmt
    : 'if' ( simpleStmt ';' )? expression block ( 'else' ( ifStmt | block ) )?
    ;

switchStmt
    : exprSwitchStmt
    | typeSwitchStmt
    ;

exprSwitchStmt
    : 'switch' ( simpleStmt ';' )? expression? '{' exprCaseClause* '}'
    ;

exprCaseClause
    : exprSwitchCase ':' statementList?
    ;

exprSwitchCase
    : 'case' expressionList
    | 'default'
    ;

typeSwitchStmt
    : 'switch' ( simpleStmt ';' )? typeSwitchGuard '{' typeCaseClause* '}'
    ;

typeSwitchGuard
    : ( IDENTIFIER ':=' )? primaryExpr '.' '(' 'type' ')'
    ;

typeCaseClause
    : typeSwitchCase ':' statementList?
    ;

typeSwitchCase
    : 'case' typeList
    | 'default'
    ;

typeList
    : type_ ( ',' type_ )*
    ;

selectStmt
    : 'select' '{' commClause* '}'
    ;

commClause
    : commCase ':' statementList?
    ;

commCase
    : 'case' ( sendStmt | recvStmt )
    | 'default'
    ;

recvStmt
    : ( expressionList '=' | identifierList ':=' )? expression
    ;

forStmt
    : 'for' ( expression | forClause | rangeClause )? block
    ;

forClause
    : simpleStmt? ';' expression? ';' simpleStmt?
    ;

rangeClause
    : ( expressionList '=' | identifierList ':=' )? 'range' expression
    ;

goStmt
    : 'go' expression
    ;

// ===== Expressions =====

expression
    : unaryExpr                                           # UnaryExpression
    | expression mul_op expression                        # MultiplicativeExpr
    | expression add_op expression                        # AdditiveExpr
    | expression '|' expression                           # ConcatenationExpr
    | expression rel_op expression                        # RelationalExpr
    | expression '&&' expression                          # LogicalAndExpr
    | expression '||' expression                          # LogicalOrExpr
    ;

primaryExpr
    : operand                                             # PrimaryOperand
    | conversion                                          # ConversionExpr
    | methodExpr                                          # MethodExpression
    | primaryExpr selector                                # SelectorExpr
    | primaryExpr index                                   # IndexExpr
    | primaryExpr slice_                                  # SliceExpr
    | primaryExpr typeAssertion                           # TypeAssertionExpr
    | primaryExpr arguments                               # CallExpr
    ;

unaryExpr
    : primaryExpr
    | unary_op unaryExpr
    ;

conversion
    : type_ '(' expression ','? ')'                       # SimpleConversion
    | '(' '*' '[' ']' type_ ')' '(' expression ')'       # SliceCastExpr
    | '(' '*' '[' ']' type_ ',' endianness ')' '(' expression ')'  # SliceCastEndianExpr
    | '&' '(' '*' '[' ']' type_ ')' '(' expression ')'   # SliceCastCopyExpr
    | '&' '(' '*' '[' ']' type_ ',' endianness ')' '(' expression ')' # SliceCastCopyEndianExpr
    ;

endianness
    : 'NativeEndian'
    | 'LittleEndian'
    | 'BigEndian'
    ;

operand
    : literal                                             # LiteralOperand
    | operandName                                         # NameOperand
    | '(' expression ')'                                  # ParenOperand
    ;

literal
    : basicLit
    | compositeLit
    | functionLit
    ;

basicLit
    : INT_LIT
    | FLOAT_LIT
    | IMAGINARY_LIT
    | RUNE_LIT
    | string_
    ;

string_
    : RAW_STRING_LIT
    | INTERPRETED_STRING_LIT
    ;

operandName
    : IDENTIFIER
    | qualifiedIdent
    ;

qualifiedIdent
    : IDENTIFIER '.' IDENTIFIER
    ;

// Composite Literals - Enhanced for Moxie
compositeLit
    : literalType literalValue
    ;

literalType
    : structType
    | arrayType
    | sliceType      // &[]T{...}
    | mapType        // &map[K]V{...}
    | channelType    // &chan T{cap: N}
    | typeName typeArgs?
    | '[' '...' ']' elementType
    ;

literalValue
    : '{' ( elementList ','? )? '}'
    ;

elementList
    : keyedElement ( ',' keyedElement )*
    ;

keyedElement
    : ( key ':' )? element
    ;

key
    : IDENTIFIER
    | expression
    | literalValue
    ;

element
    : expression
    | literalValue
    ;

functionLit
    : 'func' signature block
    ;

selector
    : '.' IDENTIFIER
    ;

index
    : '[' expression ']'
    ;

slice_
    : '[' ( expression? ':' expression? ( ':' expression )? | expression? ':' )  ']'
    ;

typeAssertion
    : '.' '(' type_ ')'
    ;

arguments
    : '(' ( ( expressionList | type_ ( ',' expressionList )? ) '...'? ','? )? ')'
    ;

methodExpr
    : type_ '.' IDENTIFIER
    ;

mul_op
    : '*' | '/' | '%' | '<<' | '>>' | '&' | '&^'
    ;

add_op
    : '+' | '-' | '^'
    ;

rel_op
    : '==' | '!=' | '<' | '<=' | '>' | '>='
    ;

unary_op
    : '+' | '-' | '!' | '^' | '*' | '&' | '<-'
    ;

// ===== Lists =====

expressionList
    : expression ( ',' expression )*
    ;

identifierList
    : IDENTIFIER ( ',' IDENTIFIER )*
    ;

// ===== End of Statement =====

eos
    : ';'
    | EOF
    | TERMINATOR
    ;

// ===== Lexer Rules =====

// Keywords
BREAK       : 'break';
CASE        : 'case';
CHAN        : 'chan';
CONST       : 'const';
CONTINUE    : 'continue';
DEFAULT     : 'default';
DEFER       : 'defer';
ELSE        : 'else';
FALLTHROUGH : 'fallthrough';
FOR         : 'for';
FUNC        : 'func';
GO          : 'go';
GOTO        : 'goto';
IF          : 'if';
IMPORT      : 'import';
INTERFACE   : 'interface';
MAP         : 'map';
PACKAGE     : 'package';
RANGE       : 'range';
RETURN      : 'return';
SELECT      : 'select';
STRUCT      : 'struct';
SWITCH      : 'switch';
TYPE        : 'type';
VAR         : 'var';

// Moxie-specific built-ins (not keywords, but recognized identifiers)
// These are parsed as IDENTIFIER but semantically checked:
// clone, copy, grow, clear, free, dlopen, dlsym, dlclose, dlerror, dlopen_mem
// NOTE: append() is NOT a built-in - use | operator for concatenation

// Types - Note: No 'int' or 'uint' in Moxie (explicit sizes required)
BOOL        : 'bool';
BYTE        : 'byte';
INT8        : 'int8';
INT16       : 'int16';
INT32       : 'int32';
INT64       : 'int64';
UINT8       : 'uint8';
UINT16      : 'uint16';
UINT32      : 'uint32';
UINT64      : 'uint64';
FLOAT32     : 'float32';
FLOAT64     : 'float64';
COMPLEX64   : 'complex64';
COMPLEX128  : 'complex128';
STRING      : 'string';
UINTPTR     : 'uintptr';
RUNE        : 'rune';

// Predeclared identifiers (zero values, etc.)
NIL         : 'nil';
TRUE        : 'true';
FALSE       : 'false';
IOTA        : 'iota';

// Literals

DECIMAL_LIT
    : '0'
    | [1-9] ('_'? [0-9])*
    ;

BINARY_LIT
    : '0' [bB] '_'? [01] ('_'? [01])*
    ;

OCTAL_LIT
    : '0' [oO]? '_'? [0-7] ('_'? [0-7])*
    ;

HEX_LIT
    : '0' [xX] '_'? HEX_DIGIT ('_'? HEX_DIGIT)*
    ;

fragment HEX_DIGIT
    : [0-9a-fA-F]
    ;

FLOAT_LIT
    : DECIMAL_FLOAT_LIT
    | HEX_FLOAT_LIT
    ;

fragment DECIMAL_FLOAT_LIT
    : DECIMALS '.' DECIMALS? EXPONENT?
    | DECIMALS EXPONENT
    | '.' DECIMALS EXPONENT?
    ;

fragment DECIMALS
    : [0-9] ('_'? [0-9])*
    ;

fragment EXPONENT
    : [eE] [+-]? DECIMALS
    ;

fragment HEX_FLOAT_LIT
    : '0' [xX] HEX_MANTISSA HEX_EXPONENT
    ;

fragment HEX_MANTISSA
    : '_'? HEX_DIGIT ('_'? HEX_DIGIT)* ( '.' ('_'? HEX_DIGIT)* )?
    | '_'? '.' HEX_DIGIT ('_'? HEX_DIGIT)*
    ;

fragment HEX_EXPONENT
    : [pP] [+-]? DECIMALS
    ;

IMAGINARY_LIT
    : ( DECIMALS | INT_LIT | FLOAT_LIT ) 'i'
    ;

INT_LIT
    : DECIMAL_LIT
    | BINARY_LIT
    | OCTAL_LIT
    | HEX_LIT
    ;

RUNE_LIT
    : '\'' ( ~['\\\r\n] | ESCAPE_SEQ | UNICODE_VALUE ) '\''
    ;

fragment UNICODE_VALUE
    : '\\u' HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT
    | '\\U' HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT
    ;

fragment ESCAPE_SEQ
    : '\\' [abfnrtv\\'"]
    | '\\' OCTAL_DIGIT OCTAL_DIGIT OCTAL_DIGIT
    | '\\x' HEX_DIGIT HEX_DIGIT
    ;

fragment OCTAL_DIGIT
    : [0-7]
    ;

RAW_STRING_LIT
    : '`' ~'`'* '`'
    ;

INTERPRETED_STRING_LIT
    : '"' ( ~["\\\r\n] | ESCAPE_SEQ | UNICODE_VALUE )* '"'
    ;

// Identifiers

IDENTIFIER
    : LETTER ( LETTER | UNICODE_DIGIT )*
    ;

fragment LETTER
    : UNICODE_LETTER
    | '_'
    ;

fragment UNICODE_LETTER
    : [\p{L}]
    ;

fragment UNICODE_DIGIT
    : [\p{Nd}]
    ;

// Operators and delimiters (these are handled by parser rules mostly)

// Comments

LINE_COMMENT
    : '//' ~[\r\n]* -> skip
    ;

BLOCK_COMMENT
    : '/*' .*? '*/' -> skip
    ;

// Whitespace

WS
    : [ \t\r\n]+ -> skip
    ;

TERMINATOR
    : [\r\n]+
    ;

// Catch-all for any other character (error)
OTHER
    : .
    ;
