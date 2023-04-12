---
category:
- Geek
creator: admin
date: 2015-11-15
permalink: /2015/11/15/chat-ops-snack-preview-script/
tags:
- Automation
- ChatOps
- ImageMagick
- Pack
- Python
- S3
title: 'ChatOps: Snack Preview Script'
type: post
wp_id: "2644"
summary: Automating review of random image selection for an iOS app.
---
One of the things we've built at Pack is [Snack](http://snack.packdog.com/), an iOS app which shows users five great dog photos every day.

Ideally, each Snack would be thematic, and curated.  However, curation doesn't always get done.  Either no theme is apparent, or we just don't get to it in time.

When there is no curated snack, we draw from our [Editors Picks](http://packdog.com/home/editors-picks), a collection of the best dog photos available.  This gives us quality content, even when we don't have the chance to pick all five by hand.

There is a downside to this though.  Some of our editors picks are specific to an event or time period.  Like Christmas or Halloween photos.  These wouldn't be very good in a Snack in the middle of May.

So how do we balance the light editorial touch while making sure the random snacks are cohesive?

My solution was a preview injected into our chat.  Every day a script grabs tomorrows Snack, makes a composite, and sends it into Slack for us to review.  If we see something is off, we jump into the admin and fix it.

It's fairly brute force, but a good example of centralizing tools around chat, which is something we try to do.  First we get the photo information, then we download each photo.  We take those and use ImageMagick to create a captioned composite.  Finally, we upload that to S3 and send the link to Slack.

![The Pack Snack from November 13, 2015](https://static.velvetcache.org/pages/2015/11/15/chat-ops-snack-preview-script/snack_2015-11-13.jpg)

This first listing is pretty simple.  We just send a request for tomorrows Snack JSON from the public API.  You might wonder why we don't just set `day_requested` to `tomorrow`, but the API doesn't support that, and neither does the logic of the app.  This runs on a cron job which mails failed runs to us using [cronic](http://habilis.net/cronic/), so we call `raise_for_status()` to explicitly fail on bad HTTP requests.

```python
tomorrow = datetime.date.today() + datetime.timedelta(days=1)

response = requests.get('http://packdog.com/api/snack.json',
                        params={'day_requested': 'today',
                                'year': tomorrow.year,
                                'month': tomorrow.month, 
                                'day': tomorrow.day})

response.raise_for_status()
result = response.json()
```

This next section shells out to download each image into a temporary file.  Order matters in a Snack, so we use a counter instead of taking random filenames.

```python
index = 0
for post in result['posts']:
    os.system("curl -sq %s > /tmp/snack_tmp_%d" % (post['image'], index))
    index += 1
```

ImageMagick is a powerful suite of tools, and in the next block we use the [montage command](http://www.imagemagick.org/script/montage.php) to stitch our photos together.

**`-title`**, as you might imagine, lets us write some title text onto the composite.  **`-geometry`** specifies the size and spacing of each image, 300px square, with 20px of vertical offset from the title. Lastly, **`-tile`** lays out the images in a 5 by 1 grid instead of stacking them in the default method.

```python
p = subprocess.call(["montage",
                     "-title",
                     "Snack for %s\n%s - %s" % (tomorrow, result['title'], result['description']),
                     "-geometry", "300x300+0+20",
                     "-tile", "5x1",
                     "/tmp/snack_tmp_*",
                     "/tmp/snack_comp"])
```

`tinys3` is an awesomely light API to S3 uploads, and then we use the [Slack Incoming Webhook integration](https://api.slack.com/incoming-webhooks) to send the message to chat.

```python
conn = tinys3.Connection("AWS_ACCESS_KEY", "AWS_SECRET_KEY", tls=True)
with open("/tmp/snack_comp", "rb") as handle:
    conn.upload("snack-previews/snack_%s.jpg" % tomorrow, handle, "pack-uploads", public=True)

snack_comp_url = "http://pack-uploads.s3.amazonaws.com/snack-previews/snack_%s.jpg" % tomorrow

requests.post('https://hooks.slack.com/services/WEBHOOK/URL',
              data=json.dumps({"text": snack_comp_url}),
              headers={'content-type': 'application/json'})
```

All in all a simple, effective piece code that drops needed information into our chat every morning.  Here's the whole listing for a complete picture.

```python
import os
import json
import tinys3
import datetime
import requests
import subprocess

os.system("rm -f /tmp/snack_*")

tomorrow = datetime.date.today() + datetime.timedelta(days=1)

response = requests.get('http://packdog.com/api/snack.json',
                        params={'day_requested': 'today',
                                'year': tomorrow.year,
                                'month': tomorrow.month, 
                                'day': tomorrow.day})

response.raise_for_status()
result = response.json()

index = 0
for post in result['posts']:
    os.system("curl -sq %s > /tmp/snack_tmp_%d" % (post['image'], index))
    index += 1

p = subprocess.call(["montage",
                     "-title",
                     "Snack for %s\n%s - %s" % (tomorrow, result['title'], result['description']),
                     "-geometry", "300x300+0+20",
                     "-tile", "5x1",
                     "/tmp/snack_tmp_*",
                     "/tmp/snack_comp"])

conn = tinys3.Connection("AWS_ACCESS_KEY", "AWS_SECRET_KEY", tls=True)
with open("/tmp/snack_comp", "rb") as handle:
    conn.upload("snack-previews/snack_%s.jpg" % tomorrow, handle, "pack-uploads", public=True)

snack_comp_url = "http://pack-uploads.s3.amazonaws.com/snack-previews/snack_%s.jpg" % tomorrow

requests.post('https://hooks.slack.com/services/WEBHOOK/URL',
              data=json.dumps({"text": snack_comp_url}),
              headers={'content-type': 'application/json'})
```

And here is an example of a random Snack that needed correction. Problem solved!

![A preview of a random Snack](http://static.velvetcache.org/pages/2015/11/15/chat-ops-snack-preview-script/snack_2015-11-15.jpg)

Many thanks to [tinys3](https://github.com/smore-inc/tinys3) and [requests](http://docs.python-requests.org/en/latest/).
