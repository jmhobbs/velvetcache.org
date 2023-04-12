---
category:
- Geek
creator: admin
date: 2006-11-22
permalink: /2006/11/22/javascript-colorshifter/
tags:
- JavaScript
- Programming
- Snippets
title: Javascript Colorshifter
type: post
wp_id: "98"
---
<script type="text/javascript" src="https://static.velvetcache.org/pages/2006/11/22/javascript-colorshifter/cs.js"></script>

I wrote this little function because I wanted a little background color shifter without having to deal with a big script animation library.  Really simple, and it takes some time to run through itself.  I'll rewrite it to be faster, more flexible later.

Actually I probably won't.

Anyway, here it is, and heres a demo.  Multiple clicks will get you psycho function timeouts all over the place and some strange color changes.  Rad.

<input type="button" value="Click Me For Color-Shift" onclick="colorShift('csTarget',255,0,0,255,255,255);" />

```javascript
function colorShift(strTarget,sR,sG,sB,tR,tG,tB) {
  if(sR == tR && sG == tG && sB == tB) { return; }
  if(sR != tR) { (sR < tR) ? sR++ : sR--; }
  if(sG != tG) { (sG < tG) ? sG++ : sG--; }
  if(sB != tB) { (sB < tB) ? sB++ : sB--; }
  document.getElementById(strTarget).style.backgroundColor = "rgb("+sR+","+sG+","+sB+")";
  setTimeout("colorShift('"+strTarget+"',"+sR+","+sG+","+sB+","+tR+","+tG+","+tB+")",1);
}
```

Now that I look, it's sorta related to the [Javascript Attention Grabber](/2006/11/03/javascript-attention-grabber/) from earlier this month.
