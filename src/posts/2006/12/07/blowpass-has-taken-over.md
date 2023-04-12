---
category:
- Geek
creator: admin
date: 2006-12-07
permalink: /2006/12/07/blowpass-has-taken-over/
tags:
- BlowPass
- Internet
- JavaScript
- Open Source
- PHP
- Programming
- Projects
- Security
title: BlowPass Has Taken Over
type: post
wp_id: "106"
---

So I [thought that I had moved past BlowPass](/2006/12/04/passletcom/).  I guess I was wrong.  I've been spending every spare moment working on it.  I found what I feel is a better Blowfish library at [www.farfarfar.com](http://www.farfarfar.com/scripts/encrypt/).  I still can't implement any of the vector tests because they're all in hex and translate into nasty characters.  This means I have no actual idea if the crypt is working. I also quickly stopped trying to write my own Twofish implementation.  I could handle it in C I think, but not JavaScript, I don't know enough of it and it's little oddities.

Regardless of all that, I've got my prototype AJAX down pat now (okay, AHAH) and I'm working up my own open source version of [passlet.com](http://www.passlet.com/).  Here's a nice list of features/todo's.

- Uses a non-proprietary algorithm (Blowfish)
- Has AJAX-y-ness
- Uses PHP
- Uses a database abstraction library TODO
- Slick animations (mootools?) TODO

You can check out the current version at [https://static.velvetcache.org/projects/blowpass/demo/](https://static.velvetcache.org/projects/blowpass/demo/) to play around.
