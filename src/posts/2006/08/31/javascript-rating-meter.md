---
category:
- Geek
creator: admin
date: 2006-08-31
permalink: /2006/08/31/javascript-rating-meter/
tags:
- CSS
- JavaScript
title: Javascript Rating Meter
type: post
wp_id: "18"
---

> **Update:** A newer version is outlined in [this](/2006/09/01/javascript-rating-meter-2/) post.

Due to a discussion I had earlier today, I was thinking about how to make one of those little star meters where you can click to set the rating and they change colors.  I took some time to code up this quick version.

I didn't want to use image swapping, because that wouldn't be as flexible.  So instead, you change the background color and have a transparent cutout of the image you want over that.  It'll make more sense once I get an example made.

For now, here's my prototype.  It needs a lot of work to make it into a good script, but it's a start.  [Example Here](https://static.velvetcache.org/pages/2006/08/31/javascript-rating-meter/)

```javascript
var sets = Array();

function hover(level) {
  for(j = 1; j <= level; j++)
  {
    document.getElementById('S'+j).style.backgroundColor = "#0FF";
  }
  for(k = (level+1); k <= 5; k++)
  {
    document.getElementById('S'+k).style.backgroundColor = "transparent";
  }
}

function unhover(group) {
  for(j = 1; j < 6; j++)
  {
    if(sets[group] == null)
    {
      document.getElementById('S'+j).style.backgroundColor = "transparent";
    }
    else
    {
      if(sets[group] >= j)
      {
        document.getElementById('S'+j).style.backgroundColor = "#F00";
      }
      else
      {
        document.getElementById('S'+j).style.backgroundColor = "transparent";
      }
    }
  }
}

function set(group,level) {
  sets[group] = level;
}
```
