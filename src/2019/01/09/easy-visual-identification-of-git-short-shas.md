---
category:
- Geek
creator: admin
date: 2019-01-09
layout: layout.njk
tags:
- git
- JavaScript
title: Easy visual identification of git short sha's.
type: post
permalink: /2019/01/09/easy-visual-identification-of-git-short-shas/
wp_id: "2845"
summary: A list of short SHA's can be hard to read, wrapping them in color can help.
---

For a recent pet project at work I had to display a bunch of git short shas.  Most of the time, these shas should match each other, and it is important to be able to quickly glance at them and evaluate if any of them are not the same.

Sure, you could count on your eyes to just notice the characters don't match, but short shas are drawing from a limited alphabet (16 characters) and we would only be adding more shas to the listing over time, so noticing one abberant item would get harder over time.

![an example of colorized shas](http://static.velvetcache.org/pages/2019/01/09/easy-visual-identification-of-git-short-shas/example.png)

The solution I landed on was to drop the final character of the sha and use that as a hex color string.  While this is imperfect, it works remarkably well.

To make the sha readable over the color I needed to find a contrasting color.  To keep it simple, I took the average of the individual channels, and if it was over 128 I used black.  Under 128 I used white.

Here's a little snippet of the JavaScript I wrote for this:

```javascript
$span.css({
  backgroundColor: '#' + sha.substring(0,6),
  color: contrastColorForSha(sha)
});

function contrastColorForSha(sha) {
  let avg = ( parseInt(sha.substring(0,2), 16) + 
              parseInt(sha.substring(2,4), 16) +
              parseInt(sha.substring(4,6), 16)
            ) / 3;
  return (avg > 128) ? '#000' : '#FFF';
}
```

This could be improved by using luminosity measures designed for eyeballs instead of a rough mean of the colors, I could constrain the colorspace for shas to throw out colors that are problematic for color blindness, etc.  I'll leav e that as an exercise for the reader.
