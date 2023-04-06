---
category:
- Geek
creator: admin
date: 2006-10-16
layout: layout.njk
permalink: /2006/10/16/simplebooks-v001/
tags:
- PHP
- Projects
- SimpleBooks
title: SimpleBooks V0.01
type: post
wp_id: "72"
---

  Like essentially every project I write SimpleBooks arose out of a need for a simple web based book reader. Late in the summer of 2006 I wrote a set of scripts that worked as a converter for plain text files to create on-line readable e-books. I also created the viewer.  At that time the system was very simple, as I just bookmarked my page when I was done reading and came back to it on any computer using Foxmarks.

  As of 10/15/06 I've overhalued the system and made it more robust. It now handles book uploads as well as covers and keeps track of book details using flatfiles.  It still needs a user system and built in bookmarking, but it's way better than it used to be.

 This initial version has some bugs floating around and is certainly not terribly attractive or user friendly (except the viewer).  The splitting method isn't very intelligent and I would like to refine it to handle HTML tags someday, but it's not terribly high on the list. Also, there might be some exploitable code in the creator.php, though I think I covered all of that stuff up (cross fingers).  There are numerous caveats, such as encoding types of the txt files and the like, but they are easily noticed and fixed.

Anyway, you can try the viewer with a copy of _Crime And Punishment_ from [Project Gutenberg](https://www.gutenberg.org/) or download the whole shebang and mess with it yourself.

**Links**

- [Demo](https://static.velvetcache.org/projects/simplebooks/demo/viewer.php?book=cap&page=1)
- [Source (tar.gz)](https://static.velvetcache.org/projects/simplebooks/simplebooks_v001.tar.gz)
- [form.php Source](https://static.velvetcache.org/projects/simplebooks/V001/form.phps)
- [creator.php Source](https://static.velvetcache.org/projects/simplebooks/V001/creator.phps)
- [viewer.php Source](https://static.velvetcache.org/projects/simplebooks/V001/viewer.phps)
