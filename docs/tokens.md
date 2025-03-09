#  Tokenisation in Monogram

In common with most programming languages, Mongram separates analysing the input
into two levels: tokenisation and then parsing the token stream. On this page 
we describe the different kinds of tokens that Monogram recognises and try
to explain the thinking behind some of the less familiar-looking tokens.

## Comments

Monogram supports `#` as the end of line comment character in the same way
as Python, Ruby and Bash. These are simply discarded by a Monogram processor.


## The different types of tokens

- Numbers: positive and negative integers and floats. 

- Strings: single, double and back quoted strings are all supported. These all 
  support string interpolation and uniformly combine with raw and multiline 
  variants.

- Symbols: these include the single character delimiters `(`, `)`, `[`, `]`, 
  `{`, `}`and also punctuation such as `,` and `;`. These
  do not 'stick' to any other character, although may appear inside strings. 
  Delimiters play the dual roles of nesting expressions and
  supporting function and method calls. 
  e.g. `(x + y) * z` vs `f[x, y]` and `x.m(p, q)`.

- Signs: these are runs of sign-characters such as `+`, `**`, `-->` and so
  on. These primarily play the role of infix operators e.g. `x := y`. In
  some places they act as supporting 'punctuation'.use

- Identifiers: the usual rule for identifiers is followed - they start with
  an alphabetical character or an underscore and continue with those, plus digits.
  However Mongram also supports the use of underscore as an identifier 
  quoting mechanism.

## Numbers in detail

Currently these are in the same format as JSONs numbers. However we intend to
extend them to include different radixes and underbars for improved readability.

## Strings in detail

_WORK IN PROGRESS, March 2025_

## Symbols in detail

_WORK IN PROGRESS, March 2025_

## Signs in detail

_WORK IN PROGRESS, March 2025_


## Identifiers in detail

There are two classes of identifiers: basic identifiers and underscore quoted
identifiers. Basic identifiers start with a letter or underscore and continue
with letters, digits and underscore. However they may not start and finish with
an underscore. Basic identifiers that do not start or end with an underscore are
called core identifiers.

Underscore quoted identifiers use underscores as delimiting quote characters
and, within that scope, requires escape sequences for any character other than
a letter or digit. Their purpose is to allow arbitrary characters to be used as
an identifier.

The escape sequences they support exclude interpolation. Hence `\(` is a 
plain open parenthesis. 

## Token rules as a railroad diagram

![Railroad diagram](images/token_rules.png)