---
category:
- Geek
creator: admin
date: 2006-10-01
permalink: /2006/09/30/drumkit-v001/
tags:
- DrumKit
- Equipment
- Java
- Music
- Programming
- Projects
title: DrumKit v0.01
type: post
wp_id: "58"
---

I finally started work on my DrumKit project today.  Essentially this is an attempt to create a electronic drum set from some pads, an old keyboard and a computer.  I got down to buisness and wrote up a Java app that plays back the drum sounds on key press events. It filters them and even handles shift for the hi-hat open/close, which I'll use to make the hi-hat pedal work.

The source is messy, and still has relative path's for the sound files.  Essentially it's a frame, a keylistener and some clips. Not too rough, though patching this together from the Java API and sparse information on javax.sound.sampled was tougher than I guessed it would be.  Anyway, it's got some bugs and features not implemented, but I can play drums with my keyboard now, which is a good start.

<figure>
  <img src="https://static.velvetcache.org/projects/drumkit/drumkit001.jpg" alt="DrumKit v0.01" />
  <figcaption>The Cutting Edge DrumKit GUI</figcaption>
</figure>

I suppose I should upload the code just in case, and to keep track of my versions.  You can get the source and the sounds for Version 0.01 below.  As a heads up, all the kit samples are absolute path'd for my machine, so you'll need to adjust them if you intend to compile it.

- [ DrumKit001.tar.gz](https://static.velvetcache.org/projects/drumkit/DrumKit001.tar.gz) - Everything
- [ DrumKit.java](https://static.velvetcache.org/projects/drumkit/v001Source/DrumKit.java) - Main class
- [ KeyHandler.java](https://static.velvetcache.org/projects/drumkit/v001Source/KeyHandler.java) - Event listener
- [ KitClip.java](https://static.velvetcache.org/projects/drumkit/v001Source/KitClip.java) - Kit sample class
