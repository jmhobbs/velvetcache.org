---
category:
- Geek
creator: admin
date: 2006-08-04
permalink: /2006/08/04/foxmarks-insecurities/
tags:
- Firefox
- Open Source
- Security
title: Foxmarks Insecurities
type: post
wp_id: "7"
---
Newsflash! [Foxmarks ](http://www.foxcloud.com/wiki/Main_Page)bookmark synchronizer transmits your username and password in cleartext.

I had [LiveHTTP Headers](http://livehttpheaders.mozdev.org/) open while trying to figure out a post error to a server at work when foxmarks went ahead and sync'd up.  I noticed the extra header info and was mildly surprised to find that it had sent my username and password in <u>cleartext</u> over an insecure connection, like so, `http://username:password@sync.foxcloud.com/home/username/foxmarks.xml`

So whats this mean for us? Well, anyone sniffing your traffic (can you say "insecure wireless network"?) will get instant access to your account.  There are no real solutions but you can do a few things to limit the damage.

- Don't use that password on any other site or service.
- Don't auto synchronize on a wireless connection, wait for a hardline if you can.
- Don't put sensitive links or information into foxmarks
