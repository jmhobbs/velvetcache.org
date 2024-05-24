---
category:
- Uncategorized
creator: admin
date: 2024-05-23
tags:
- golang
- Gaming
- protocols
title: DayZ Server Browsers
type: post
permalink: /2024/05/23/dayz-server-browsers/
opengraph_image: /static/og/2024-05-23-dayz-server-browsers.jpg
summary: >
    Recently I've been playing a lot of DayZ, through the DayZSA launcher. But, where
    do these server listings live, and how does a thrid-party launcher get them?
---


Recently I've been playing a lot of DayZ, often through the DayZSA launcher. But I was curious, where do these server listings live, and how does a thrid-party launcher get them?  Further,
how does a website like [BattleMetrics](https://www.battlemetrics.com/servers/dayz) get them and list it online?

This question led me to some very, very weird data structures, jammed into weird protocols, that feel utterly undesigned. And yet, they work!

# A2S, or Steam Server Queries

Many multiplayer games use a protocol provided for in the Steam SDK commonly known as A2S. Some light searching revealed that DayZ uses this protcol. It is UDP based, and loosely documented, 
the best I could find was [this page on the Valve wiki](https://developer.valvesoftware.com/wiki/Server_queries)

If you skim that, you'll see it feels very ad-hoc, like it started small and simple and things just got glommed onto it over time.  So let's try to talk to a DayZ server!

First, we need a server.  To keep things simple, we will use an official server, [`NY 6053`](https://www.battlemetrics.com/servers/dayz/5489323).  At time of writing it's query port was at `31.214.136.147:10101`.

Now, we need to know how what it expects.  We will start with `A2S_INFO`, to get the map, players etc.

The query starts with a four byte packet type header, `0xFFFFFFFF`.  This signifies it is a single packet request or response, not split across multiple packets.

Next, there is a single byte header indicating the request type, in this case `0x54` for `A2S_INFO`.  Finally, we have the query string, which is `Source Engine Query\0`.  Strings in this protocol are null terminated, like in C.

So the expected query packet is:

{%- # a2s_query.bin -%}

<pre><code class="ansi">┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│<span class="ansi ansi-fg-bright-black">00000000</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-yellow">ff ff ff ff </span><span class="ansi ansi-fg-cyan">54 53 6f 75</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">72 63 65 </span><span class="ansi ansi-fg-green">20 </span><span class="ansi ansi-fg-cyan">45 6e 67 69</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-yellow">××××</span><span class="ansi ansi-fg-cyan">TSou</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">rce</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-cyan">Engi</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000010</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">6e 65 </span><span class="ansi ansi-fg-green">20 </span><span class="ansi ansi-fg-cyan">51 75 65 72 79</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-bright-black">00                     </span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">ne</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-cyan">Query</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-bright-black">⋄       </span><span class="ansi ansi-fg-default">│
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘</code></pre>

Let's send that with netcat,

```
$ printf '\xFF\xFF\xFF\xFF\x54Source Engine Query\0' | nc --udp -x 31.214.136.147 10101
Sent 25 bytes to the socket
00000000  FF FF FF FF  54 53 6F 75  72 63 65 20  45 6E 67 69  ....TSource Engi
00000010  6E 65 20 51  75 65 72 79  00                        ne Query.
AlReceived 9 bytes from the socket
00000000  FF FF FF FF  41 6A 81 08  6C                        ....Aj..l
```

That...is not a long response.  It turns out, some A2S servers use a challenge on the query, something that was tacked onto the protocol later on.  This is a DDoS prevention mechanism for a reflection attack.

Before we can handle this challenge, let's take a look at the response packet.  We have the 4 byte header indicating a single packet response.  Then `0x41`, which is our response type, `S2C_CHALLENGE`. Finally, there is the four byte challenge. This is the value we need to append to our original packet to authenticate with the server.

{%- # a2s_query_challenge.bin -%}

<pre><code class="ansi">┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│<span class="ansi ansi-fg-bright-black">00000000</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-yellow">ff ff ff ff </span><span class="ansi ansi-fg-cyan">41 6a </span><span class="ansi ansi-fg-yellow">81 </span><span class="ansi ansi-fg-green">08</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">6c                     </span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-yellow">××××</span><span class="ansi ansi-fg-cyan">Aj</span><span class="ansi ansi-fg-yellow">×</span><span class="ansi ansi-fg-green">•</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">l       </span><span class="ansi ansi-fg-default">│
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘</code></pre>

At this point, we've reached the end of what is really viable to do with netcat, since we need to send raw bytes in response, and our little `printf` shenanigans won't get us there anymore.  Let's write a little client in Go instead!

```golang
package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func main() {
	var queryPacket bytes.Buffer
	// packet type
	queryPacket.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF})
	// request type
	queryPacket.WriteByte(0x54)
	// payload
	queryPacket.WriteString("Source Engine Query\000")

	conn, err := net.Dial("udp", "31.214.136.147:10101")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	response := sendPacket(conn, queryPacket.Bytes())

	if response[4] == 0x41 {
		fmt.Println("Challenge packet received")
		challenge := response[5:]
		fmt.Printf("Challenge: %v\n", challenge)

		// add the challenge bytes
		queryPacket.Write(challenge)

		response = sendPacket(conn, queryPacket.Bytes())
	} else {
		fmt.Println("Challenge packet not received")
	}

	fmt.Printf("Response: %v\n", response)
}

func sendPacket(conn net.Conn, pkt []byte) []byte {
	_, err := conn.Write(pkt)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1024)

	n, err := conn.Read(response)
	if err != nil {
		panic(err)
	}

	return response[:n]
}
```

There's a lot of fluff in there, and a lot of assumptions, but look, it works!

```
$ go run main.go
Challenge packet received
Challenge: [53 154 72 189]
Response: [255 255 255 255 73 17 68 97.....]
```

## A2S_INFO - Basic server information

Now that we have some bytes, let's decode it.  Here it is again in as a hex dump.

{%- # a2s_info.bin -%}

<pre><code class="ansi">┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│<span class="ansi ansi-fg-bright-black">00000000</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-yellow">ff ff ff ff </span><span class="ansi ansi-fg-cyan">49 </span><span class="ansi ansi-fg-green">11 </span><span class="ansi ansi-fg-cyan">44 61</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">79 5a </span><span class="ansi ansi-fg-green">20 </span><span class="ansi ansi-fg-cyan">55 53 </span><span class="ansi ansi-fg-green">20 </span><span class="ansi ansi-fg-cyan">2d </span><span class="ansi ansi-fg-green">20</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-yellow">××××</span><span class="ansi ansi-fg-cyan">I</span><span class="ansi ansi-fg-green">•</span><span class="ansi ansi-fg-cyan">Da</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">yZ</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-cyan">US</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-cyan">-</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000010</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">4e 59 </span><span class="ansi ansi-fg-green">20 </span><span class="ansi ansi-fg-cyan">36 30 35 33 </span><span class="ansi ansi-fg-green">20</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">28 31 73 74 </span><span class="ansi ansi-fg-green">20 </span><span class="ansi ansi-fg-cyan">50 65 72</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">NY</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-cyan">6053</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">(1st</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-cyan">Per</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000020</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">73 6f 6e </span><span class="ansi ansi-fg-green">20 </span><span class="ansi ansi-fg-cyan">4f 6e 6c 79</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">29 </span><span class="ansi ansi-fg-bright-black">00 </span><span class="ansi ansi-fg-cyan">63 68 65 72 6e 61</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">son</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-cyan">Only</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">)</span><span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-cyan">cherna</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000030</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">72 75 73 70 6c 75 73 </span><span class="ansi ansi-fg-bright-black">00</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">64 61 79 7a </span><span class="ansi ansi-fg-bright-black">00 </span><span class="ansi ansi-fg-cyan">44 61 79</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">rusplus</span><span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">dayz</span><span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-cyan">Day</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000040</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">5a </span><span class="ansi ansi-fg-bright-black">00 00 00 </span><span class="ansi ansi-fg-cyan">23 3c </span><span class="ansi ansi-fg-bright-black">00 </span><span class="ansi ansi-fg-cyan">64</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">77 </span><span class="ansi ansi-fg-bright-black">00 </span><span class="ansi ansi-fg-green">01 </span><span class="ansi ansi-fg-cyan">31 2e 32 33 2e</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">Z</span><span class="ansi ansi-fg-bright-black">⋄⋄⋄</span><span class="ansi ansi-fg-cyan">#<</span><span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-cyan">d</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">w</span><span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-green">•</span><span class="ansi ansi-fg-cyan">1.23.</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000050</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">31 35 37 30 34 35 </span><span class="ansi ansi-fg-bright-black">00 </span><span class="ansi ansi-fg-yellow">b1</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">74 27 </span><span class="ansi ansi-fg-green">03 </span><span class="ansi ansi-fg-cyan">3c </span><span class="ansi ansi-fg-yellow">ad 93 b4 </span><span class="ansi ansi-fg-cyan">62</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">157045</span><span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-yellow">×</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">t'</span><span class="ansi ansi-fg-green">•</span><span class="ansi ansi-fg-cyan"><</span><span class="ansi ansi-fg-yellow">×××</span><span class="ansi ansi-fg-cyan">b</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000060</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">40 </span><span class="ansi ansi-fg-green">01 </span><span class="ansi ansi-fg-cyan">62 61 74 74 6c 65</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">79 65 2c 6e 6f 33 72 64</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">@</span><span class="ansi ansi-fg-green">•</span><span class="ansi ansi-fg-cyan">battle</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">ye,no3rd</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000070</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">2c 73 68 61 72 64 30 30</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">31 2c 6c 71 73 30 2c 65</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">,shard00</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">1,lqs0,e</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000080</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">74 6d 34 2e 32 30 30 30</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">30 30 2c 65 6e 74 6d 34</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">tm4.2000</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">00,entm4</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">00000090</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-cyan">2e 30 30 30 30 30 30 2c</span><span class="ansi ansi-fg-default"> ┊ </span><span class="ansi ansi-fg-cyan">31 34 3a 30 39 </span><span class="ansi ansi-fg-bright-black">00 </span><span class="ansi ansi-fg-yellow">ac </span><span class="ansi ansi-fg-cyan">5f</span><span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-cyan">.000000,</span><span class="ansi ansi-fg-default">┊</span><span class="ansi ansi-fg-cyan">14:09</span><span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-yellow">×</span><span class="ansi ansi-fg-cyan">_</span><span class="ansi ansi-fg-default">│
│</span><span class="ansi ansi-fg-bright-black">000000a0</span><span class="ansi ansi-fg-default">│ </span><span class="ansi ansi-fg-green">03 </span><span class="ansi ansi-fg-bright-black">00 00 00 00 00      </span><span class="ansi ansi-fg-default"> ┊                        <span class="ansi ansi-fg-default"> │</span><span class="ansi ansi-fg-green">•</span><span class="ansi ansi-fg-bright-black">⋄⋄⋄⋄⋄  </span><span class="ansi ansi-fg-default">┊        <span class="ansi ansi-fg-default">│
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘</code></pre>

It starts with the "single packet" preamble which we've seen on every message so far.  The response type on this one is `0x49`, which is the `A2S_INFO` response value.  Next we have a single byte for "protocol". I am unsure what this is honestly, there is nothing in the docs. After that is the name, a `null` terminated  string, `DayZ US - NY 6053 (1st Person Only)`.  Then map, folder, and game name, all strings.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━━━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000000</span>│ <span class="ansi ansi-fg-bright-yellow">ff</span> <span class="ansi ansi-fg-bright-yellow">ff</span> <span class="ansi ansi-fg-bright-yellow">ff</span> <span class="ansi ansi-fg-bright-yellow">ff</span>             ┊                         │<span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-yellow">x</span>    ┊        │ Response Type ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000004</span>│ <span class="ansi ansi-fg-bright-cyan">49</span>                      ┊                         │<span class="ansi ansi-fg-bright-cyan">I</span>       ┊        │ Packet Type   ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000005</span>│ <span class="ansi ansi-fg-bright-yellow">11</span>                      ┊                         │<span class="ansi ansi-fg-bright-yellow">x</span>       ┊        │ Protocol      ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000006</span>│ <span class="ansi ansi-fg-bright-cyan">44</span> <span class="ansi ansi-fg-bright-cyan">61</span> <span class="ansi ansi-fg-bright-cyan">79</span> <span class="ansi ansi-fg-bright-cyan">5a</span> <span class="ansi ansi-fg-green">20</span> <span class="ansi ansi-fg-bright-cyan">55</span> <span class="ansi ansi-fg-bright-cyan">53</span> <span class="ansi ansi-fg-green">20</span> ┊ <span class="ansi ansi-fg-bright-cyan">2d</span> <span class="ansi ansi-fg-green">20</span> <span class="ansi ansi-fg-bright-cyan">4e</span> <span class="ansi ansi-fg-bright-cyan">59</span> <span class="ansi ansi-fg-green">20</span> <span class="ansi ansi-fg-bright-cyan">36</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">35</span> │<span class="ansi ansi-fg-bright-cyan">D</span><span class="ansi ansi-fg-bright-cyan">a</span><span class="ansi ansi-fg-bright-cyan">y</span><span class="ansi ansi-fg-bright-cyan">Z</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-bright-cyan">U</span><span class="ansi ansi-fg-bright-cyan">S</span><span class="ansi ansi-fg-green"> </span>┊<span class="ansi ansi-fg-bright-cyan">-</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-bright-cyan">N</span><span class="ansi ansi-fg-bright-cyan">Y</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-bright-cyan">6</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">5</span>│ Name          ┃
┃<span class="ansi ansi-fg-bright-black">00000022</span>│ <span class="ansi ansi-fg-bright-cyan">33</span> <span class="ansi ansi-fg-green">20</span> <span class="ansi ansi-fg-bright-cyan">28</span> <span class="ansi ansi-fg-bright-cyan">31</span> <span class="ansi ansi-fg-bright-cyan">73</span> <span class="ansi ansi-fg-bright-cyan">74</span> <span class="ansi ansi-fg-green">20</span> <span class="ansi ansi-fg-bright-cyan">50</span> ┊ <span class="ansi ansi-fg-bright-cyan">65</span> <span class="ansi ansi-fg-bright-cyan">72</span> <span class="ansi ansi-fg-bright-cyan">73</span> <span class="ansi ansi-fg-bright-cyan">6f</span> <span class="ansi ansi-fg-bright-cyan">6e</span> <span class="ansi ansi-fg-green">20</span> <span class="ansi ansi-fg-bright-cyan">4f</span> <span class="ansi ansi-fg-bright-cyan">6e</span> │<span class="ansi ansi-fg-bright-cyan">3</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-bright-cyan">(</span><span class="ansi ansi-fg-bright-cyan">1</span><span class="ansi ansi-fg-bright-cyan">s</span><span class="ansi ansi-fg-bright-cyan">t</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-bright-cyan">P</span>┊<span class="ansi ansi-fg-bright-cyan">e</span><span class="ansi ansi-fg-bright-cyan">r</span><span class="ansi ansi-fg-bright-cyan">s</span><span class="ansi ansi-fg-bright-cyan">o</span><span class="ansi ansi-fg-bright-cyan">n</span><span class="ansi ansi-fg-green"> </span><span class="ansi ansi-fg-bright-cyan">O</span><span class="ansi ansi-fg-bright-cyan">n</span>│               ┃
┃<span class="ansi ansi-fg-bright-black">00000038</span>│ <span class="ansi ansi-fg-bright-cyan">6c</span> <span class="ansi ansi-fg-bright-cyan">79</span> <span class="ansi ansi-fg-bright-cyan">29</span> <span class="ansi ansi-fg-bright-black">00</span>             ┊                         │<span class="ansi ansi-fg-bright-cyan">l</span><span class="ansi ansi-fg-bright-cyan">y</span><span class="ansi ansi-fg-bright-cyan">)</span><span class="ansi ansi-fg-bright-black">⋄</span>    ┊        │               ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000042</span>│ <span class="ansi ansi-fg-bright-cyan">63</span> <span class="ansi ansi-fg-bright-cyan">68</span> <span class="ansi ansi-fg-bright-cyan">65</span> <span class="ansi ansi-fg-bright-cyan">72</span> <span class="ansi ansi-fg-bright-cyan">6e</span> <span class="ansi ansi-fg-bright-cyan">61</span> <span class="ansi ansi-fg-bright-cyan">72</span> <span class="ansi ansi-fg-bright-cyan">75</span> ┊ <span class="ansi ansi-fg-bright-cyan">73</span> <span class="ansi ansi-fg-bright-cyan">70</span> <span class="ansi ansi-fg-bright-cyan">6c</span> <span class="ansi ansi-fg-bright-cyan">75</span> <span class="ansi ansi-fg-bright-cyan">73</span> <span class="ansi ansi-fg-bright-black">00</span>       │<span class="ansi ansi-fg-bright-cyan">c</span><span class="ansi ansi-fg-bright-cyan">h</span><span class="ansi ansi-fg-bright-cyan">e</span><span class="ansi ansi-fg-bright-cyan">r</span><span class="ansi ansi-fg-bright-cyan">n</span><span class="ansi ansi-fg-bright-cyan">a</span><span class="ansi ansi-fg-bright-cyan">r</span><span class="ansi ansi-fg-bright-cyan">u</span>┊<span class="ansi ansi-fg-bright-cyan">s</span><span class="ansi ansi-fg-bright-cyan">p</span><span class="ansi ansi-fg-bright-cyan">l</span><span class="ansi ansi-fg-bright-cyan">u</span><span class="ansi ansi-fg-bright-cyan">s</span><span class="ansi ansi-fg-bright-black">⋄</span>  │ Map           ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000056</span>│ <span class="ansi ansi-fg-bright-cyan">64</span> <span class="ansi ansi-fg-bright-cyan">61</span> <span class="ansi ansi-fg-bright-cyan">79</span> <span class="ansi ansi-fg-bright-cyan">7a</span> <span class="ansi ansi-fg-bright-black">00</span>          ┊                         │<span class="ansi ansi-fg-bright-cyan">d</span><span class="ansi ansi-fg-bright-cyan">a</span><span class="ansi ansi-fg-bright-cyan">y</span><span class="ansi ansi-fg-bright-cyan">z</span><span class="ansi ansi-fg-bright-black">⋄</span>   ┊        │ Folder        ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000061</span>│ <span class="ansi ansi-fg-bright-cyan">44</span> <span class="ansi ansi-fg-bright-cyan">61</span> <span class="ansi ansi-fg-bright-cyan">79</span> <span class="ansi ansi-fg-bright-cyan">5a</span> <span class="ansi ansi-fg-bright-black">00</span>          ┊                         │<span class="ansi ansi-fg-bright-cyan">D</span><span class="ansi ansi-fg-bright-cyan">a</span><span class="ansi ansi-fg-bright-cyan">y</span><span class="ansi ansi-fg-bright-cyan">Z</span><span class="ansi ansi-fg-bright-black">⋄</span>   ┊        │ Game          ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━━━━━━━┛</code></pre>

After that is a short, a two byte integer, for the Steam App ID.  In this response that is `00 00`, which doesn't match up with DayZ's actual app id, [221100](https://steamdb.info/app/221100/info/).  This exposes a flaw in the protocol.  The maximum value you can store in an unsigned short is `65535`, which is way smaller than our app id.  So since it doesn't fit, the server returns a `0`.  We will come back to this later on.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━━━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000066</span>│ <span class="ansi ansi-fg-bright-black">00</span> <span class="ansi ansi-fg-bright-black">00</span>                   ┊                         │<span class="ansi ansi-fg-bright-black">⋄</span><span class="ansi ansi-fg-bright-black">⋄</span>      ┊        │ Steam App ID  ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━━━━━━━┛</code></pre>

Following this are byte fields for number of players, maximum player count, and bots on the server.  At the time there were 35/60 players and no bots.  We also get bytes indicating if the server is dedicated, what OS it's hosted on, if it is password protected, and if it is [VAC](https://developer.valvesoftware.com/wiki/Valve_Anti-Cheat) enabled.  There's a weird mix here between using ASCII characters (`w` for Windows servers!) and numeric values (`1` means VAC is on!).  I understand the reasoning, but this is not a human readable format, so it's a bit odd to mix it up like that.  Lastly we get a version number, as a string.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000068</span>│ <span class="ansi ansi-fg-bright-cyan">23</span>                      ┊                         │<span class="ansi ansi-fg-bright-cyan">#</span>       ┊        │ Players     ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000069</span>│ <span class="ansi ansi-fg-bright-cyan">3c</span>                      ┊                         │<span class="ansi ansi-fg-bright-cyan"><</span>       ┊        │ Max Players ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000070</span>│ <span class="ansi ansi-fg-bright-black">00</span>                      ┊                         │<span class="ansi ansi-fg-bright-black">⋄</span>       ┊        │ Bots        ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000071</span>│ <span class="ansi ansi-fg-bright-cyan">64</span>                      ┊                         │<span class="ansi ansi-fg-bright-cyan">d</span>       ┊        │ Server Type ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000072</span>│ <span class="ansi ansi-fg-bright-cyan">77</span>                      ┊                         │<span class="ansi ansi-fg-bright-cyan">w</span>       ┊        │ Environment ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000073</span>│ <span class="ansi ansi-fg-bright-black">00</span>                      ┊                         │<span class="ansi ansi-fg-bright-black">⋄</span>       ┊        │ Visibility  ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000074</span>│ <span class="ansi ansi-fg-bright-yellow">01</span>                      ┊                         │<span class="ansi ansi-fg-bright-yellow">x</span>       ┊        │ VAC         ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃<span class="ansi ansi-fg-bright-black">00000075</span>│ <span class="ansi ansi-fg-bright-cyan">31</span> <span class="ansi ansi-fg-bright-cyan">2e</span> <span class="ansi ansi-fg-bright-cyan">32</span> <span class="ansi ansi-fg-bright-cyan">33</span> <span class="ansi ansi-fg-bright-cyan">2e</span> <span class="ansi ansi-fg-bright-cyan">31</span> <span class="ansi ansi-fg-bright-cyan">35</span> <span class="ansi ansi-fg-bright-cyan">37</span> ┊ <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">34</span> <span class="ansi ansi-fg-bright-cyan">35</span> <span class="ansi ansi-fg-bright-black">00</span>             │<span class="ansi ansi-fg-bright-cyan">1</span><span class="ansi ansi-fg-bright-cyan">.</span><span class="ansi ansi-fg-bright-cyan">2</span><span class="ansi ansi-fg-bright-cyan">3</span><span class="ansi ansi-fg-bright-cyan">.</span><span class="ansi ansi-fg-bright-cyan">1</span><span class="ansi ansi-fg-bright-cyan">5</span><span class="ansi ansi-fg-bright-cyan">7</span>┊<span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">4</span><span class="ansi ansi-fg-bright-cyan">5</span><span class="ansi ansi-fg-bright-black">⋄</span>    │ Version     ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━━━━━┛</code></pre>

Now we get to an interesting byte, the EDF, or Extra Data Flag. This is a bitfield which indicates more (and which) fields follow.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000087</span>│ <span class="ansi ansi-fg-bright-yellow">b1</span>                      ┊                         │<span class="ansi ansi-fg-bright-yellow">x</span>       ┊        │ EDF ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━┛</code></pre>

Our EDF byte is `0xb1`, which is `10110001`.  We take this value and apply a bitwise and operation, `&`, with other values to check which flags are set.  First up is `0x80`, or `10000000`.

```
  10110001
& 10000000
----------
  10000000
```

This flag is on! This means there will be a short containing the servers game port.  We do this with `0x10`, `0x40`, `0x20`, `0x01` and can parse the remaining bytes.  We have game port, Steam ID, keywords and game id.

Game port gives us our first meaningful use of a multi-byte integer value (since Steam App ID was just `00 00`).  We have two ways to interpret this value: little-endian or big-endian.  Essentially, which "end" of a number comes first? The docs tell us what the byte order will be little-endian, which is unusual for a network protocol, they tend towards big-endian.  This means `74 27` is the integer value `10100` (don't be fooled, it's not binary!).  If it were big-endian the value would be `29735`.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000088</span>│ <span class="ansi ansi-fg-bright-cyan">74</span> <span class="ansi ansi-fg-bright-cyan">27</span>                   ┊                         │<span class="ansi ansi-fg-bright-cyan">t</span><span class="ansi ansi-fg-bright-cyan">'</span>      ┊        │ Game Port ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━━━┛</code></pre>

Next we have a long, which is a four byte, or 32 bit integer.  It is `03 3c ad 93`, which comes out to `2477603843`.  This is the [SteamID](https://developer.valvesoftware.com/wiki/SteamID) of the server, which is a unique identifier for a Steam account, not for an app.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000090</span>│ <span class="ansi ansi-fg-bright-yellow">03</span> <span class="ansi ansi-fg-bright-cyan">3c</span> <span class="ansi ansi-fg-bright-yellow">ad</span> <span class="ansi ansi-fg-bright-yellow">93</span>             ┊                         │<span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-cyan"><</span><span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-yellow">x</span>    ┊        │ Steam ID ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━━┛</code></pre>

Following that we have a string of "keywords".  These are tags for describing the server.  The string we get is comma separated, and rather opaque.  Some things that stand out are `shard001` and what looks like a time, `14:09`.  A hive in DayZ is a group of servers that a character can migrate in between, and the player state is stored on a shard. Beyond that I'm not sure what any of the rest of it is.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000094</span>│ <span class="ansi ansi-fg-bright-yellow">b4</span> <span class="ansi ansi-fg-bright-cyan">62</span> <span class="ansi ansi-fg-bright-cyan">40</span> <span class="ansi ansi-fg-bright-yellow">01</span> <span class="ansi ansi-fg-bright-cyan">62</span> <span class="ansi ansi-fg-bright-cyan">61</span> <span class="ansi ansi-fg-bright-cyan">74</span> <span class="ansi ansi-fg-bright-cyan">74</span> ┊ <span class="ansi ansi-fg-bright-cyan">6c</span> <span class="ansi ansi-fg-bright-cyan">65</span> <span class="ansi ansi-fg-bright-cyan">79</span> <span class="ansi ansi-fg-bright-cyan">65</span> <span class="ansi ansi-fg-bright-cyan">2c</span> <span class="ansi ansi-fg-bright-cyan">6e</span> <span class="ansi ansi-fg-bright-cyan">6f</span> <span class="ansi ansi-fg-bright-cyan">33</span> │<span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-cyan">b</span><span class="ansi ansi-fg-bright-cyan">@</span><span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-cyan">b</span><span class="ansi ansi-fg-bright-cyan">a</span><span class="ansi ansi-fg-bright-cyan">t</span><span class="ansi ansi-fg-bright-cyan">t</span>┊<span class="ansi ansi-fg-bright-cyan">l</span><span class="ansi ansi-fg-bright-cyan">e</span><span class="ansi ansi-fg-bright-cyan">y</span><span class="ansi ansi-fg-bright-cyan">e</span><span class="ansi ansi-fg-bright-cyan">,</span><span class="ansi ansi-fg-bright-cyan">n</span><span class="ansi ansi-fg-bright-cyan">o</span><span class="ansi ansi-fg-bright-cyan">3</span>│ Keywords ┃
┃<span class="ansi ansi-fg-bright-black">00000110</span>│ <span class="ansi ansi-fg-bright-cyan">72</span> <span class="ansi ansi-fg-bright-cyan">64</span> <span class="ansi ansi-fg-bright-cyan">2c</span> <span class="ansi ansi-fg-bright-cyan">73</span> <span class="ansi ansi-fg-bright-cyan">68</span> <span class="ansi ansi-fg-bright-cyan">61</span> <span class="ansi ansi-fg-bright-cyan">72</span> <span class="ansi ansi-fg-bright-cyan">64</span> ┊ <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">31</span> <span class="ansi ansi-fg-bright-cyan">2c</span> <span class="ansi ansi-fg-bright-cyan">6c</span> <span class="ansi ansi-fg-bright-cyan">71</span> <span class="ansi ansi-fg-bright-cyan">73</span> <span class="ansi ansi-fg-bright-cyan">30</span> │<span class="ansi ansi-fg-bright-cyan">r</span><span class="ansi ansi-fg-bright-cyan">d</span><span class="ansi ansi-fg-bright-cyan">,</span><span class="ansi ansi-fg-bright-cyan">s</span><span class="ansi ansi-fg-bright-cyan">h</span><span class="ansi ansi-fg-bright-cyan">a</span><span class="ansi ansi-fg-bright-cyan">r</span><span class="ansi ansi-fg-bright-cyan">d</span>┊<span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">1</span><span class="ansi ansi-fg-bright-cyan">,</span><span class="ansi ansi-fg-bright-cyan">l</span><span class="ansi ansi-fg-bright-cyan">q</span><span class="ansi ansi-fg-bright-cyan">s</span><span class="ansi ansi-fg-bright-cyan">0</span>│          ┃
┃<span class="ansi ansi-fg-bright-black">00000126</span>│ <span class="ansi ansi-fg-bright-cyan">2c</span> <span class="ansi ansi-fg-bright-cyan">65</span> <span class="ansi ansi-fg-bright-cyan">74</span> <span class="ansi ansi-fg-bright-cyan">6d</span> <span class="ansi ansi-fg-bright-cyan">34</span> <span class="ansi ansi-fg-bright-cyan">2e</span> <span class="ansi ansi-fg-bright-cyan">32</span> <span class="ansi ansi-fg-bright-cyan">30</span> ┊ <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">2c</span> <span class="ansi ansi-fg-bright-cyan">65</span> <span class="ansi ansi-fg-bright-cyan">6e</span> <span class="ansi ansi-fg-bright-cyan">74</span> │<span class="ansi ansi-fg-bright-cyan">,</span><span class="ansi ansi-fg-bright-cyan">e</span><span class="ansi ansi-fg-bright-cyan">t</span><span class="ansi ansi-fg-bright-cyan">m</span><span class="ansi ansi-fg-bright-cyan">4</span><span class="ansi ansi-fg-bright-cyan">.</span><span class="ansi ansi-fg-bright-cyan">2</span><span class="ansi ansi-fg-bright-cyan">0</span>┊<span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">,</span><span class="ansi ansi-fg-bright-cyan">e</span><span class="ansi ansi-fg-bright-cyan">n</span><span class="ansi ansi-fg-bright-cyan">t</span>│          ┃
┃<span class="ansi ansi-fg-bright-black">00000142</span>│ <span class="ansi ansi-fg-bright-cyan">6d</span> <span class="ansi ansi-fg-bright-cyan">34</span> <span class="ansi ansi-fg-bright-cyan">2e</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">30</span> ┊ <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">2c</span> <span class="ansi ansi-fg-bright-cyan">31</span> <span class="ansi ansi-fg-bright-cyan">34</span> <span class="ansi ansi-fg-bright-cyan">3a</span> <span class="ansi ansi-fg-bright-cyan">30</span> <span class="ansi ansi-fg-bright-cyan">39</span> <span class="ansi ansi-fg-bright-black">00</span> │<span class="ansi ansi-fg-bright-cyan">m</span><span class="ansi ansi-fg-bright-cyan">4</span><span class="ansi ansi-fg-bright-cyan">.</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">0</span>┊<span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">,</span><span class="ansi ansi-fg-bright-cyan">1</span><span class="ansi ansi-fg-bright-cyan">4</span><span class="ansi ansi-fg-bright-cyan">:</span><span class="ansi ansi-fg-bright-cyan">0</span><span class="ansi ansi-fg-bright-cyan">9</span><span class="ansi ansi-fg-bright-black">⋄</span>│          ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━━┛
</code></pre>

Lastly, we have Game ID, which is a 32-bit integer, representing the Appd ID.  App ID's apparently were originally 16 bits, as mentioned above.  Since the number of apps grew too large, they moved to 32 bits, but couldn't alter the protocol.  So when Game ID is present in the EDF, it supercedes the value in the App ID field.  According to the docs, the App ID is in the lower 24 bits of this field, so just `03 5F AC`.

I'm not sure why they chose to just use three bytes for App ID, since we're sending all four bytes anyway. It is fairly safe to assume that this field will have 0's in the first byte, so you can cast it as a long int, and the whole thing will be valid.  The value we get is `00000158ac5f03`, which is `221100` in decimal.  This is the App ID for DayZ.

<pre><code class="ansi">┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━┓
┃<span class="ansi ansi-fg-bright-black">00000158</span>│ <span class="ansi ansi-fg-bright-yellow">ac</span> <span class="ansi ansi-fg-bright-cyan">5f</span> <span class="ansi ansi-fg-bright-yellow">03</span> <span class="ansi ansi-fg-bright-black">00</span>             ┊                         │<span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-cyan">_</span><span class="ansi ansi-fg-bright-yellow">x</span><span class="ansi ansi-fg-bright-black">⋄</span>    ┊        │ Game ID ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━┛</code></pre>

# Next Steps

That is it for A2S_INFO. There's a lot in there, but not everything we see in tools like DayZSA.  This post has gotten long, so I'm going to stop here, but in the next post I'll cover `A2S_RULES` which has a lot more data crammed into it.  Eventually I may get into server discovery and CM Client, but we will see.
