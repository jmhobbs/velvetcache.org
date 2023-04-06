---
category:
- Geek
creator: admin
date: 2006-08-18
layout: layout.njk
permalink: /2006/08/17/sed-perl-rename/
tags:
- Linux
- Perl
title: sed, perl, rename
type: post
wp_id: "9"
---

This is a post from a previous system.  Information, links and images may not be vaild.

Yesterday I spent 10, maybe 20 minutes poking around the internet looking for some sort, any sort, of batch renaming utility for linux.  I knew there had to be one somewhere, I mean, come on.  I eventually found a batch file that ran into some sed that wasn't really what I wanted, but I figured I could pick apart the sed and figure it out.

Bad idea, I ran it without really testing it, and it ate all of my files, every single one. Luckily I had backups, and I started pulling them off the server.  Today I wrote a Perl script to do the job, and finished off with one more Google search to see if I couldn't find something simpler.

I guess my problem was searching for "**batch** rename" because there were no good results, but I found the nix command _rename_ this time and check out the man page.  Wow, stupid me, I re-invented the wheel today,and didn't do it nearly as well.

Whatever, I can live with that, and it got me writing some Perl, which I haven't done in a long time.  So here for your consumption is my Perl script that <u>will</u> let you blow your foot off, but can get the job done for substitution.


```perl
#!/usr/bin/perl

  if($#ARGV < 2) {
    print "Directory: ";
    $directory = <>;
    chomp($directory);
    print "Replace: ";
    $find = <>;
    chomp($find);
    print "With: ";
    $with = <>;
    chomp($with);
  } else {
    $directory = $ARGV[0];
    $find = $ARGV[1];
    $with = $ARGV[2];
  }

  opendir(DIR, $directory) || die "can't opendir $directory: $!";
  @dots = readdir(DIR);
  closedir DIR;
  foreach(@dots) {
    $copy = $_;
    if( $_ =~ s/$find/$with/ ) {
      system("mv $copy $_");
    }
  }
```

