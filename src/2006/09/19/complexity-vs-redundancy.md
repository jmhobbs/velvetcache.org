---
category:
- Geek
creator: admin
date: 2006-09-20
layout: layout.njk
permalink: /2006/09/19/complexity-vs-redundancy/
tags:
- Computers
- Java
- Musing
- Programming
- School
title: Complexity Vs. Redundancy
type: post
wp_id: "38"
---
I just wrote the following line of code as part of an implementation of a bucket sort.

```cpp
bucket[(array[j]/static_cast<int>(pow(10.0,y)))%BUCKETS].Enqueue(array[j]);
```

Somehow I feel wrong inside about it.  Look at it, it's fairly complicated.  In the context of the program you _could_ work out what it was doing, but the effort would be made easier if it was rewritten as this.

```cpp
int temp = ( array[j] / static_cast<int>( pow(10.0, y) ) ) % BUCKETS;
bucket[temp].Enqueue(array[j]);
```

Just splitting up that nasty math part from the method call clear up a bit of the ambiguity of what the line is doing, and the next version is even better.

```cpp
// Returns a set number of digits from an
// integer in the array, based on the number
// of buckets and the current iteration.
// For example:
//     array[j] = 123456;
//     y = 4;
//     BUCKETS = 100;
// Then:
//     temp = 12;
int temp = ( array[j] / static_cast<int>( pow(10.0, y) ) ) % BUCKETS;
bucket[temp].Enqueue(array[j]);
```

While this doesn't _really_ tell you whats going on, it can lead you to finding the pattern involved.  That pattern is described thus:

_For an integer **N**, the digit **X** positions to the left of the least significant digit is: **(N/10^X)%(10^Y)**  Where **Y** is a "look-ahead" to retrieve the next-**Y** digits as well._

Thats my informal definition mix for base 10 numbers, derived from a few dozen minutes of fighting with [bc](http://www.die.net/doc/linux/man/man1/bc.1.html).  I _could_ put this information into the program, and then deciphering it's workings would be a snap.  But why?  I mean, I'm never going to use this again, and if I need to, I worked it out the first time, the second shouldn't be any harder.  How much documentation is too much?  How much is too little?

I hate bloated files, sources that are so filled with comments you never see any code.  I think (like many before me) that code should be self-documenting as much as possible.  I don't belive in using in-line comments.  They get in the way and don't often help.  But I've also always thought you shouldn't break out the comment blocks for anything that wasn't "serious", so do I lower the definition of a "serious" comment, or do I throw away my would-be in-liners?

I suppose this is one of those things that comes with time, but I've tried to maintain old code before, that someone else wrote with zero comments, and it's impossible.  I think I'll err on the side of caution, for whoever comes after me.
