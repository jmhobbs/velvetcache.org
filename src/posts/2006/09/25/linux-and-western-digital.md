---
category:
- Geek
creator: admin
date: 2006-09-26
permalink: /2006/09/25/linux-and-western-digital/
tags:
- Computers
- Hardware
- Linux
title: Linux And Western Digital
type: post
wp_id: "45"
---

About a month ago I bought a [Western Digital 320 GB "Extreme Lighted" USB/IEEE 1394 Hard Drive](http://www.westerndigital.com/en/products/products.asp?driveid=153&language=en) as a backup (and general storage) drive.  And yes, I use the term IEEE 1394 because I despise Apple.  Anyway, moving on. It's bigger than both my drive in my main machine put together, two 120 GB Maxtor's that've lasted long enough that I'm assuming they'll die at any moment, nothing against Maxtor, but I had a 120 GB drive from them that died under even less time and stress.  Agian, off topic here.  I bought it because it was cheap for the storage amount, plus it included bling-bling lights.


After a bit of a struggle I managed to format the drive as ext3, but I couldn't do it with `parted`, and ended up using Paragon Partition Manager to coax it into it.  Originally there was two partitions on there, and no disk label.  They were both VFAT's and one was really small.  The lack of a disk label threw me, as I thought you needed that to have a partition table.  Anyway, I got it formatted and hooked up just fine.


I set up a daily backup of my home directory and my photo's directory, which is on a different disk. I use the excellent [Simple Backup Solution](http://sourceforge.net/projects/sbackup/) I got with Automatix.  It backed up fine and didn't complain.  I then started moving some content off of my "Secondary" drive because it was overloaded.  There were a some things like DVD rips (of family occasions, I assure you ;), and some small things like old web and school files.


The small files handle just fine, even over Samba. But when I go to access anything over about 375 Megs, it sticks partway through and requires a hard reset.  I can't even drag things back off onto my main drive, they get stuck partway.  I tried using `split` to break down a 500 Meg movie into 100 Meg chunks, and it got stuck too.


Right now I'm running `fsck` to see if maybe it got messed up in the partitioning process. Oh, and WD doesn't have _any_ Linux help or support beyond "jumper settings".  Guess I'll just have to build my own networked fileserver if this fails.

### 9/26/06, 11:15 pm

I ripped the drive out of it's super cool casing over the last hour or so.  It was harder than it looked, but I was taking special care not to scratch or crack the clear case.  Once I got the drive out I found it to be completely unmarked.  No jumper settings, nothing.  I hooked it up without a jumper as a single device on one of my IDE channels.  BIOS told me it was a "WD3200BB" which came up as a 320 Gb drive from the "Caviar" line.  I went and got the [jumper settings](http://wdc.custhelp.com/cgi-bin/wdc.cfg/php/enduser/std_adp.php?p_faqid=1308&p_created=1106764846#jumper) hooked everything back up and now it's working.

I just transfered about 6 gigs off of the device over Samba to my XP machine, so it was obviously something in either the USB/EIDE interface or else a glitch in the mass-media drivers I have installed in Linux.  Either way, I'm leaving it in the case.  I might take the tiny 16 gig backup drive out of the XP machine and put it in the case, just to see if it will take it.  If it doesn't, whatever.
