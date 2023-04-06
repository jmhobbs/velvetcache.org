---
category:
- Geek
creator: admin
date: 2006-12-04
layout: layout.njk
permalink: /2006/12/04/passletcom/
tags:
- BlowPass
- Internet
- JavaScript
- Open Source
- PHP
- Projects
- Security
title: Passlet.com
type: post
wp_id: "105"
---

I saw on my [Ajaxian](http://ajaxian.com/archives/passlet-ajax-password-manager-with-aes-client-side-encryption) feed today a neat service called [Passlet](https://www.passlet.com/).  Essentially it is a password keeper, like [KisKis](http://kiskis.sourceforge.net/) or the one built into Firefox.  The novelty here is that it uses JavaScript to handle all the encrypting and decrypting on the client side.  That means no transmission of clear text information, not even over SSL.

I happily admit I'd been thinking about this concept for at least 4 months.  See, I liked KisKis a lot.  It was Java, used good, solid encryption and had a nice interface.  Problem was, it's hard to keep my thumb drive version synced to my box versions, and I rarely remembered to anyway.  So I thought, why not make a web based password keeper that used JavaScript to keep it secure?

The result was [BlowPass](https://static.velvetcache.org/projects/blowpass/) which uses a JavaScript implementation of the Blowfish cipher.  I was working on the Ajax stuff when I got frustrated with mootools and left it alone. It has several key weaknesses, and I suppose I could learn from Passlet, but, I may as well just use it instead of finishing BlowPass.  If you want the source to BlowPass leave me a note.  Thats my GPL disclaimer since the Blowfish implementation was GPL'd.

#### Update (01/11/07)

BlowPass is semi-active now, you can get more information and try it out at [https://static.velvetcache.org/projects/blowpass](https://static.velvetcache.org/projects/blowpass).  It's still a rather raw version though.  If you aren't concerned about the "open-source" aspect (e.g. don't want to host it and mod it yourself) I'd go use passlet or passpack.
