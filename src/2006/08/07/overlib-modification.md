---
category:
- Geek
creator: admin
date: 2006-08-07
layout: layout.njk
permalink: /2006/08/07/overlib-modification/
tags:
- JavaScript
title: OverLIB Modification
type: post
wp_id: "8"
---

This is a post from a previous system.  Information, links and images may not be vaild.

There are a number of scripts I use on a semi-regular basis, but I hate having to do a bunch (or even a little) setup for them.  I prefer the lightbox way, load the script, add some attributes to tags and you're good to go.

In that spirit I recently made a little adjustment to Erik Bosrups [OverLIB](http://www.bosrup.com/web/overlib/), which adds nice little tooltips  to your page. Basically I just moved all the init stuff into the overLIB file. Here's the code:

```javascript
  addLoadEvent(
   function () {
     var objBody = document.getElementsByTagName("body").item(0);
     var objTiplayer = document.createElement("div");
     objTiplayer.setAttribute('id','overDiv');
     objTiplayer.style.visibility = 'hidden';
     objTiplayer.style.position = 'absolute';
     objTiplayer.style.zIndex = '1000';
     objBody.insertBefore(objTiplayer, objBody.firstChild);
     TText = new Array(0);
     if (!document.getElementsByTagName){ return; }
     var anchors = document.getElementsByTagName("a");
     for (var i=0; i < anchors.length; i++) {
       var anchor = anchors[i];
       if(anchor.getAttribute("rel") == "tooltip") {
        var titleC = anchor.getAttribute("data");
        var content = titleC.split("|");
        anchor.c1 = content[0];
        anchor.c2 = content[1];
        if(anchor.c1!=undefined) {
          anchor.onmouseover = function () { return overlib(this.c2, CAPTION,     this.c1); }
        }
        else
        {
          anchor.onmouseover = function () {return overlib(this.c2);}
        }
        anchor.onmouseout = function() { return nd(); }
      }
    }
  });
```

Once that's in there you just add two attributes to your tags, `rel="tooltip"` and `data="Title|Content"`. To get a title-free tooltip you put `data="|Content"`.  You can also still initialize the tooltips the old way for more advanced features

I realize this is not the best solution, as it creates invalid markup and can only do two types of tooltip (at this point) but it comes in handy and it's nice to have a simple drop-in javascript file to do all this.

[Live Example](https://static.velvetcache.org/projects/js/overlib_mod/)

[OverLIB.js with the mod installed](https://static.velvetcache.org/projects/js/overlib_mod/overlibMod.js)
