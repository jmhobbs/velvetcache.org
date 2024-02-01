---
category:
- Uncategorized
creator: admin
date: 2024-01-31
tags:
- golang
- Programming
title: CHIP-8 BCD
type: post
permalink: /2024/01/31/chip-8-bcd
summary: While implementing a CHIP-8 emulator, I came across a new concept.
---

Recently I've been implementing a [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) emulator in Go, from just
[the spec here](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM).  While implementing, I learned a new
concept from one of the instructions:

> **Fx33 - LD B, Vx**
>
> Store BCD representation of Vx in memory locations I, I+1, and I+2.

BCD here stands for Binary-coded Decimal, where a byte is used for each decimal digit in a number. In this case,
the 8 bit register `Vx` has it's hundreds place moved to `I`, tens to `I+1` and ones to `I+2`.

Heres an example, if your value to encode is `238`:


```
Vx  | 11101110 | 0xEE | 238
----+----------+------+----
I   | 00000010 | 0x02 |   2
I+1 | 00000011 | 0x03 |   3
I+2 | 00001000 | 0x08 |   8
```

Neat! But, I'm not sure what it's useful for in CHIP-8.  The [Wikipedia page](https://en.wikipedia.org/wiki/Binary-coded_decimal#Comparison_with_pure_binary)
has a list of advantages, but I'll have to poke around some CHIP-8 programs to see how it is used.

# Deriving BCD digits

My implementation is the most straightforward I could think of.  There is probably some bit twiddly way to 
do it, but I'm not great at those things.  My method was to modulo out the upper digits, and divide what remains
to get each digit independently.

```golang
m.memory[m.I] = m.V[x(op)] / 100
m.memory[m.I+1] = (m.V[x(op)] % 100) / 10
m.memory[m.I+2] = ((m.V[x(op)] % 100) % 10) / 1
```

# Clock

After seeing the visualization on the Wikipedia page for BCD, I realized I have [a clock from ThinkGeek](https://web.archive.org/web/20070210215459/http://www.thinkgeek.com/homeoffice/lights/59e0/) 
(remember ThinkGeek?!) that was a BCD representation of the time.  I'm sure I still have it in a box 
somewhere, but it's at least 15+ years old at this point.

![BCD Clock](https://static.velvetcache.org/pages/2024/01/31/chip-8-bcd/led-binclock-described.jpg)
