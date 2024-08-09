---
date: 2024-08-08
tags:
- Plaintext
- golang
- Unicode
- Encoding
title: Exploring UTF-8 Encoding
type: post
permalink: /2024/08/08/exploring-utf-8-encoding/
summary: Let's use Go to poking at UTF-8 encoded text and better understand the most common encoding in the world.
---

[A coworker](https://michaelarestad.com/) shared a great conference talk by Dylan Beattie called "[Plain Text](https://www.youtube.com/watch?v=gd5uJ7Nlvvo)".  It's well worth your time to watch if you're at all interested in the internals of your computer, and the history and quirks behind them.

I've long known the design of UTF-8, but never played with it, so I wrote a quick demo that picks apart input byte by byte and groups the bytes according to UTF-8 encoding.

The first check I do is to see if it is a "normal" ASCII range byte.  This is denoted in UTF-8 by the high bit being a `0`. If it is, I print it as is and continue.

```go
if currentByte&0b10000000 == 0 {
  fmt.Printf("U+%04X\t%s\t%08b\n", currentByte, string(currentByte), currentByte)
}
```

After that, I check if the byte is a new marker byte, the first byte in a UTF-8 sequence.  These bytes have the highest two bits set.  If it is, I add it to a slice for later, and then determine how many continuation bytes we should expect.  This is based on how many of the high bits in a row are `1`.

```go
if currentByte&0b11000000 == 0b11000000 {
  if len(bytes) != 0 {
    fmt.Println("ERROR: Unexpected multibyte sequence")
    break
  }
  bytes = []byte{currentByte}

  if currentByte>>5 == 0b110 {
    expectedLength = 2
  } else if currentByte>>4 == 0b1110 {
    expectedLength = 3
  } else if currentByte>>3 == 0b11110 {
    expectedLength = 4
  } else {
    fmt.Println("ERROR: Invalid multibyte sequence")
    break
  }
}
```

Finally, anything left over is going to be a continuation bytes, one that starts with `10` in the highest two bits.  We make sure we are in the midst of a multibyte sequence, and then these are tacked onto a slice of bytes which are being built into the unicode value. Finally, if the sequence is complete we print it out and reset everything.

```go
if len(bytes) == 0 {
  fmt.Println("ERROR: Unexpected continuation byte")
  break
}
bytes = append(bytes, currentByte)

if len(bytes) == expectedLength {
  printMultibyte(bytes)
  bytes = []byte{}
}
```

The `printMultibyte` grabs the bits from the initial byte, and continuation bytes, and glues them together to create the complete unicode codepoint byt shifting the initial byte based on the length of the sequence to remove the high flag bits, and then appending each continuation byte, again, without the two high flag bits.

```go
var length int = len(bytes)

var codepoint int = (int(bytes[0]<<(length+1)) & 0b11111111) >> (length + 1)
for _, b := range bytes[1:] {
  codepoint = codepoint << 6
  codepoint |= int(b & 0b00111111)
}
```

As an example, here is the math building up the [white heavy check mark](https://codepoints.net/U+2705):

```
Bytes: 11100010 10011100 10000101

Codepoint = 11100010 << (3 + 1) = 00001110 00100000
Codepoint = 00001110 00100000 & 11111111 = 00100000
Codepoint = 00100000 >> (3 + 1) = 00000010

Continuation Byte 1: 10011100

  Codepoint = 00000010 << 6 = 10000000
  Continuation Byte Masked: 10011100 & 00111111 = 00011100
  Codepoint = 10000000 | 00011100 = 10011100

Continuation Byte 2: 10000101

  Codepoint = 10011100 << 6 = 00100111 00000000
  Continuation Byte Masked: 10000101 & 00111111 = 00000101
  Codepoint = 00100111 00000000 | 00000101 = 00100111 00000101

Final Codepoint: 00100111 00000101 = 0x2705
```

Putting this all together, we can send in some UTF-8 text and inspect the composition:

```
$ echo -n "Hello UTF-8! ✅" | go run main.go
Codepoint    Print    Bytes
U+0048       H        01001000
U+0065       e        01100101
U+006C       l        01101100
U+006C       l        01101100
U+006F       o        01101111
U+0020                00100000
U+0055       U        01010101
U+0054       T        01010100
U+0046       F        01000110
U+002D       -        00101101
U+0038       8        00111000
U+0021       !        00100001
U+0020                00100000
U+2705       ✅       11100010 10011100 10000101
```

Here's the complete source. I expect there are some clever tricks I'm missing, but I've never been great at bit twiddling, so those are left as an exercise for the reader 😉

```go
package main

import (
  "bufio"
  "fmt"
  "io"
  "os"
)

func main() {
  fmt.Println("Codepoint\tPrint\tBytes")

  var (
    bytes          []byte
    expectedLength int

    currentByte byte
    err         error
    stdin       *bufio.Reader = bufio.NewReader(os.Stdin)
  )

  for {
    currentByte, err = stdin.ReadByte()
    if err == io.EOF {
      break
    }

    if currentByte&0b10000000 == 0 {
      fmt.Printf("U+%04X\t%s\t%08b\n", currentByte, string(currentByte), currentByte)
    } else if currentByte&0b11000000 == 0b11000000 {
      if len(bytes) != 0 {
        fmt.Println("ERROR: Unexpected multibyte sequence")
        break
      }

      bytes = []byte{currentByte}

      if currentByte>>5 == 0b110 {
        expectedLength = 2
      } else if currentByte>>4 == 0b1110 {
        expectedLength = 3
      } else if currentByte>>3 == 0b11110 {
        expectedLength = 4
      } else {
        fmt.Println("ERROR: Invalid multibyte sequence")
        break
      }
    } else {
      if len(bytes) == 0 {
        fmt.Println("ERROR: Unexpected continuation byte")
        break
      }
      bytes = append(bytes, currentByte)

      if len(bytes) == expectedLength {
        printMultibyte(bytes)
        bytes = []byte{}
      }
    }
  }
}

func printMultibyte(bytes []byte) bool {
  var length int = len(bytes)

  var codepoint int = (int(bytes[0]<<(length+1)) & 0b11111111) >> (length + 1)
  for _, b := range bytes[1:] {
    codepoint = codepoint << 6
    codepoint |= int(b & 0b00111111)
  }

  var args []any
  var fmtstring string = "U+%04X\t%s\t"
  args = append(args, codepoint)
  args = append(args, string(bytes))

  for i := 0; i < length; i++ {
    fmtstring += "%08b "
    args = append(args, bytes[i])
  }
  fmtstring += "\n"

  fmt.Printf(fmtstring, args...)

  return true
}
```

