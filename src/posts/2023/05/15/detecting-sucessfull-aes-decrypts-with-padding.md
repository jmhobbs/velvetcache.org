---
tags:
- crypto
- golang
title: Detecting Sucessfull AES Decrypts With Padding
type: post
date: 2023-05-15
summary: >
  Unauthenticated AES works even if you have the wrong key. Here's a way to tell you have the right key.
---

When I was working on my [Authy CLI](/2023/05/12/the-authy-backup-system) I had a stack of existing ciphertexts, and needed to ensure the key used was correct.

AES in CBC mode doesn't do anything for authentication, it works purely on the blocks and doesn't care about what comes out.  As such, you can decrypt a message with the wrong key, and it will not be an error, it will just return garbage.

An example:

```go
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {
	correctPassphrase := "So I ask you, how now brown cow?"
	incorrectPassphrase := "I say unto thee, my friend. Moo."

	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	ciphertext, err := base64.StdEncoding.DecodeString("w5hNXNwLOHkVF0mmA4Mwz57oMIktUaAjJiiFX/gE770=")
	if err != nil {
		panic(err)
	}

	plaintext, err := decrypt([]byte(correctPassphrase), ciphertext, iv)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(plaintext))

	plaintext, err = decrypt([]byte(incorrectPassphrase), ciphertext, iv)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(plaintext))
}

func decrypt(passphrase, ciphertext, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}
```

```console
$ go run .
Thank you for the milk.
�>7pfO!��V�
           %
```

Luckily, there is a way to be sure things decrypted _if_ we know what padding was used.  This cipher must work on whole blocks, that is 16 bytes at a time.  If you pass it a buffer whose length is not divisible by 16, that is an error.

You can't see it in the output above, but the correct version is actually printing some invisible characters, since I'm not removing the padding:

```console
$ hexyl plaintext
┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│00000000│ 54 68 61 6e 6b 20 79 6f ┊ 75 20 66 6f 72 20 74 68 │Thank yo┊u for th│
│00000010│ 65 20 6d 69 6c 6b 2e 09 ┊ 09 09 09 09 09 09 09 09 │e milk._┊________│
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘
```

Those `0x09` bytes are the padding, as this message is using [PKCS #5](https://www.rfc-editor.org/rfc/rfc2898).

PKCS #5 is a very simple but effective padding scheme.  The number of padding bytes required is also the value of the padding bytes.

So in our case, the block size is `16`, our message length `% 16` is `7`, so our padding is `9`.

In the case where your message length matches up exactly to your block size, you pad it with the whole block size.

Some examples, all with block size of `16`, or `0x10`:

```
┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│00000000│ 48 65 6c 6c 6f 20 57 6f ┊ 72 6c 64 05 05 05 05 05 │Hello Wo┊rld•••••│
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘
┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│00000000│ 48 65 6c 6c 6f 20 57 6f ┊ 72 6c 64 2c 20 49 74 27 │Hello Wo┊rld, It'│
│00000010│ 73 20 4d 65 21 0b 0b 0b ┊ 0b 0b 0b 0b 0b 0b 0b 0b │s Me!•••┊••••••••│
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘
┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│00000000│ 48 65 6c 6c 6f 20 42 6c ┊ 75 65 20 57 6f 72 6c 64 │Hello Bl┊ue World│
│00000010│ 10 10 10 10 10 10 10 10 ┊ 10 10 10 10 10 10 10 10 │••••••••┊••••••••│
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘
```

We can leverage this to know when we have used the correct decryption key, since the last `N` bytes of the message should match the last byte of the message, where `N` is the integer value of the last byte of the message.

```golang
func removePKCS5Padding(size int, plaintext []byte) ([]byte, error) {
	padding := plaintext[len(plaintext)-1]
	paddingLength := int(padding)

	if paddingLength > size {
		return nil, fmt.Errorf("Padding value %02X larger than block size", padding)
	}

	for _, b := range plaintext[len(plaintext)-int(paddingLength):] {
		if b != padding {
			return nil, fmt.Errorf("Incorrect padding %02X", b)
		}
	}
	return plaintext[:len(plaintext)-paddingLength], nil
}
```

With this in our `decrypt` function, we now return an error for the incorrect passphrase.

```console
$ go run .
Thank you for the milk.
panic: Padding value EB larger than block size

goroutine 1 [running]:
main.main()
	/Users/jmhobbs/Working/decrypt/main.go:49 +0x37c
exit status 2
```

Additionally, we can see we have removed the padding on our message as well, returning us to the true original plaintext.

```console
$ hexyl plaintext-stripped
┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│00000000│ 54 68 61 6e 6b 20 79 6f ┊ 75 20 66 6f 72 20 74 68 │Thank yo┊u for th│
│00000010│ 65 20 6d 69 6c 6b 2e    ┊                         │e milk. ┊        │
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘
```
