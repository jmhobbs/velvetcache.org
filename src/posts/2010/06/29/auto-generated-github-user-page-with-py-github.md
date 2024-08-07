---
category:
- Geek
creator: admin
date: 2010-06-29
permalink: /2010/06/29/auto-generated-github-user-page-with-py-github/
tags:
- Github
- Programming
- Python
- Snippets
- Tools
title: Auto-Generated Github User Page With py-github
type: post
wp_id: "1211"
---

<aside>

**Update (2010-06-30)**

So I got antsy about this and I upgraded to using [pystache](http://github.com/defunkt/pystache) instead of my homebrew templating system.  This was my first run in with mustache, and I have to say I like it, even though I used the bare minimum feature set.

New code is at [http://github.com/jmhobbs/jmhobbs.github.com](http://github.com/jmhobbs/jmhobbs.github.com)

</aside>

Github has a cool feature called "[Github Pages](http://pages.github.com/)" that let you host static content on a subdomain of github, e.g. [http://jmhobbs.github.com/](http://jmhobbs.github.com/).

They also provide an auto-generator for project pages that has a nice clean format which I really like.  So I decided to make my user page match the look and feel of the project pages.  And to boot I wanted to be able have it auto-generate since I want it to be "hands free", otherwise I'll forget to update it.

To make this happen I whipped up my template and then grabbed the excellent [py-github](http://github.com/dustin/py-github) from Dustin Sallings, which I have [used before](http://jmhobbs.github.com/github-watcher/).

Without furthur ado I'll just show you the source. It's not complicated, just some API calls then search replace on a template file.  If you want to use it, be sure to get the most recent version from [http://github.com/jmhobbs/jmhobbs.github.com](http://github.com/jmhobbs/jmhobbs.github.com).

Throw in a cron job and you are set. Beware of lot's of "page build" notices from Github though.

```python
# -*- coding: utf-8 -*-

import github.github as github
import yaml
import time
from datetime import datetime

def repo_date_to_epoch ( date ):
  epoch = time.mktime(
    time.strptime(
      date[0:-6],
      "%Y-%m-%dT%H:%M:%S"
    )
  )
  return int( epoch )

def main ():

  print "Loading settings...."
  f = open( 'settings.yaml' )
  settings = yaml.load( f )
  f.close()

  gh = github.GitHub()

  print "Fetching user information..."
  user = gh.users.show( settings['username'] )

  print "Fetching repository information..."
  repos = gh.repos.forUser( settings['username'] )

  print "Sorting repositories..."
  repos = sorted( repos, cmp=lambda a, b: repo_date_to_epoch( b.pushed_at ) - repo_date_to_epoch( a.pushed_at ) )

  print "Loading template..."
  f = open( 'index.html.tpl' )
  template = f.read()
  f.close()

  print "Mangling template..."
  template = template.replace( '<% username %>', settings['username'] )
  template = template.replace( '<% fullname %>', user.name )
  template = template.replace( '<% email %>', user.email )
  template = template.replace( '<% following %>', str( user.following_count ) )
  template = template.replace( '<% followers %>', str( user.followers_count ) )
  template = template.replace( '<% publicrepos %>', str( user.public_repo_count ) )

  repo_string = ''

  for repo in repos:
    if repo.private:
      continue

    repo_string = repo_string + '<div class="repo"><h3><a href="' + repo.url + '">' + repo.name + '</a>'

    try:
      repo_string = repo_string + ' - <span class="small"><a href="' + repo.homepage + '">' + repo.homepage + '</a></span>'
    except AttributeError:
      pass

    repo_string = repo_string + '</h3>'

    repo_string = repo_string + "Forks: " + str( repo.forks ) + " - Watchers: " + str( repo.watchers ) + ' | '

    if repo.has_issues:
      repo_string = repo_string + ' <a href="' + repo.url + '/issues">Issues</a> |'

    if repo.has_wiki:
      repo_string = repo_string + ' <a href="http://wiki.github.com/' + settings['username'] + '/' + repo.name + '">Wiki</a> |'

    if repo.has_downloads:
      repo_string = repo_string + ' <a href="' + repo.url + '/downloads">Downloads</a> |'

    repo_string = repo_string + '<br/>Last Push: ' + datetime.fromtimestamp( repo_date_to_epoch( repo.pushed_at ) ).ctime()

    try:
      repo_string = repo_string + '<pre>' + repo.description + '< /pre>'
    except AttributeError:
      repo_string = repo_string + '<br/><br/>'
      pass

    repo_string = repo_string + "</div><!--// .repo //-->\n"

  template = template.replace( '<% repos %>', repo_string )

  ga = """
    <script type="text/javascript">
      var _gaq = _gaq || [];
      _gaq.push(['_setAccount', '<% ga_code %>']);
      _gaq.push(['_trackPageview']);
      (function() {
        var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
        ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
        var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
      })();
    </script>
  """

  if False != settings['google_analytics']:
    template = template.replace( '<% google_analytics %>', ga )
    template = template.replace( '<% ga_code %>', settings['google_analytics'] )
  else:
    template = template.replace( '<% google_analytics %>', '' )

  print "Writing file..."
  f = open( 'index.html', 'w' )
  f.write( template )
  f.close()

  print "Done!"

if __name__ == "__main__":
  main()
```

Wow. You actually scrolled through all of that. Amazing.
