---
category:
- Geek
creator: admin
date: 2006-08-30
layout: layout.njk
permalink: /2006/08/29/lastfm-web-services/
tags:
- PHP
- XML
title: Last.fm Web Services
type: post
wp_id: "17"
---
Huzzah!  I was looking around the [audioscrobbler.net web services](https://www.audioscrobbler.net/data/webservices/) and found my holy grail!

I've tried a bunch of ways to retrieve a cover for the most recent track.Â  I've never been able to figure out Amazon's image naming system, though I found neat hacks at this cool site: [https://aaugh.com/imageabuse.html](https://aaugh.com/imageabuse.html) that absorbed way too much time.Â  I've also tried retrieving Google results and regex'ing the _whole_ file, which was slow and rarely worked.

Well, today I found the artist info XML at audioscrobbler. [Example](https://ws.audioscrobbler.com/1.0/album/Metallica/Metallica/info.xml)
I'm going to incorporate this bad boy into my last.fm scripts and put the system up here on the homepage.Â  Maybe I'll try my hand at a WP plugin, just for the heck of it.
