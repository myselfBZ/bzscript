## Bzscript

BzScript - a toy programming language designed for beginners to grasp the basics of programming.

## Basic syntax

<style>
  .bzscript-code {
    background-color: #1e1e1e;
    color: #d4d4d4;
    font-family: 'Consolas', 'Fira Code', 'Courier New', monospace;
    font-size: 14px;
    line-height: 1.5;
    padding: 16px;
    border-radius: 8px;
    overflow-x: auto;
    white-space: pre;
  }
  .bz-kw { color: #569cd6; font-weight: bold; }   /* Keywords (var, fun, struct, if, else, while, return) */
  .bz-type { color: #4ec9b0; }                    /* Data Types/Constructs (map, Human) */
  .bz-str { color: #ce9178; }                     /* Strings */
  .bz-num { color: #b5cea8; }                     /* Numbers */
  .bz-bool { color: #569cd6; font-weight: bold;}  /* Booleans (true, false) */
  .bz-com { color: #6a9955; font-style: italic; } /* Comments */
  .bz-fn { color: #dcdcaa; }                      /* Built-in Functions */
  .bz-prop { color: #9cdcfe; }                    /* Map/Struct Keys & Properties */
</style>

<div class="bzscript-code">
<span class="bz-com">// Variable declaration with different types</span>
<span class="bz-kw">var</span> x = <span class="bz-num">1234</span> <span class="bz-com">// integers</span>
<span class="bz-kw">var</span> pi = <span class="bz-num">3.14</span> <span class="bz-com">// floats</span>
<span class="bz-kw">var</span> name = <span class="bz-str">"bzscript"</span> <span class="bz-com">// strings</span>
<span class="bz-kw">var</span> isNew = <span class="bz-bool">true</span> <span class="bz-com">// booleans</span>
<span class="bz-kw">var</span> forProd = <span class="bz-bool">false</span> <span class="bz-com">// booleans</span>

<span class="bz-com">// Functions</span>
<span class="bz-kw">fun</span> myfunc() {....}

<span class="bz-kw">fun</span> hello() {
    <span class="bz-fn">puts</span>(<span class="bz-str">"hello"</span>)
}

<span class="bz-kw">fun</span> addOne(x) {
    <span class="bz-kw">return</span> x + <span class="bz-num">1</span>
}

<span class="bz-kw">fun</span> isEven(x) {
    <span class="bz-kw">return</span> x % <span class="bz-num">2</span> == <span class="bz-num">0</span>
}

<span class="bz-com">// Looping</span>
<span class="bz-kw">var</span> x = <span class="bz-num">1</span>
<span class="bz-kw">while</span> x &lt; <span class="bz-num">10</span> {
    <span class="bz-fn">puts</span>(<span class="bz-str">"x is"</span>, x)
    x = x + <span class="bz-num">1</span>
}

<span class="bz-com">// Control flow</span>
<span class="bz-kw">var</span> age = <span class="bz-num">17</span>
<span class="bz-kw">if</span> age &lt; <span class="bz-num">18</span> {
    <span class="bz-fn">puts</span>(<span class="bz-str">"You are not allowed"</span>)
} <span class="bz-kw">else</span> {
    <span class="bz-fn">puts</span>(<span class="bz-str">"Welcome"</span>)
}

<span class="bz-com">// Structures, or grouped data</span>
<span class="bz-kw">struct</span> <span class="bz-type">Human</span> {
    name;
    age;
    social_security;
}

<span class="bz-kw">var</span> me = <span class="bz-type">Human</span>{<span class="bz-prop">name</span>: <span class="bz-str">"myselfBZ"</span>, <span class="bz-prop">age</span>: <span class="bz-num">19</span>, <span class="bz-prop">social_security</span>: <span class="bz-str">"$$%#$$%#"</span>}

<span class="bz-fn">puts</span>(<span class="bz-str">"Name: "</span>, me.name)
<span class="bz-fn">puts</span>(<span class="bz-str">"Age: "</span>, me.age)
<span class="bz-fn">puts</span>(<span class="bz-str">"Social Security: "</span>, me.social_security)

<span class="bz-com">// Maps</span>
<span class="bz-kw">var</span> scores = <span class="bz-type">map</span>{
    <span class="bz-str">"John"</span>: <span class="bz-num">89</span>,
    <span class="bz-str">"Sarah"</span>: <span class="bz-num">88</span>,
}

<span class="bz-kw">var</span> entity = scores[<span class="bz-str">"John"</span>]

<span class="bz-kw">if</span> entity.exists {
    <span class="bz-fn">puts</span>(<span class="bz-str">"John has"</span>, entity.value, <span class="bz-str">"scores"</span>)
} <span class="bz-kw">else</span> {
    <span class="bz-fn">puts</span>(<span class="bz-str">"John is not on the scores map"</span>)
}

<span class="bz-com">// Arrays</span>
<span class="bz-kw">var</span> shopping_list = [<span class="bz-str">"apples"</span>, <span class="bz-str">"bananas"</span>, <span class="bz-str">"pineapple"</span>]
<span class="bz-kw">var</span> first_item = shopping_list[<span class="bz-num">0</span>]
<span class="bz-kw">var</span> second_item = shopping_list[<span class="bz-num">1</span>]
<span class="bz-kw">var</span> third_item = shopping_list[<span class="bz-num">2</span>]

<span class="bz-com">// Visiting every element</span>
<span class="bz-kw">var</span> i = <span class="bz-num">0</span>
<span class="bz-kw">while</span> i &lt; <span class="bz-fn">len</span>(shopping_list) {
    <span class="bz-fn">puts</span>(shopping_list[i])
    i = i + <span class="bz-num">1</span>
}
</div>
