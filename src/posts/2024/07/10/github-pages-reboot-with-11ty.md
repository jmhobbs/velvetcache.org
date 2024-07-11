---
date: 2024-07-10
tags:
- nodejs
- Github
title: GitHub Pages Reboot With 11ty 
type: post
permalink: /2024/07/10/github-pages-reboot-with-11ty/
opengraph_image: /static/og/2024-07-10-github-pages-reboot-with-11ty.png
summary: I've had my repos listed on my GitHub pages for 14 years, it was time to revisit it.
---

I've had a list of my most recently updated repositories on my GitHub user page long enough that the repo is `jmhobbs.github.com` instead of `jmhobbs.github.io`.  Fourteen years ago I wrote a little Python script to [generate this page every day](/2010/06/29/auto-generated-github-user-page-with-py-github/)  It's worked well for all that time, with minor updates here and there. Recently I got a dependabot notice that some packages were out of date, so I took another look at it.  The page was pretty dated looking, so I decided to refresh it. After that was done, I decided to port it to [11ty](https://www.11ty.dev/), because I love 11ty and this is a chance to play with eleventy-fetch which I haven't tried yet.

# Global Data Files

The 11ty cascade allows for a global data files, which can be in several formats, including executable JavaScript.  To start, I added two files, a `config.json` with my username in it, and `user.js` which would fetch my user details from GitHub via the REST API endpoint for [getting a user](https://docs.github.com/en/rest/users/users?apiVersion=2022-11-28#get-a-user)

```javascript
// _data/user.js
const fetch = require('node-fetch');
const config = require('./config.json');

module.exports = async () => {
  const response = await fetch(`https://api.github.com/users/${config.username}`);
  return response.json();
};
```

```javascript
// _data/config.json
{
  "username": "jmhobbs"
}
```

With that done, I can now use my GitHub user API response in my template:

{% raw %}
```html
{# index.njk #}
<h1>@{{ user.login }}</h1>
```
{% endraw %}

![I am @jmhobbs](https://static.velvetcache.org/pages/2024/07/10/github-pages-reboot-with-11ty/i-am-jmhobbs.png)

That's great, but it's going to make that API request every time it needs to render the template, which is a bit agressive.  Which takes a whopping ~250ms each time!

![250ms, whoa is me!](https://static.velvetcache.org/pages/2024/07/10/github-pages-reboot-with-11ty/slow-users-fetch.png)

Luckily, there's an easy way to cache it.

# @11ty/eleventy-fetch

The [eleventy-fetch plugin](https://www.11ty.dev/docs/plugins/fetch/) exists to fix this problem.  Not only will it fetch the data for me (so long `node-fetch`! ok, so long direct dependency on `node-fetch`...), it will convert the response to JSON and cache that data for as long as I would like.

```javascript
// _data/user.js
const EleventyFetch = require('@11ty/eleventy-fetch');
const config = require('./config.json');

module.exports = async () => {
  return EleventyFetch(`https://api.github.com/users/${config.username}`, {
    duration: "7d",
    type: "json",
  });
}
```

I plug in the URL, set the cache life to 7 days, and it just works!

![Cache is king.](https://static.velvetcache.org/pages/2024/07/10/github-pages-reboot-with-11ty/eleventy-fetch-fast.png)

# Repositories

Next, I need access to the repositories.  Unfortunately, I have way more repositories than the GitHub API will send in one page, so I have to paginate through them.

Doubly unfortunate, GitHub uses the [`link` header for pagination](https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api?apiVersion=2022-11-28#using-link-headers), and eleventy-fetch does not expose reponse headers.

We can use [octokit.js](https://github.com/octokit/octokit.js) to fetch them by pages, so we will focus on that first.

```javascript
// _data/repos.js
const config = require('./config.json');

module.exports = async function () {
  const { Octokit } = await import('octokit');

  const octokit = new Octokit();

  const iterator = octokit.paginate.iterator(octokit.rest.repos.listForUser, {
    username: config.username,
    per_page: 50,
    sort: 'pushed',
  });

  const repos = [];

  for await (const page of iterator) {
    for( const repo of page.data) {
      if(repo.private) {
        continue
      }
      repos.push(repo);
    }
  }
  return repos;
};
```
{% raw %}
```html
{# index.njk #}
<h1>@{{ user.login }}</h1>

<ul>
  {% for repo in repos %}
  <li>{{ repo.name }}</li>
  {% endfor %}
</ul>
```
{% endraw %}

It works! 🎉

![Repos galore!](https://static.velvetcache.org/pages/2024/07/10/github-pages-reboot-with-11ty/repository-works.png)

But it's slow 😿

![3.5 seconds for this?!](https://static.velvetcache.org/pages/2024/07/10/github-pages-reboot-with-11ty/repository-very-slow.png)

# AssetCache

You didn't think 11ty would leave us hanging, did you?  It doesn't, eleventy-fetch has a way to access the cache portion all by itself, [AssetCache](https://www.11ty.dev/docs/plugins/fetch/#advanced-usage).

```javascript
// _data/repos.js
const { AssetCache } = require('@11ty/eleventy-fetch');
const config = require('./config.json');

module.exports = async function () {
  let asset = new AssetCache(`repos_${config.username}`);

  if (asset.isCacheValid('12h')) {
    return asset.getCachedValue();
  }
  const { Octokit } = await import('octokit');

  const octokit = new Octokit();

  const iterator = octokit.paginate.iterator(octokit.rest.repos.listForUser, {
    username: config.username,
    per_page: 50,
    sort: 'pushed',
  });

  const repos = [];

  for await (const page of iterator) {
    for( const repo of page.data) {
      if(repo.private) {
        continue
      }
      repos.push(repo);
    }
  }

  await asset.save(repos, 'json');

  return repos;
};
```

Much better!

![AssetCache to the rescue.](https://static.velvetcache.org/pages/2024/07/10/github-pages-reboot-with-11ty/asset-cache-to-the-rescue.png)

# Ship It!

That's it!  That's the guts of the replacement, everything else is just template work and CSS.

I've put all the code samples, as a working demo, online at [github.com/jmhobbs/11ty-fetch-demo](https://github.com/jmhobbs/11ty-fetch-demo).  Each commit matches up with a heading in this post, so you can follow along at home if you like.

My GitHub page is, of course, up to date and available at [jmhobbs.github.io](https://jmhobbs.github.io/), and you can ge the source for that at [github.com/jmhobbs/jmhobbs.github.com](https://github.com/jmhobbs/jmhobbs.github.com).

Thanks for reading!
