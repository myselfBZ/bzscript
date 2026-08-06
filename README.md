## Bzscript

BzScript - a toy programming language designed for beginners to grasp the basics of programming.


### How to get up and running
1. Compile the code
```sh
go build -o bin/bzscript ./cmd/bzscript
```
2. Run it with a bzscript file (which should end in .bzs)
```sh
./bin/bzscript [your_file].bzs
```

## Basic syntax
```js
// variable decleration with different types
var x = 1234 // intigers
var pi = 3.14 // floats
var name = "bzscript" // strings
var isNew = true // booleans
var forProd = false // booleans
```

```js
// functions
fun hello() {
    print("hello")
}

fun addOne(x) {
    return x + 1
}

fun isEven(x) {
    return x % 2 == 0
}
// higher-order functions, and closures
var f = fun (x) {
    return fun(y) {
        return x + y
    }
}
var addTwo = f(2)
addTwo(5) // 7
```

```js
// looping
var x = 1
while x < 10 {
    print("x is", x)
    x = x + 1
}
```

```js
// control flow
var age = 17
if age < 18 {
    print("You are not allowed")
} else {
    print("Welcome")
}
```

```js
// arrays
var shopping_list = ["apples", "bananas", "pineapple"]
var first_item = shoppingList[0]
var second_item = shoppingList[1]
var third_item = shoppingList[2]

// visiting every element
var i = 0
while i < len(shopping_list) {
    print(shopping_list[i])
        i = i + 1

}
```
```go
// structures, or grouped data
struct Human {
    name;
    age;
    social_security;

}
// var me = Huma{name: "myselfBZ", age: 19, social_security: "1235"} this not yet supported
var me = new(Human)
me.age = 19
me.name = "myselfBZ"
me.social_security = "1235"
print("Name: ", me.name)
print("Age: ", me.age)
print("Social Security: ", me.social_security)
```

# Coming soon...

```js
// maps
var scores = map{
    "John":89,
    "Sarah":88,

}
var entity = scores["John"]
if entity.exists {
    print("John has", entity.value, "scores")
} else {
    print("John is not on the scores map")
}
```

