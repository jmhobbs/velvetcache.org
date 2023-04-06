---
category:
- Geek
creator: admin
date: 2006-09-08
layout: layout.njk
permalink: /2006/09/08/eclipse-on-a-usb-stick/
tags:
- Computers
- Eclipse
- Java
- Portable
- Programming
- Software
title: Eclipse On A USB Stick
type: post
wp_id: "26"
---

I'm taking a Java course this semester and it was annoying not having access to Eclipse on the various computers in labs around campus.  So I applied some of my Google-fu and found a way to install Eclipse on a USB stick.

The original instructions are [here](http://portableapps.com/node/929), way at the bottom of the page.  I relocated them here and changed them to mirror my experience.

1. Download [Eclipse](http://www.eclipse.org/downloads/) and install it on your USB drive.  I used [EasyEclipse For Desktop Java](http://easyeclipse.org/site/distributions/desktop-java.html) from nexB.
2. Install a [JDK](http://java.sun.com/javase/downloads/) on your or any other PC. The default installation path should be something like &quot;C:\Programs\Java&quot;, containing a folder named &quot;jdk1.5.0_xx&quot;
3. Create a subfolder &quot;JDKs&quot; in your Eclipse folder (Depending on your drive letter and extraction path this should look like E:\Portables\Eclipse\JDKs
4. Copy the jdk1.5.0_xx folder into the JDKs folder. (like `E:\Portables\Eclipse\JDKs\jdk1.5.0_xx`) This can take some time.
5. Edit the file &quot;E:\Portables\Eclipse\eclipse.ini&quot; to
```
-vm
..\JDKsjdk1.5.0_xx\bin\javaw
-vmargs
-Xms40m
-Xmx256m
```
(replace &quot;jdk1.5.0_xx&quot; with the actual folder name) This is more elegant than using a batch file due to being independent from drive letters and path variables.
6. Run &quot;E:\Portables\Eclipse\eclipse.exe&quot;
7. When asked for a workspace location, you can enter &quot;.\workspace&quot;, which will create a workspace in your eclipse folder (like &quot;E:\Portables\Eclipse\workspace\&quot;)
8. Additionally I trimmed down my JDK install by removing the samples and demos directories.

