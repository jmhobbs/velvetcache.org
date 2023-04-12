---
category:
- Geek
creator: admin
date: 2006-11-04
permalink: /2006/11/03/more-nokia-goodness/
tags:
- Hardware
- Nokia 770
title: More Nokia Goodness
type: post
wp_id: "88"
---

I've been messing with my 770 some more.  Such a sweet device.  Here's a quick rundown of what I've loaded up so far:

- Leafpad
- GFrotz
- Maemodrac (Solitaire)
- LXDoom (And Doom wad)
- Maemosweeper
- FBReader (E-Books)
- Puchi (Themer)
- Maemo Mapper (Map Reader)
- osso X-term (Terminal Emu)
- Tuner (Musical)
- Gaim
- GPE Calendar
- GPE Contacts
- GPE Todo
- MaemoPad
- Media Streamer
- MPlayer
- Kismet
- dsniff

I've also got the stuff to mount SMB/CIFS shares and wrote two small mount/unmount scripts for my home network.  It all worked nicely until I tried playing a few MP3's over the CIFS mount. It didn't buffer very intelligently, I watched my network monitor on the serving machine and it would spike every time the song "skipped" on the Nokia.

Looking for an answer I hopped onto Synaptic and found a UPnP server, [GMediaServer](http://www.gnu.org/software/gmediaserver/) since I knew there was a nice UPnP audio player for the Nokia, from the company itself actually.  Hooked that up and it's behaved so far.  Don't know what I'm going to do about streaming out video to it though.  Oh well, I was going to buy a bigger rs-mmc card anyway, I can just dump videos onto that if I need to.

Here's some screens of the UPnP media player and an XTerm.  I greened up the layout with puchi.

[![Media Streamer](https://static.velvetcache.org/pages/2006/11/03/more-nokia-goodness/screenshot00.png)](https://static.velvetcache.org/pages/2006/11/03/more-nokia-goodness/screenshot00.png)
[![Home](https://static.velvetcache.org/pages/2006/11/03/more-nokia-goodness/screenshot01.png)](https://static.velvetcache.org/pages/2006/11/03/more-nokia-goodness/screenshot01.png)
