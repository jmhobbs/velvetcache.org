---
category:
- Geek
creator: admin
date: 2006-08-29
layout: layout.njk
permalink: /2006/08/28/more-wordpress-stuff/
tags:
- PHP
- Themes
- Wordpress
title: More WordPress Stuff
type: post
wp_id: "14"
---
So I've gotten all the posts over from the old system.  I've had alot of trobule with getting [iG:Syntax Hiliter](http://blog.igeek.info/wp-plugins/igsyntax-hiliter/) to behave, though it's mostly my own fault for not understanding the options and the styling.  The "plain text" links still weren't working, so I took them off.

I've been finding little flaws in the theme here and there, and just today created a home.php to override the index and give this nice custom job, more like the old page.  I also had problems getting [scrobbler](http://leflo.de/projekte/wordpress/scrobbler) to work.  I finally realized it wasn't scrobbler that was messed up, it was [wp-cache](http://mnm.uib.es/gallir/wp-cache-2/) making the dynamic [last.fm](http://www.last.fm/) content into cached copies.

I figured that out after pulling scrobbler and inserting my own last.fm code straight into the template.  Not the best solution, but I never intend on releasing this template as public anyway.

Still chugging away at the static site, should have some of my content back soon.
