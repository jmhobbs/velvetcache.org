---
date: 2023-04-05
layout: layout.njk
tags:
- go
- golang
- WordPress
- 11ty
title: Moving from WordPress to 11ty
type: post
summary: >
  After 17 years with WordPress, it's time for something different. Here's how
  I migrated 500 posts to a static site generator.
---

In April 2006, I set up the first version of my blog.  It was written from scratch in PHP, with minimal admin tooling and few features.  By July, I had moved over to WordPress.  It was great!  Writing new content was a breeze, modifying my theme or adding functionality was no work at all.

Since then, I've written plugins, made themes as a job, and worked on the infrastructure of (what was) the best WordPress host around, Flywheel.
But WordPress was wearing on me.  I didn't blog as much anymore, and when I did there was always updates to run before I got started, or a plugin which had stopped working and needed to be replaced.  My posts were full of old, stale shortcodes that didn't work.  The blog became toil, instead of an outlet.

I needed something new and simple.  Static site generators have been around for ages, I've used [Hugo](https://gohugo.io/), [Jekyll](https://jekyllrb.com/) and others, but I really enjoy [11ty](https://www.11ty.dev/).  It was obvious that is what I should move to, and it wouldn't lock me in too bad if I ever decided to go to another setup.

## wp-to-11ty

All that remained was to get the content out of WordPress and into 11ty.  When I first got this itch in 2019, there wasn't really much around, so I started writing it.  WordPress has an [XML export format](https://wordpress.com/support/export/), which, while thorough, is a mess.  But it's what we have, so I wrote some tooling in Go to read through the XML soup and write out files for each post and page.  It worked, but there was maybe 20% left to do.  Life happened, I got distracted, and it lingered until last month.  It's now "finished" enough that I felt comfortable exporting out and starting to build up the new site on 11ty.

That tooling is [wp-to-11ty](https://github.com/jmhobbs/wp-to-11ty) and it works fairly well.  There are a number of hoops to jump through moving tags, categories, etc. from WordPress to 11ty, and I made some comprimises to handle it.  But it does work, and the blog you are reading was exported using that tool.

## Crawl and Compare

After the export, I wanted to get _some_ confidence I hadn't missed anything obvious.  To check that, I used the [colly package](https://github.com/gocolly/colly) to crawl my locally hosted 11ty site, and the existing WordPress site on velvetcache.org. It just recorded every URL it visited, which I dropped into a file.

```go
package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains("velvetcache.org", "www.velvetcache.org"))

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.Path)
	})

	c.Visit("https://velvetcache.org/")
}
```

Some additonal munging on the command line  and I had two files containing discovered paths that I could compare.  It was pretty comparable, and anything I missed wasn't something I was worried about keeping.

## Clean Up

My blog posts are...messy.  Over the last ~15 years I've used a couple different syntax highlighting plugins, I've changed my design, and I've embedded a _lot_ of custom styles in my post HTML.  On top of that, I've had my WordPress hacked a couple times, and I didn't manage to clean it all the way out.

I needed a clean, and simple new format.  I'm a big fan of [Markdown](https://www.markdownguide.org/) (who isn't?) so that's the obvious choice.  I had (and still have) a lot of converting to do, but I'm working through it a post at a time.

I have a lot of quick fixes that I apply with regex substitutions, then I go through and manually clean up what it missed (which is often a lot).  I considered using something more complex like [turndown](https://github.com/mixmark-io/turndown) but I want to review all my pages manually anyway and unify the styles, so this solution is good enough.

```bash
#!/bin/env bash

cat $1 |
  perl -pe 's|http://static.velvetcache.org|https://static.velvetcache.org|g' |
  perl -pe 's|https?://(www.)?velvetcache.org/|/|g' |
  perl -pe 's|<tt>(.*?)</tt>|`\1`|g' |
  perl -pe 's|<h1>(.*?)</h1>|# \1\r|g' |
  perl -pe 's|<h2>(.*?)</h2>|## \1\r|g' |
  perl -pe 's|<h3>(.*?)</h3>|### \1\r|g' |
  perl -pe 's|<h4>(.*?)</h4>|#### \1\r|g' |
  perl -pe 's|<h5>(.*?)</h5>|##### \1\r|g' |
  perl -pe 's|<strong>(.*?)</strong>|**\1**|g' |
  perl -pe 's|<b>(.*?)</b>|**\1**|g' |
  perl -pe 's|<em>(.*?)</em>|_\1_|g' |
  perl -pe 's|<i>(.*?)</i>|_\1_|g' |
  perl -pe 's|<p.*?>(.*?)</p>|\n\1\n|g' |
  perl -pe 's|<img src="(.*?)".*?alt="(.*?)".*?>|![\2](\1)|g' |
  perl -pe 's|<a href="(.*?)".*?>(.*?)</a>|[\2](\1)|g' > ${1/.html/.md}
```

Oddly, I had to use `cat` instead of `< $1`. When piping perl together it sticks and doesn't finish unless I use `cat`.  I assume perl must be waiting on something from `stdin` it doesn't get with redirection.

## Cloudflare Pages

I've used Netlify and Cloudflare together [in the past](/2020/01/29/netlify-cloudflare-crazy-delicious/) and it's worked out nicely.  On my last static site I just used [Cloudflare Pages](https://pages.cloudflare.com/).  It is not as slick and fancy as Netlify, but it is fast and simple.  I already run my DNS on Cloudflare, so to skip the hassle I just went with Cloudflare.

## WebFinger

My blog had very few "live" features that I wanted to ensure I brought over to the new system.  The one I decided to mess with was [WebFinger](https://webfinger.net/).  WebFinger is "...discover information about people or other entities on the Internet that are identified by a URI...".  My use case is to have a reference I can share for Mastadon, which can redirect if I ever move servers.  So if you [look up `jmhobbs@velvetcache.org` with WebFinger](https://webfinger.net/lookup/?resource=jmhobbs%40velvetcache.org), you'll get my current Mastadon server information at [noc.social](https://noc.social/@jmhobbs).

Due to how WebFinger works, a static page isn't really up to the task.  Luckily, Cloudflare has functions for their pages sites, which works nicely.

```javascript
// functions/.well-known/webfinger.js
const allowed = [
  'acct:john@velvetcache.org',
  'acct:jmhobbs@velvetcache.org',
];

export function onRequest(context) {
  const url = new URL(context.request.url);
  console.log(url.searchParams.get('resource'));
  if(allowed.includes(url.searchParams.get('resource'))) {
    return Response.redirect('https://noc.social/.well-known/webfinger?resource=acct:jmhobbs@noc.social', 302);
  }
  return new Response('400 Bad Request', { status: 400, statusText: 'Bad Request' });
}
```

## ¯\\\_(シ)_/¯

That's pretty much it.  I've got a lot of cleaning up posts to do, and the CSS styles need some attention, but overall it is all working neatly.  I'm happy with it, and if you're currently a light user of WordPress who's tired of the hustle, you might be happy with it too.
