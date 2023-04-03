---
category:
- Geek
creator: admin
date: 2016-08-30
layout: layout.njk
permalink: /2016/08/29/anonymous-code-of-conduct-reporting-with-twilio/
tags:
- NEJSConf
- PHP
- SMS
- Twilio
title: Anonymous Code of Conduct Reporting With Twilio
type: post
wp_id: "2693"
summary: Integrating with Twilio to allow reporting code of conduct violations over SMS or voice.
---
I'm lucky to be one of the organizers of [NEJS Conf](https://nejsconf.com/), a great little JavaScript & frontend conference here in Omaha, NE.


One of the core values we've held since the beginning of our conference was for it to be diverse, respectful and safe.  To that end we adopted a [Code of Conduct](https://nejsconf.com/code-of-conduct/) from the very beginning, based on the excellent [JSConf example](http://confcodeofconduct.com/).


Our first year, we identified specific volunteers as our CoC points of contact.  It seemed like a good plan, but our only report that year came via a circuitous route, which may have been a result of the face-to-face reporting we had defaulted to.


This spring I got to attend [Twilio's SIGNAL conference](https://www.twilio.com/signal), and one neat thing they had in their Code of Conduct was an anonymous reporting phone line.  Sounds like a good idea, and something fun to build!


The plan is simple: add a Twilio backed phone number which anonymizes incoming SMS and calls, then forwards them to the code of conduct reporting volunteer.  Twilio makes this easy.  At it's core it's just two TwiML files, one for SMS and one for Voice.


The SMS response contains the original message, a destination number to send it to (i.e. the CoC volunteer), a unique ID per reporter, and a link to the web interface.  Behind the scenes we are doing a little but of work to assign the ID, match up numbers to destinations, etc, but not a lot of work total.


{% raw %}
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<Response>
  <Message to="{{ destination_number }}">[{{ reporter_id }}] {{ message_body }}

{{ link }}</Message>
</Response>
```
{% endraw %}

Voice is even simpler. Here we just connect the call to the CoC volunteer, and spoof the caller ID with the hotline's number.

{% raw %}
```xml
<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <Say>Connecting, one moment.</Say>
  <Dial record="true" callerId="{{ caller_id }}">{{ number.destination }}</Dial>
</Response>
```
{% endraw %}

That's the core of it, only took an evening to get things running.  As I hinted above, I added a web interface for replying to reporters, as well as seeing the entire interaction in one place.

![SMS Reporting](https://dl.dropboxusercontent.com/u/21819015/simple-coc/IMG_4120.PNG)
![Web Interface](https://dl.dropboxusercontent.com/u/21819015/simple-coc/IMG_4121.PNG)
![Voice Reporting](https://dl.dropboxusercontent.com/u/21819015/simple-coc/IMG_4122.PNG)

If you'd like to run this for your event, the source is all on github. [https://github.com/jmhobbs/simple-coc](https://github.com/jmhobbs/simple-coc)

