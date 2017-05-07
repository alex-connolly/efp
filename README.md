The element-field parser is a method of generating file validation parsers with a clean and consistent syntax.

## Basic Concepts

### Elements

Elements are one of the two key structures in the ```efp```.

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

# Process

First, a prototype tree is generated. This is the format which is specified in your ```.efp``` file. All files will then be compiled against this prototype tree for validation. An example prototype tree:

```go
fireVM {
    name : string!
}
```

This is then enforced in our files - all valid files must contain a top level ```fireVM``` element, and a ```name``` field, which takes a string value. In a prototype file, the ```!``` denotes a compulsory field.

Field declarations are made with the following syntax:

```go
key : type
```

Types may be in one or more of the following formats:

- Standard Aliases
- Regex
- Array Notation

## Standard Aliases

There are several default types in the ```efp``` spec:

```go
id = "[a-zA-Z_]+"
string = "\"[^()]\""
float = "[0-9]*.[0-9]+"
bool = "true|false"
int = "[0-9]+"
```

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

## Aliasing

It is possible to declare aliases within a efp file, with the normal scope boundaries. Aliases are tantamount to C macros, in that they simply perform a text replace. If the text contains an element, that element will be evaluated and added to the tree.

```go
alias x = 5
alias divs = divisions("name") {
    x
}
```

### Recursion

As ```efp``` elements are lazily validated against the prototype tree, recursion will not cause an infinite loop. Recursion may be accomplied through the use of aliasing:

```go
alias hello = hello {
    hello
}
```
