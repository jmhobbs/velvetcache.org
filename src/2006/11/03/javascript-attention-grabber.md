---
category:
- Geek
creator: admin
date: 2006-11-03
layout: layout.njk
permalink: /2006/11/03/javascript-attention-grabber/
tags:
- Internet
- JavaScript
- Programming
- Work
title: Javascript Attention Grabber
type: post
wp_id: "87"
---

I honestly don't know what else to call this, so "Javascript Attention Grabber" will have to do.  We have a site at UNO that has an alphabetic list and you can click on a letter at the top of the page to jump to it.  However, the pages aren't always long enough for the jump to happen, yet are still cluttered enough to get lost without it.  Thus I whipped up this simple little guy to highlight the div for a short time.

I hardcoded in the return background color because `tempColor = tempObj.style.backgroundColor;` wasn't pulling off the old color.  I'll have to figure that one out later.  I'd also like to add a fader, but that was too much work at the time.  ANyway, here's the first version.

```javascript
  function hilighter (targetID,color) {
    targetObj = document.getElementById(targetID);
    tempColor = "#F6F6F6";
    targetObj.style.backgroundColor = color;
    setTimeout('unhiliter("'+targetID+'","'+tempColor+'")',1000);
  }

  function unhiliter (targetID,color) {
    document.getElementById(targetID).style.backgroundColor = color;
  }
```
```html
<a href="#anchorName">Link</a>
```

P.S. Yet again I'm faced with the ugly usability problems of this website.  I need to take the time to fix it up.  Let me go make a note of that on my [Nokia 770](/2006/11/03/nokia-770-2/). :)
