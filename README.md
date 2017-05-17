[![Build Status](https://travis-ci.org/boennemann/badges.svg?branch=master)](https://travis-ci.org/boennemann/badges)
[![GitHub version](https://badge.fury.io/gh/boennemann%2Fbadges.svg)](http://badge.fury.io/gh/boennemann%2Fbadges)
[![CodeCoverage](https://scrutinizer-ci.com/g/boennemann/badges/badges/coverage.png?s=909c9b9364a927cc44392eda274de31a30b9360b)](https://scrutinizer-ci.com/g/boennemann/badges/)

# Element Field Parser

The element-field parser is a method of generating configuration file validators with a clean and consistent syntax.

## Basic Concepts

## Fields

Fields are used to assign particular pieces of data to a key.

A basic field:

```go
key = 5
key = "hi"
key = some_text
```

Fields can also store arrays:

```go
key = [5, 4, 3, 2]
key = ["hi", "me", "not"]
```

### Elements

Elements contain fields and other elements, and are used to express hierachies and tie fields together.

A basic element:

```go
key {

}
```

An element with parameters:

```go
key("25", 25, "25"){

}
```

Of course, the use of parameters provides no significant practical benefit - it is merely a different stylistic choice, and can help to emphasise particular fields over others.

# Process

First, a prototype tree is generated. This is the format which is specified in your ```.efp``` file. All files will then be compiled against this prototype tree for validation. An example prototype tree:

```go
fireVM {
    name : string!
}
```

This is then enforced in our files - all valid files must contain a top level ```fireVM``` element, and a ```name``` field, which takes a string value. In a prototype file, the ```!``` denotes a compulsory field. This is actually shorthand for a declaration like this:

```go
fireVM {
    <1:name:1> : string
}
```


Types may be in one or more of the following formats:

- Standard Aliases
- Regex
- Array Notation

## Standard Aliases

There are several default types in the ```efp``` spec:

| Alias     | Regex         |
| :-------------: |:-------------:|
| id | [a-zA-Z_]+ |
| string | "\"[^()]\"" |
| float | [0-9]*.[0-9]+    |
| bool | true|false    |
| int | [0-9]+    |

## Regex

The element-field parser supports golang regex.

```go
key = "[5-8]{4}"
```

## Array Notation

Fields can also have array values:

```go
key = [int] // any number of integers
key = [2:string] // at least two strings
key = ["###[0-7][0-7][0-7]###":4] // at most 4 regex matching sequences
key = [3:bool:3] // precisely 3 boolean values
```

Arrays are not bound to one regex sequence, and the following are legal and enforceable declarations:

```go
key = [int|string]                      // array of ints or strings
harder = [string]|[int]                 // int array or string array
twod = [[string]]                       // two-dimensional string array
harder_mixed = [string|[string]]        // array of strings or 2d array of strings
limits = [string|[3:string:5]]          // array of strings or 2d array of strings (2nd bounded by 3,5)
complex = [string|[3:string:5]|[3:[int]:3]|["x"|"[a-zA-Z]+"|[[bool]]]
```

Possible examples matching the ```complex``` field:

```go
complex = ["hi"]
complex = [["one", "two", "three"], ["one", "two", "three", "four"]]
complex = [[1],[1,2,3],[3,2,1]]
complex = [["x", "x"], ["hello this is", "dog"], [[true, true, true], [false]]]
```

These declarations are suddenly getting very complicated, surely there's some way to make them more concise?

## Aliasing

It is possible to declare aliases within a efp file, with the normal scope boundaries. Aliases are tantamount to C macros, in that they simply perform a text replace. If the text contains an element, that element will be evaluated and added to the tree.

```go
// text alias
alias name = "ender"
// field alias
alias x = num : int
// element alias
alias divs = divisions("name") {
    x
}
```

To use aliases:

```go
alias name = "ender"
alias x = num : int
alias divs = divisions("name") {
    x
}


base {
    name = name
    divs
}
```

To simplify the complex declaration:

```go
alias 3ints = [3:int:3]
alias some_strings = [3:string:5]
alias weird_regex = "x"|"[a-zA-Z]+"
alias 2Dbool = [[bool]]

complex = [string|some_strings|3ints|[weird_regex]|2Dbool]
```

### Recursion

As ```efp``` elements are lazily validated against the prototype tree, recursion will not cause an infinite loop.  Recursion may be accomplied through the use of aliasing:

```go
alias h = hello {
    h
}
```

## Usage/Installation

To install, simply use:

```go
go get github.com/end-r/efp
```

There are two methods which must be called:

```go

import "github.com/end-r/efp"

func main(){
    p := efp.Prototype("standard.efp")
    e, errs := p.Parse("file.txt")
}
```

The full godoc for the efp can be found at:

For some example applications, check out:

- [FireVM](https://github.com/end-r/fireVM), a costly VM generator
- [Vox](https://github.com/end-r/vox), a configuration system for online elections
- [LexisTextUs](https://github.com/end-r/lexis-text-us), a legal form specification for a messenger chat bot

## Variable Elements and Fields

Sometimes, you might want an element with a user-defined key. The prototype syntax is as follows:

```go
"a-z" {

}
```

Where ```"a-z"``` is an example regular expression.

This system can also be used to allow duplicates within a parent element, using the ```<``` and ```>``` operators.

```go
// minimum of two elements matching this regex
<2:"a-z"> {

}

// precisely three elements matching this regex
<3:"a-i":3>(string, string){

}
```

The same syntax is valid for fields, e.g:

```go
<2:"key":3> = string
```

Note that the above declaration is identical to ```go <2:key:3> = string```, but the latter is validated through string equality instead of regex.

## Regex Overlaps

There is significant potential for error when dealing with field frequency in the following prototype:

```go
parent {
    "a-zA-Z" : string
    name : string
}
```

```go
parent{
    name = "hi"
}
```

As ```name``` matches the regex of the earlier field, there is a chance it would be interpreted as . However, in the ```efp```, string equality comparisons always take precedence over regex evaluations.

This gives rise to the issue of competing regex strings, e.g.:

```go
parent {
    "a-zA-Z" : string
    "a-z" : string
}
```

Of course, it would be possible to assign the field to the "best" or "simplest" regex key (by whatever metric), but this would result in additional overhead for each regex declaration, as well as ambiguity regarding these complexity comparisons. The expected behaviour, therefore, is to assign to the first matched regex. As the regexes are stored in a map (and therefore in an arbitrary order), this allocation is pseudo-random, and should NOT be relied upon or predicted. The best solution, therefore, is to avoid overlapping regex declarations within an element.

NOTE: This may change in the future, depending on feedback, or this control may become an option.

## Type System

When defining types, they are stored using the following structure:

```go
type typeDeclaration struct {
	isArray bool
	types   []*typeDeclaration
}
```

The following definitions and their representations:

```go
x : string

```

![Alt text](https://g.gravizo.com/svg?
  digraph G {
      aize ="4,4";
      string [shape=box];
  }
)

```go
x : string|int

```

![Alt text](https://g.gravizo.com/svg?
  digraph G {
      aize ="4,4";
      string [shape=box];
      int [shape=box];
  }
)

```go
x : [string|int]

```
![Alt text](https://g.gravizo.com/svg?
  digraph G {
    aize ="4,4";
    array [shape=box]
    array -> string [shape=box];
    array -> int [shape=box];
  }
)

```go
x : [string]|[int]

```

![Alt text](https://g.gravizo.com/svg?
  digraph G {
      aize ="4,4";
      array [shape=box]
      array -> array -> string [shape=box];
      array [shape=box]
      array -> int [shape=box];
  }
)

```go
x : [[string]]|[int]

```

![Alt text](https://g.gravizo.com/svg?
  digraph G {
      aize ="4,4";
      array [shape=box]
      array -> array [shape=box]
      array -> string [shape=box];
      array [shape=box]
      array -> int [shape=box];
  }
)

## Accessing Values

## Errors

| Error      | Explanation        |
| :-------------: |:-------------:|
| Alias Not Found | The alias is not visible in the current scope. |
| Duplicate Element |    |
| Duplicate Field |     |
| Invalid Token |     |
| Invalid Regex | The string specified cannot be transformed into golang regex.   |
| Duplicate Alias | The following alias has already been declared in the current element. |

## Full Example

```go
alias LETTERS = "a-zA-Z"

alias EVENT = <1:LETTERS> {
    gender : "[MF]"
}

alias PERSON = <1:LETTERS>{
    events {
        EVENT
    }
    accomodation {
        name : string!
        room : int
    }
}

olympics {
    city : string!
    year : int!
    athletes {
        PERSON
    }
    officials {
        PERSON
    }
}
```

A corresponding file:

```go
olympics {
    city = "Athens"
    year = 2004
    athletes {
        "Bill Gates" {
            events {
                "100m" {
                    gender = "M"
                }
            }
        }
        "Steve Jobs" {
            events {
                "100m" {
                    gender = "M"
                }
            }
        }
    }
}
```
