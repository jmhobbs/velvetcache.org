---
category:
- Uncategorized
creator: admin
date: 2024-02-12
tags:
- CSS 
- Fonts
- Python
title: Box drawing on the web
type: post
permalink: /2024/02/12/box-drawing-on-the-web
summary: "Boxes: Harder than you thought."
---

I enjoy a nice TUI, and [box-drawing characters](https://en.wikipedia.org/wiki/Box-drawing_character) are the heart and soul of a good one.  However, I noticed that my box-drawing elements were not rendering on this blog with the correct alignment. It looks _fine_ I guess, but not as tight as in a terminal. See how it doesn't line up on the right?

![Slightly misaligned box drawing characters.](https://static.velvetcache.org/pages/2024/02/12/oh-no.png)

Something was clearly up, and as I was working on some other posts that use lots of boxes, I felt like I should fix this.  At first I thought that the font I was using for my source code blocks did not have properly aligned characters.  This wasn't quite convincing because I use(d) [Source Code Pro off of Google Fonts](https://fonts.google.com/specimen/Source+Code+Pro), the same font I use in my editors and terminal.  But, it was my best guess.  I replaced the `font-family` with [GNU Unifont](https://unifoundry.com/unifont/index.html) and it all lined up.

However, when inspecting the page in Firefox I noticed it was using two fonts.  Source Code Pro, and Menlo.

![Two fonts are not always better than one.](https://static.velvetcache.org/pages/2024/02/12/firefox-font-inspector.png)

The box glyphs are being rendered from Menlo, and the contents from Source Code Pro. As an aside, Firefox dev tools are much better than Chrome for fonts.

![Source Code Pro in the sheets](https://static.velvetcache.org/pages/2024/02/12/source-code-pro.png)

![Menlo in the streets](https://static.velvetcache.org/pages/2024/02/12/menlo.png)

The only reason for that to happen (that I know of) is if Source Code Pro did not have the characters in the font.  Google Fonts subsets their fonts, to great benefit using [unicode-range](https://developer.mozilla.org/en-US/docs/Web/CSS/@font-face/unicode-range).  The Unicode block for box drawing is `U+2500..U+257F`.  Looking at the (many) ranges supported by the subsets, box drawing is not included in any them.

```console
$ grep unicode-range < source-code-pro.css
  unicode-range: U+0460-052F, U+1C80-1C88, U+20B4, U+2DE0-2DFF, U+A640-A69F, U+FE2E-FE2F;
  unicode-range: U+0301, U+0400-045F, U+0490-0491, U+04B0-04B1, U+2116;
  unicode-range: U+1F00-1FFF;
  unicode-range: U+0370-0377, U+037A-037F, U+0384-038A, U+038C, U+038E-03A1, U+03A3-03FF;
  unicode-range: U+0102-0103, U+0110-0111, U+0128-0129, U+0168-0169, U+01A0-01A1, U+01AF-01B0, U+0300-0301, U+0303-0304, U+0308-0309, U+0323, U+0329, U+1EA0-1EF9, U+20AB;
  unicode-range: U+0100-02AF, U+0304, U+0308, U+0329, U+1E00-1E9F, U+1EF2-1EFF, U+2020, U+20A0-20AB, U+20AD-20C0, U+2113, U+2C60-2C7F, U+A720-A7FF;
  unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
  unicode-range: U+0460-052F, U+1C80-1C88, U+20B4, U+2DE0-2DFF, U+A640-A69F, U+FE2E-FE2F;
  unicode-range: U+0301, U+0400-045F, U+0490-0491, U+04B0-04B1, U+2116;
  unicode-range: U+1F00-1FFF;
  unicode-range: U+0370-0377, U+037A-037F, U+0384-038A, U+038C, U+038E-03A1, U+03A3-03FF;
  unicode-range: U+0102-0103, U+0110-0111, U+0128-0129, U+0168-0169, U+01A0-01A1, U+01AF-01B0, U+0300-0301, U+0303-0304, U+0308-0309, U+0323, U+0329, U+1EA0-1EF9, U+20AB;
  unicode-range: U+0100-02AF, U+0304, U+0308, U+0329, U+1E00-1E9F, U+1EF2-1EFF, U+2020, U+20A0-20AB, U+20AD-20C0, U+2113, U+2C60-2C7F, U+A720-A7FF;
  unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
```

# DIY font subsets

Time to do it ourselves!  There are some higher level options out there for font manipulation, but it felt like most roads led to `pyftsubset` from [fonttools](https://github.com/fonttools/fonttools).  I started by looking through the list of Unicode blocks and picking out which ones I thought I would need, including box drawing.  I also popped the current web font that was being served into the excellent [Wakamai Fondue](https://wakamaifondue.com/) to take a look at all the characters that Google was serving and make sure they were in the ranges I selected.  Since it is only going to be used for representing code and terminal output on here, I could afford to be a little stingy, and I could probably reduce these ranges further.

Once I had the ranges, I ran them through `pyftsubset`, dropping all the alternative glyphs and other layout features I did not want, and I had a woff2 font file ready to test!  The [`ccmp` feature](https://learn.microsoft.com/en-us/typography/opentype/spec/features_ae#tag-ccmp) is retained (for now), as it is used for glyph composition. Essentially, it allows the font to assemble multiple characters using shared glyphs.

```bash
UNICODES=(
  "U+0020-007E" # ASCII Printables
  "U+00A1-00BF" # Latin-1 symbols
  "U+00D7"      # Multiplication sign
  "U+00F7"      # Division sign
  "U+2000-206F" # General Punctuation
  "U+2190-21FF" # Arrows
  "U+2200-22FF" # Mathematical Operators
  "U+2500-257F" # Box Drawing
)

unicodes_string="${UNICODES[*]}"

pyftsubset SourceCodePro-Regular.ttf \
  --output-file="SourceCodePro-Regular-subset.woff2" \
  --flavor=woff2 \
  --no-hinting \
  --desubroutinize \
  --layout-features="ccmp" \
  --unicodes="${unicodes_string//${IFS:0:1}/,}" 
```

It works!

![So close, but not quite.](https://static.velvetcache.org/pages/2024/02/12/almost.png)

Well, almost.

We still have a missing character that is falling back and messing up one row of the output. That little symbol `⋄` on the bad line is the Diamond Operator and part of the default output from [Hexyl](https://github.com/sharkdp/hexyl), my favorite hex viewer.  It's Unicode value is `U+22C4`, which is from [Mathematical Operators block](https://en.wikipedia.org/wiki/Mathematical_Operators_(Unicode_block)). But, didn't we include that in our subset?  Loading our subset font into Wakamai Fondue, I could see it was not there, even though we included the range!  So I loaded the full, original [`SourceCodePro-Regular.ttf` from Github](https://github.com/adobe-fonts/source-code-pro), and it does not have that character either. It looks like Adobe has not implemented that one. Well now what?!

# Set the cmap table, I'm hungry

Characters in a font are different than glyphs in a font.  Glyphs are arbitrary shapes, and can be used to assemble characters (remember `ccmp`?) for multiple characters, in case you really don't want to make O and 0 look different.  The way this happens internally is something called a [cmap table](https://learn.microsoft.com/en-us/typography/opentype/spec/cmap).

> This table defines mapping of character codes to a default glyph index. Different subtables may be defined that each contain mappings for different character encoding schemes. The table header indicates the character encodings for which subtables are present.

I don't particularly feel like I should be adding my own glyph to this font, and I have little desire to learn [Font Forge](https://fontforge.org/en-US/).  So, what can I do?  Well, how about finding a glyph that is Close Enough™ and mapping it to `U+22C4`?  Poking around in the font I found `U+25CA`, aka "lozenge".  Is it too large? Yes.  Do I care?  Not really, not yet at least.

So how can we edit a cmap table?  Turns out [`fonttools` is our friend here again](https://fonttools.readthedocs.io/en/latest/ttLib/tables/_c_m_a_p.html). Not only does it come with tools, it can be used as a Python module directly to edit fonts.  We load the font, then loop over the cmap subtables and patch in our mapping. It's worth noting that I'm only patching specific tables.

The `table.platformID` refers to what "platform" the mapping is for.  `0` is Unicode, `1` is Macintosh, `3` is Windows, etc.  I only care about Unicode, so I patch those. `table.platEncID` is the encoding of the table.  I'm not concerned with the encoding, but I'm printing it to show how many subtables we are encoding.

```python
from fontTools import ttLib

tt = ttLib.TTFont("SourceCodePro-Regular.ttf")

for table in tt["cmap"].tables:
    if table.platformID == 0:
        if table.cmap[0x25CA] == "lozenge":
            print("Patching table", table.platformID, table.platEncID)
            table.cmap[0x22C4] = "lozenge"

tt.save("SourceCodePro-Regular-patched.ttf")
```

At last!  Everything lines up horizontally!

![Finally](https://static.velvetcache.org/pages/2024/02/12/sweet-relief.png)

# Vertical Alignment

The last thing to do is reduce that `line-height` and close the vertical gaps.  I've scoped this to `code.language-console` so that my code examples keep a little more breathing room between lines, but terminal output will be snug.

```console
┌────────┬─────────────────────────┬─────────────────────────┬────────┬────────┐
│00000000│ ff ff ff ff 54 53 6f 75 ┊ 72 63 65 20 45 6e 67 69 │××××TSou┊rce Engi│
│00000010│ 6e 65 20 51 75 65 72 79 ┊ 00                      │ne Query┊⋄       │
└────────┴─────────────────────────┴─────────────────────────┴────────┴────────┘
```

# Conclusion

I've learned a fair amount about fonts, webfonts and how many things about them I still do not know.  Fonts are complex, and I'm grateful for people who make and release things like fonttools and Font Forge.  I'm looking forward to tweaking my patched font over time, and maybe even creating my own glyph for the Diamond Operator someday.

You can find the full code here: [https://gist.github.com/jmhobbs/b15516c7586d1d14fe08c925783fd9d6](https://gist.github.com/jmhobbs/b15516c7586d1d14fe08c925783fd9d6)
