---
category:
- Geek
- Life
creator: admin
date: 2018-03-05
permalink: /2018/03/04/chicken-cam-incubator-edition/
tags:
- Chickens
- Electronics
- go
- golang
- Projects
- Raspberry Pi
title: 'Chicken Cam: Incubator Edition'
type: post
wp_id: "2753"
summary: Why let a chicken hatch your eggs when you could make a Raspberry Pi do it?
---
[![Lizzy feeding the chickens scratch grains](https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/lizzy.sm.jpg)](https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/lizzy.jpg)

It's been over a year since we've had chickens and we've missed them, so this Christmas we got Lizzy and Charlotte an incubator so that we could try hatching some this spring.

When we went to purchase eggs, we found that you could most easily get them 10 at a time from the hatchery we have used in the past, [Murray McMurray](https://www.mcmurrayhatchery.com/).  Since the incubator we got the girls could only hold seven, we would need something for the other three.  Some searching found that you could use a styrofoam cooler and a lamp to create a makeshift incubator, so I planned on that.

Once I had a plan to create an incubator, I knew I would have to overcomplicate things.  Four years ago I built a [webcam for our chicks](/2013/10/08/building-the-chicken-cam) so I figured I would do that this time too.  Also, just setting a lamp and thermometer in and hoping for the best seemed like a potential waste of good eggs, so I wanted to monitor the temperature and humidity, and regulate them.

My initial design was a Raspberry Pi connected to a cheap DHT11 temperature and humidity sensor, controlling a relay that could turn the light on and off.  All of it would be hooked up through a PID controller to keep the temperatures right where we want them. Eventually, I added a thermocouple with a MAX6675 for more accurate temperature readings.

[![Raspberry Pi, Relay and a mess of wires.](https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/IMG_0090.sm.jpg)](https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/IMG_0090.jpg)

The server side would be designed similarly to the previous chicken cam, except written in Go.  The stats would be tracked in InfluxDB and Grafana would be used for viewing them.

After I got all the parts I did a little testing, then soldered things up and tested it to see how it ran.

[![Live view from inside the incubator](https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/live.jpg)](https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/live.jpg)

Initially I wrote everything in Go, but the DHT11 reading was very spotty.  Sometimes it would respond once every few seconds, and sometimes it would go a minute or more failing to read.  I wired on a second DHT11 and tried reading from both, but I didn't get that much better performance.

Eventually I tried them from the Adafruit Python library and had much better luck, so I decided to just read those from Python and send them to my main Go application for consumption. I still have trouble with the DHT11's, but I suspect it's my fault more than the sensors fault.

My next issue was that it was extremely jittery, the readings would vary by degrees one second to another, so I collected readings in batches of 5 seconds then averaged them. That smoothed it out enough that graphs looked reasonable.

<figure>
  <img src="https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/jittery.png" alt="A spikey graph of DHT11 readings"
  <figcaption>On. Off. On. Off. On. Off.</figcaption>
</figure>

Temperature was now well regulated, but the air wasn't humid enough.  I switched to a sponge and found I could manage it much easier that way. I briefly tried a 40W bulb thinking I could spend more time with the lamp off, but temperatures still plunged at the same rate when the light was off, so I mostly just created quicker cycles.

After putting the 25W bulb back in, I still wanted a longer, smoother cycle, so I wrapped up a brick (for cleanliness) and stuck that in there.  That got me longer cycles with better recovery at the bottom, it didn't get too cold before the lamp came back on. Some slight improvements to the seal of my lid helped as well. I had trouble with condensation and too much humidity, but some vent holes and better water management took care of that.

<figure>
  <img src="https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/30mins-without-brick-incubator.png" alt="Highly fluctuating temperature graph" />
  <figcaption>Before the brick.</figcaption>
</figure>

<figure>
  <img src="https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/30mins-with-brick-incubator.png" alt="More gradually fluctuating temperature graph" />
  <figcaption class="caption muted">After the brick.</figcaption>
</figure>

For the server side, I mostly duplicated the code from the previous Chicken cam, but in Go.  Then I used the InfluxDB library to get the most recent temperature and humidity readings for display.

![Screenshot of the incubator](http://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/web-interface.png)

At this point, I felt ready for the eggs, which was good because they had arrived!  We placed them in the incubator and we're just waiting now.  On day 8 we candled them with a homebuilt lamp i.e. a cardboard box with a hole cut in it.

[![Candling an egg to view the embryo inside](http://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/candling.sm.jpg)](http://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/candling.jpg)

Things _seem_ to be progressing well so far, so here's hoping something hatches!

## Update!

In May of 2019 we hatched a round of eggs from our own chickens.  We caught the first hatch on camera, and streamed it as it happened with OBS.

[![A chick hatching in the incubator](https://static.velvetcache.org/pages/2018/03/04/chicken-cam-incubator-edition/hatched.png)](https://m.twitch.tv/videos/425869188)
