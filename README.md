## Bzscript

BzScript - a toy programming language designed for beginners to grasp the basics of programming.

## Basic syntax

```bzscript
// variable decleration with different types
var x = 1234 // intigers
var pi = 3.14
var name = "bzscript"
var isNew = true
var forProd = false
```

```bzscript
// functions
fun myfunc() {....}

fun hello() {
    puts("hello")
}

fun addOne(x) {
    return x + 1
}

fun isEven(x) {
    return x % 2 == 0
}
```


```bzscript
// looping

var x = 1
while x < 10 {
    puts("x is", x)
    x = x + 1
}
```

```bzscript
// control flow
var age = 17
if age < 18 {
    puts("You are not allowed")
} else {
    puts("Welcome")
}
```

```bzscript
// structures, or grouped data
struct Human {
    name
    age
    social_security
}

var me = Human{name: "myselfBZ", age:19, social_security:"$$%#$$%#"}

puts("Name: ", me.name)
puts("Age: ", me.age)
puts("Social Security: ", me.social_security)
```

```bzscript
// maps
var scores = map{
    "John":89,
    "Sarah":88,
}

var entity = scores["John"]

if entity.exists {
    puts("John has", entity.val, "scores")
} else {
    puts("John is not on the scores map")
}
```


```bzscript
// arrays
var shopping_list = ["apples", "bananas", "pineapple"]
var first_item = shoppingList[0]
var second_item = shoppingList[1]
var third_item = shoppingList[2]

// visiting every element
var i = 0
while i < len(shopping_list) {
    puts(shopping_list[i])
    i = i + 1
}
```
