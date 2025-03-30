

```
declaration    → classDecl
               | funDecl
               | varDecl
               | statement ;

classDecl      → "class" IDENTIFIER "{" function* "}" ;
```

```
function       → IDENTIFIER "(" parameters? ")" block ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
assignment     → ( call "." )? IDENTIFIER "=" assignment
               | logic_or ;
```