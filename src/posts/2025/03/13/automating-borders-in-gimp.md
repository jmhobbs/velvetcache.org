---
date: 2025-03-13
tags:
- Python
- GIMP
title: Automating Borders in GIMP
type: post
permalink: /2025/03/13/automating-borders-in-gimp/
summary: Saving time and my wrist editing gifs using Python-Fu
---

Often when making or converting an image to use as a Slack emoji, I will outline it so it is visible on dark and light backgrounds.  This is especially important for text at small sizes.

When there is a lot of frames, selecting, adding border and then filling can take forever.  So, with a 76 frame gif in front of me, I decided to finally try scripting GIMP.

Turns out, it's pretty easy!

Here is the script I wrote to do it;

```python
import gimpcolor

# Cheat way to get the image handle, assuming only one is open
currentImage = gimp.image_list()[0]

# gif's are indexed mode, so convert to RGB
pdb.gimp_image_convert_rgb(currentImage)

# This is the color we will select on, black in this case
selectColor = gimpcolor.RGB(0,0,0)

# For each layer...
for layer in currentImage.layers:
	# Select all of this color on the image (i.e. all my black text)
	pdb.gimp_image_select_color(
		currentImage,
		CHANNEL_OP_ADD,
		layer,
		selectColor
	)
	# Convert the selection into a 4 pixel border
	pdb.gimp_selection_border(currentImage, 4)
	# Fill the while selection with the foreground color
	pdb.gimp_edit_bucket_fill_full(
		layer,
		BUCKET_FILL_FG,
		LAYER_MODE_NORMAL,
		100,
		255,
		False,
		True,
		SELECT_CRITERION_COMPOSITE,
		0,
		0
	)
	# Deselect everything before proceeding to the next layer
	pdb.gimp_selection_none(currentImage)
```

Here's an example, but with border width of 1 instead of 4 since the image is smaller.

**Before**

<img src="https://static.velvetcache.org/pages/2025/03/13/automating-borders-in-gimp/wow.gif" alt="Before" loading="lazy" fetchpriority="low" />

**After**

<img src="https://static.velvetcache.org/pages/2025/03/13/automating-borders-in-gimp/wow-outline.gif" alt="After" loading="lazy" fetchpriority="low" />

There are possibilities for improvement or adaptation, like feathering the select and antialiasing things, but this suits my purpose for now.  The breadth of things you can script in GIMP really surprised me, and once you get ahold of the conventions, it's really pretty simple.
