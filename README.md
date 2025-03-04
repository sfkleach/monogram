# Monogram

Monogram is a "no batteries" notation for writing domain-specific programs and
configuration files. It is easy for humans to read and write. It is easy for
machines to parse and generate. It deliberately borrows from many programming
languages but feels familiar to Python and Ruby programmers.

## _"It's source code, Jim. But not as we know it!"_

Here's an initial example to help explain what we mean by 'batteries not included'.
To experienced programmers, the following code looks a lot like the definition
of the factorial function:
```py
def f(n):
    if n <= 1:
        1
    else:
        n * f(n - 1)
    endif
enddef
```

However, the twist is that Monogram has no idea what `def` or `if` might mean!
Nor does it have a clue about `*` or `-` either. And it definitely cannot
execute this program. 

And yet Monogram can easily translate this example into neatly structured XML
(shown below). Or it can translate to [JSON](docs/json.md) or [YAML](docs/yaml.md).
```xml
<form>
    <part keyword="def">
        <apply kind="parentheses" separator="undefined">
            <identifier name="f"/>
            <arguments>
                <identifier name="n"/>
            </arguments>
        </apply>
    </part>
    <part keyword="_">
        <form>
            <part keyword="if">
                <operator name="&lt;=">
                    <identifier name="n"/>
                    <number value="1"/>
                </operator>
            </part>
            <part keyword="_">
                <number value="1"/>
            </part>
            <part keyword="else">
                <operator name="*">
                    <identifier name="n"/>
                    <apply kind="parentheses" separator="undefined">
                        <identifier name="f"/>
                        <arguments>
                            <operator name="-">
                                <identifier name="n"/>
                                <number value="1"/>
                            </operator>
                        </arguments>
                    </apply>
                </operator>
            </part>
        </form>
    </part>
</form>
```

Alternatively it can render the code as a diagram using Mermaid (below) or 
[Graphviz](docs/dot.md). Here's the same structure visualised as a graph.

```mermaid
graph TD
  126587637945888["form"]:::custom_form;
  126587637945968["part: def"]:::custom_part;
  126587637945888 --> 126587637945968;
  126587637946048["apply"]:::custom_apply;
  126587637945968 --> 126587637946048;
  126587637946128["identifier: f"]:::custom_identifier;
  126587637946048 --> 126587637946128;
  126587637946208["arguments"]:::custom_arguments;
  126587637946048 --> 126587637946208;
  126587637946288["identifier: n"]:::custom_identifier;
  126587637946208 --> 126587637946288;
  126587637946368["part: _"]:::custom_part;
  126587637945888 --> 126587637946368;
  126587637946448["form"]:::custom_form;
  126587637946368 --> 126587637946448;
  126587637946528["part: if"]:::custom_part;
  126587637946448 --> 126587637946528;
  126587637946608["operator: <="]:::custom_operator;
  126587637946528 --> 126587637946608;
  126587637946688["identifier: n"]:::custom_identifier;
  126587637946608 --> 126587637946688;
  126587637946768["number: 1"]:::custom_number;
  126587637946608 --> 126587637946768;
  126587637946848["part: _"]:::custom_part;
  126587637946448 --> 126587637946848;
  126587637946928["number: 1"]:::custom_number;
  126587637946848 --> 126587637946928;
  126587637947008["part: else"]:::custom_part;
  126587637946448 --> 126587637947008;
  126587637947088["operator: *"]:::custom_operator;
  126587637947008 --> 126587637947088;
  126587637947168["identifier: n"]:::custom_identifier;
  126587637947088 --> 126587637947168;
  126587637947248["apply"]:::custom_apply;
  126587637947088 --> 126587637947248;
  126587637947408["identifier: f"]:::custom_identifier;
  126587637947248 --> 126587637947408;
  126587637947568["arguments"]:::custom_arguments;
  126587637947248 --> 126587637947568;
  126587637947728["operator: -"]:::custom_operator;
  126587637947568 --> 126587637947728;
  126587637947888["identifier: n"]:::custom_identifier;
  126587637947728 --> 126587637947888;
  126587637948048["number: 1"]:::custom_number;
  126587637947728 --> 126587637948048;

classDef custom_form fill:lightpink,stroke:#000,stroke-width:2px;
classDef custom_part fill:#FFD8E1,stroke:#000,stroke-width:2px;
classDef custom_apply fill:lightgreen,stroke:#000,stroke-width:2px;
classDef custom_identifier fill:Honeydew,stroke:#000,stroke-width:2px;
classDef custom_arguments fill:PaleTurquoise,stroke:#000,stroke-width:2px;
classDef custom_operator fill:#C0FFC0,stroke:#000,stroke-width:2px;
classDef custom_number fill:lightgoldenrodyellow,stroke:#000,stroke-width:2px;
```

In other words, Monogram is just a notation for writing program-like "code" but
comes without any built-in meanings. Although it is not infinitely flexible, it 
can often save you the effort of designing the syntax and implementing a parser
when you want an application/domain-specific language.

For more examples and more output formats (like JSON, YAML, PNG) see the 
[examples page](docs/examples.md).


## Monogram grammar

### Overview of tokens

The basic building blocks of a Monogram document are tokens - that is to say
numbers (`123`, `-0.12`), strings (`"hello, world"`), symbols (`{`, `}`), signs
(`:`, `++`) and various kinds of identifiers (`true`, `x`, `while`). These will
be largely familiar to anyone used to working with JSON or any mainstream
programming language.

Full details of tokenisation are given on [this page](docs/tokens.md) but
because these are generally so familiar to most programmers we highlight just a
few aspects that will be less familiar here:

- **Strings** support all three quote characters: single , double and back quotes.
    - All three are completely symmetrical in their design.
    - And support escape sequences, string interpolation, and raw and multiline
      versions.

- **Symbols** include parentheses, brackets and braces as well as punctuation such
  as `,` and `;` (but not `.`)
    - The three different brackets are treated symmetrically
    - So these are all valid expressions, for instance: `m.f(x)`,  `m.f[x]`, `m.f{x}`.

- **Operators** are runs of sign-characters. In addition to familiar single-character
  operators such as `+`, `*`, `^`, Monogram allows for arbitrary combinations
  such as `:=`, `-->` or even `++^=!$$`. 
    - These primarily play the role of infix operators.
    - Operator precedence is decided on the first character of the sign and follows
      the precedence rules of the C-programming language. As a consequence,
      we can use sequences such as `s = x + y * z` and get expected results.
    - N.B. If the first character is repeated then the precedence is slightly 
      adjusted so it binds slightly more tightly. Which is why `p = a == b`
      binds the expected way. 

- **Identifiers** 
  - Support string-like quotes using underscores e.g. `_hello, world_`
  is a valid identifier. 
  - Identifiers starting `end` are key to the way the grammar works as they
    mark reserved words.

### Overview of the grammar

_WORK IN PROGRESS, March 2025_

### Railroad diagrams

Here's the grammar for Monogram as a railroad diagram; also available in
[HTML](docs/grammar.html), [PDF](docs/images/grammar.pdf) and
[PNG](docs/images/grammar.png).

![Monogram Grammar PDF](docs/images/grammar.png) 
