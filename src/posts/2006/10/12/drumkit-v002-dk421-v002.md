---
category:
- Geek
creator: admin
date: 2006-10-12
permalink: /2006/10/12/drumkit-v002-dk421-v002/
tags:
- DK421
- DrumKit
- Java
- Programming
- Projects
title: DrumKit V0.02 -> DK421 V0.02
type: post
wp_id: "67"
---

I've been working on my [DrumKit](/2006/09/30/drumkit-v001/) again, since I [don't get to play with my Nokia 770](/2006/10/11/nokia-770-review/). I came up with a better name for this revision, DK421.  It doesn't mean anything.  Well, it kinda does.  It came from a Stormtrooper in ANH ("TK421, Why aren't you at your post? TK421, respond.")

There isn't a lot of change in this version.  I re-read the API for javax.sound.sampled and re-wrote the `SoundClip` class (formerly `KitClip`) from scratch.  It think it's cleaner and more usefull for other projects now.

I also added in a file chooser class, `DrumChooser` so there aren't any absolute paths in the source anymore.  I also added to that class a minor key-fetching function.  I hope to make the kit's expandable, as many drums as you can until you hit a `LineUnavailableException`.

The GUI hasn't changed a lick, and I need to clean it up and develop a new layout for it and for the "multiple-drum" expandable version.  It's surprising what you can learn from an API.
