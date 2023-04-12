---
category:
- Geek
creator: admin
date: 2015-10-07
permalink: /2015/10/06/android-dev-phone/
tags:
- Android
title: $20 Android Dev Phone
type: post
wp_id: "2632"
summary: Setting up a cheap prepaid phone as an alternative to using the emulator.
---
I've been doing a bit of Android dev work lately, just tinkering mostly, and I needed hardware to test on.

Following the advice of local Android guru [Mark](https://twitter.com/MarkCorrado), I just went to Amazon and bought something cheap running at least 4.0.  What I ended up with was a Verizon prepaid [LG Optimus Exceed 2](http://www.amazon.com/gp/product/B00K2XX4OY)

![The $20 Dev Phone](http://static.velvetcache.org/pages/2015/10/06/android-dev-phone/LG-Optimus-Exceed-2.png)

The hardware is decent, especially for a $20 phone.  I've not used it a ton yet, but I would say this would almost be a tolerable daily use phone.  Camera sucks, but who cares.

Here's how I got it set up.

## Skip Activation

Right after boot, you can skip activation with this sequence:

1. Volume Up
2. Volume Down
3. Back Button
4. Home Button

Source: [XDA](http://forum.xda-developers.com/showthread.php?t=2601533)

## Developer Mode

Go to `Settings > General  > About Phone  > Software Information` then tap `Build Number` seven times.

You should find a `Developer Options` in your main settings now.

There's one more step though on OS X. When connecting to USB, you only trigger debug mode by choosing `Internet Connection > Ethernet` for the USB connection method.

![Use ethernet.](http://static.velvetcache.org/pages/2015/10/06/android-dev-phone/ethernet.png)

Source: [Stack Overflow](http://stackoverflow.com/questions/24685768/lg-device-not-listed-in-adb-devices#32516590)

## Root It

The internet would have you believe that rooting your Optimus Exceed 2 is super easy.  For me, it was not.  There is a [one-click guide](http://forum.xda-developers.com/lg-g3/general/guide-root-lg-firmwares-kitkat-lollipop-t3056951) but I had to fall back on the [original set of scripts](http://forum.xda-developers.com/android/development/guide-root-method-lg-devices-t3049772) to make it work, which took a couple tries.

But it worked!

![Look! SuperSu!](http://static.velvetcache.org/pages/2015/10/06/android-dev-phone/SuperSU.png)

## Conclusion

I've not actually used it for a lot of development yet, but I've loaded up some apps with the Play store and it's reasonably responsive.  I like it.


