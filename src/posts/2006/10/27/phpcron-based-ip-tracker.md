---
category:
- Geek
creator: admin
date: 2006-10-27
permalink: /2006/10/27/phpcron-based-ip-tracker/
tags:
- Linux
- PHP
- Programming
title: PHP/Cron Based IP Tracker
type: post
wp_id: "79"
---

I connect to my home machine from work and school fairly regulary, I'd say at least three times a week, often more.  The problem is, COX runs DHCP and every once in a while my lease expires and I get a new IP.  It's not too often, but when I can no longer get to my machine it's frustrating.

No more comrades.  I wrote this little PHP script to handle the problem, and the really neat part (I think) is that it uses itself for the data storage.  Slick, I know.  It's an idea much like the no-ip and dyndns services.  Anyway, here it is, you can just hook up a cron job to keep it up to date.

If you wanted to you could add simple authentication, though if a script-kiddie wants to waste his time flooding the thing with incorrect IP's, what do I care?  It should fix itself in (insert the time spacing of your cron job here) anyway.

P.S.  One caveat here is that thanks to some weird bug or combination of bugs in WP I can't post PHP code that has fopen, fgets, fclose in it.  Thus everywhere you see `// HERE` in the code I've added a space to one of those function calls that needs fixing.  It'd probably be easier just to grab the source [at this link](https://static.velvetcache.org/pages/2006/10/27/php-cron-based-ip-tracker/iptracker.phps).

```php
<?php
  $curIP = "0.0.0.0";
  $newIP = $_SERVER['REMOTE_ADDR'];
  // This should be this file
  $filename = "iptracker.php";
  if($_GET['set'] == 1 && $curIP != $newIP)
  {
    $handle = f open($filename, "r"); // HERE
    // These throw away the first two lines
    f gets($handle); // HERE
    f gets($handle); // HERE
    // Put the rest of the file into a buffer
    while (!feof($handle))
    {
      $buffer .= f gets($handle); // HERE
    }		
    f close($handle); // HERE
    // Write it back
    $handle = f open($filename, 'w+'); // HERE
    $buffer = "
<?php
    \$curIP = \"$newIP\";
    ".$buffer;
    f write($handle, $buffer); // HERE
    f close($handle); // HERE
    print "Updated.";
  }
  else
  {
    print "Current IP: $curIP";
  }
?>
```

#### Update (11/01/06)

I didn't check the operation of this script enough and it turns out there was a bug, at least in my implementation.  It was adding a newline and not removing the old IP address, so it was essentially not working.  Easily fixed, and my new source has replaced the old at the link (though it has not been changed above).

I also realized that I didn't throw in the cron part of this in the first posting.  I use lynx for the job, though wget works too.  Here's my command, I cron it every hour.

```console
lynx -dump http://www.example.com/iptracker.php?set=1 > /dev/null
```
