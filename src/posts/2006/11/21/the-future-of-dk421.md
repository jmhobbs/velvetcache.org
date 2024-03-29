---
category:
- Geek
creator: admin
date: 2006-11-21
permalink: /2006/11/21/the-future-of-dk421/
tags:
- DK421
- DrumKit
- Projects
title: The Future Of DK421
type: post
wp_id: "97"
---

I'm currently re-thinking the future of [DK421](/category/projects/dk421/).  It's the hardware that has become the hang up. I'm 100% sure that the design I have for the pads is physically impossible.   In point of fact, the design is mentally obtuse to begin with. A drum stick just doesn't hit down with any reasonable amount of force. To think it could be enough to move a not insubstantial amount of plywood any real distance was silly. More so considering I also sought to compress four springs.

Where that leaves me is with a strong desire to use piezo transducers, or some other cheap, microphone-like option, combined with either a micro controller or some relays or something.  The hard part is I don't know any electronics really, I'm just winging it so far.

The vector I'd be most happy with pursuing is using an [Ardunio](http://www.arduino.cc/) board to read the transducers (voltage? wattage? I dunno...) and when they peak a certain amount, fire off some information over a serial connection.  The problems are many for this.  One, I don't have a micro controller, and don't know how to program and use them.  Two, I've never done any serial port programming, and am not completely sure if it's even a viable option for passing the data. Three, I don't know what to measure on the transducers or how to do so. Four, I think Java might be out of the question with the serial.  There is a javax.comm, but I'm not so sure about it.

Regardless, I want to try this way.  I've wanted to buy a micro controller to play with for a while now, and the Ardunio is open-source and way cheaper than a basic stamp setup. I also kinda want to get me C++ back, and maybe learn how to use QT or GTK.  This is either going to be long and fun, or die in shame. Going to buy the board tonight I think.
