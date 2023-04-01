---
category:
- Geek
creator: admin
date: 2018-06-13T22:32:46
layout: layout.njk
tags:
- GIF
- go
- golang
- Party Parrot
title: Party Gopher!
type: post
---
The [Go slack](https://invite.slack.golangbridge.org/) has a cute little dancing Gopher that appears to have come from [Egon Elbre](https://github.com/egonelbre/gophers).  I love it!

<!-- todo: too big -->
![Dancing Gopher](http://static.velvetcache.org/pages/2018/06/13/party-gopher/dancing-gopher.gif)

This little dancing Gopher made me think of [Party Parrot](http://cultofthepartyparrot.com), so I wanted to parrot-ize him.  Normally I might just open up Gimp and start editing, but this is the Go Gopher, we can do better than that!

My plan was to use Go's image packages to edit each frame and replace the blue with the correct parrot color for that frame by walking over the pixels in each frame.

Once I got into the package docs however, I realized that since gif's are paletted, I can just tweak the palette on each frame and be done.  Much simpler.  Let's get into then, shall we?

### Colors!

First things first, I needed to declare the party parrot frame colors, and the light and dark blue that the dancing gopher uses.  I grabbed the blues with [Sip](https://sipapp.io/) and I already had the parrot colors on hand.  Sure, I could precompute these and declare, but let's keep it interesting.

Note that I have a `DarkParrotColors` slice as well, this is for the corresponding dark blue replacements.  I generate these with `darken` which I'll show in a moment.
<!-- todo: start-line="11" mark="3,28" -->
```go
var (
  ParrotColors     []color.Color
  DarkParrotColors []color.Color
  LightGopherBlue  color.Color
  DarkGopherBlue   color.Color
)

func init() {
  var err error

  for _, s := range []string{
    "FF6B6B",
    "FF6BB5",
    "FF81FF",
    "FF81FF",
    "D081FF",
    "81ACFF",
    "81FFFF",
    "81FF81",
    "FFD081",
    "FF8181",
  } {
    c, err := hexToColor(s)
    if err != nil {
      log.Fatal(err)
    }
    ParrotColors = append(ParrotColors, c)
    DarkParrotColors = append(DarkParrotColors, darken(c))
  }

  LightGopherBlue, err = hexToColor("8BD0FF")
  if err != nil {
    log.Fatal(err)
  }
  DarkGopherBlue, err = hexToColor("82C2EE")
  if err != nil {
    log.Fatal(err)
  }
}
```

Also notable is the `hexToColor` which just unpacks an HTML hex RGB representation into a `color.Color`.

<!-- todo: start-line="89" -->
```go
func hexToColor(hex string) (color.Color, error) {
  c := color.RGBA{0, 0, 0, 255}

  r, err := strconv.ParseInt(hex[0:2], 16, 16)
  if err != nil {
    return c, err
  }

  g, err := strconv.ParseInt(hex[2:4], 16, 16)
  if err != nil {
    return c, err
  }

  b, err := strconv.ParseInt(hex[4:6], 16, 16)
  if err != nil {
    return c, err
  }

  c.R = uint8(r)
  c.G = uint8(g)
  c.B = uint8(b)

  return c, nil
}
```

Here is the `darken` function, pretty simple.

<!-- todo: start-line="115" -->
```go
func darken(c color.Color) color.Color {
  r, g, b, a := c.RGBA()
  r = r - 15
  g = g - 15
  b = b - 15
  return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
```

Now I need to pull in the gif and decode it, all very boilerplate.

<!-- todo: start-line="57" -->
```go
  // Open the dancing gopher gif
  f, err := os.Open("dancing-gopher.gif")
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  // Decode the gif so we can edit it
  gopher, err := gif.DecodeAll(f)
  if err != nil {
    log.Fatal(err)
  }
```

After that, I iterate over the frames and edit the palettes.

<!-- todo: start-line="73" -->
```go
  for i, frame := range gopher.Image {
    lbi = frame.Palette.Index(LightGopherBlue)
    dbi = frame.Palette.Index(DarkGopherBlue)

    frame.Palette[lbi] = ParrotColors[i%len(ParrotColors)]
    frame.Palette[dbi] = DarkParrotColors[i%len(DarkParrotColors)]
  }
```

Lastly, more boilerplate to write it out to disk.

<!-- todo: start-line="83" -->
```go
  o, _ := os.OpenFile("party-gopher.gif", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
  defer o.Close()
  gif.EncodeAll(o, gopher)
```

![Party Gopher](http://static.velvetcache.org/pages/2018/06/13/party-gopher/party-gopher.gif)

You can grab the code on [Github](https://github.com/jmhobbs/party-gopher), and thanks again to [Egon Elbre](http://egonelbre.com/) for the excellent original gif!

