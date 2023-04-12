---
category:
- Geek
creator: admin
date: 2006-11-16
permalink: /2006/11/16/wlan-router-monitor-bash-script/
tags:
- BASH
- Programming
- Security
title: WLAN Router Monitor BASH Script
type: post
wp_id: "96"
---

I was reading some material on WEP and WPA cracking, and decided to write a monitor for our router. I was curious if anyone other than us hooked up.  I've turned off the MAC filtering on it and got my BASH script working.

I'm kinda proud of it actually.  I wrote it from scratch, just hit up the man pages on my system for hints.  I wanted to use `lynx -dump` on the "Attached Devices" page, but I couldn't get lynx to authenticate from the command line.  I decided to use `wget` instead, since it worked just fine.  I also knew I had `html2txt` installed from something else, so that was good too.

Here's the script, password removed of course. The slickness is all in that last line, pardon my WP plug-in's poor highlighting.

```bash
#!/bin/bash
date >> $HOME/System/routerLog
echo >> $HOME/System/routerLog
wget -O - http://admin:PASSWORD@192.168.1.1/DEV_device.htm | html2text >> $HOME/System/routerLog
```

The one tough feature to find was getting `wget` to print to the stdout instead of to a file.  Thats what the `-O - ` does.

It works nice, but it has a lot of extra spacing in it.  I tried to do a `sed` line to filter out multiple newlines, but I've never actually used `sed`, and I couldn't get it working.  Maybe it's my regex: `s/\n{2,}//g`.  Dunno, not a biggie.  I hooked up a cron job for every 15 minutes, we'll see what I catch (and how bloated that log file will get)

P.S. The router is a NetGear WGR614 v5

P.P.S. I got to thinking about that comment about MAC filtering.  With a big network you could camp out with Kismet, grab some attached devices MAC, wait until it disconnects and change your MAC to it's.  While you wait you can crack the WEP too.  MAC filtering really isn't as good as I thought.  Same with [fakeap](http://www.blackalchemy.to/project/fakeap/).  If only one ap has attached devices...uh, that'd be the real one...
