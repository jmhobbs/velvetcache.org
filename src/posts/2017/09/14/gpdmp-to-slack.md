---
category:
- Geek
creator: admin
date: 2017-09-15
permalink: /2017/09/14/gpdmp-to-slack/
tags:
- go
- golang
- Programming
- Projects
title: gpdmp-to-slack
type: post
wp_id: "2731"
summary: >
  Show your currently playing song as your Slack status from Google Play
  Desktop Music Player.
---
When Rdio shut down, I tried a few services before landing on Google Play.  It's not perfect, but it's good enough and it's better than Spotify.  One thing that seemed lacking was a desktop application, but that need was neatly filled by the excellent [GPDMP](https://www.googleplaymusicdesktopplayer.com/).

One lesser known feature of GPDMP is the JSON API, which manifests as a simple JSON file that the application updates with information about the playback.  When Slack announced custom statuses, I though back to the days of instant messaging and the integrations that set your status to the song you were playing.

![Demo](http://static.velvetcache.org.s3.amazonaws.com/temp/slack-gpmdp-demo.gif)

Implementing the link from GPDMP to Slack was, in all, a fairly simple matter.  First, I looked at the JSON file to get a feel for the structure.

```json
{
  "playing": true,
    "song": {
      "title": "Freeze Me",
      "artist": "Death From Above 1979",
      "album": "Outrage! Is Now",
      "albumArt": "https://lh3.go...-e100"
    },
    "rating": {
      "liked": false,
      "disliked": false
    },
    "time": {
      "current": 363509,
      "total": 198000
    },
    "songLyrics": null,
    "shuffle": "NO_SHUFFLE",
    "repeat": "NO_REPEAT",
    "volume": 100
}
```

Short and sweet! Now to represent that in Go for decoding.

```go
type Song struct {
	Title    string
	Artist   string
	Album    string
	AlbumArt string
}

type PlaybackJSON struct {
	Playing bool
	Song    Song
	Rating  struct {
		Liked    bool
		Disliked bool
	}
	Time struct {
		Current int
		Total   int
	}
	SongLyrics string
	Shuffle    string
	Repeat     string
	Volume     int
}
```

I didn't _need_ to represent all the elements, but it's a small structure so I went ahead with it. I didn't embed `Song` because I wanted to write an equality test for that struct on it's own. That will get used later on.

```go
func (a Song) Equal(b Song) bool {
	return a.Title == b.Title && a.Artist == b.Artist && a.Album == b.Album
}
```

Next, I needed a way to monitor that file for updates, which GPDMP does fairly often.  [fsnotify](https://github.com/fsnotify/fsnotify) was the obvious choice, and an easy drop in. I added a time based debounce so that we don't read the file on every update, which would be excessive.  This will delay updates by up to whatever `debounce` is set to, but I'm okay with that trade off.

```go
watcher, err := fsnotify.NewWatcher()
if err != nil {
	log.Fatal(err)
}
defer watcher.Close()

go func() {
	var lastRead time.Time

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				if time.Now().After(lastRead.Add(debounce)) {
					lastRead = time.Now()
					...
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}()

err = watcher.Add(gp.Path)
if err != nil {
	log.Fatal(err)
}
<-done
```

Inside that debounce (at line 16) we open the file, decode it to a new struct and, if it's playing, pass it off to a channel.

```go
f, err := os.Open(event.Name)
if err != nil {
	log.Println(err)
	continue
}

dec := json.NewDecoder(f)
pb := PlaybackJSON{}

err = dec.Decode(&pb)
if err != nil {
	log.Println(err)
	continue
}

if pb.Playing {
	updates <- pb.Song
}
```

So, that's it for getting updates from GPDMP! Less than 100 lines, formatted.  Now I needed to watch that `update` channel and post changes in status to Slack.

I found an excellent [Slack API client](https://github.com/nlopes/slack) on a different project, so I grabbed that.  I started by building a little struct to hold my client and state.

```go
type Slack struct {
	Client       *slack.Client
	CurrentSong  Song
	Set          bool
	InitialText  string
	InitialEmoji string
}
```

Then, during client initialization, we get the current custom status for the user and save it.  This way, when you pause your music, it will revert to whatever you had set before.

```go
func (s *Slack) Init() {
	auth, err := s.Client.AuthTest()
	if err != nil {
		log.Fatal(err)
	}

	user, err := s.Client.GetUserInfo(auth.UserID)
	if err != nil {
		log.Fatal(err)
	}

	s.InitialText = user.Profile.StatusText
	s.InitialEmoji = user.Profile.StatusEmoji
	log.Printf("Initial status: %s %s", s.InitialEmoji, s.InitialText)
}
```

Once it is initialized, we just need to range over our updates channel and post them to Slack when it changes.  We set a timeout, because the GPDMP client won't send updates when the song is paused, or if the app quits updating the file (i.e. you quit GPDMP).  By putting the logic for the timeout on this side, we have less to pass over the channel, and we can revert properly if something goes awry in the api reading goroutine.

```go
func (s *Slack) Sync(emoji string, updates chan Song, revert_after time.Duration) {
	for {
		select {
		case song := <-updates:
			if !s.CurrentSong.Equal(song) {
				log.Printf("Sync: %s by %s\n", song.Title, song.Artist)
				s.Client.SetUserCustomStatus(fmt.Sprintf("%s by %s", song.Title, song.Artist), emoji)
				s.CurrentSong = song
				s.Set = true
			}
		case <-time.After(revert_after):
			if s.Set {
				log.Printf("Reverting Status: %s %s\n", s.InitialEmoji, s.InitialText)
				s.Client.SetUserCustomStatus(s.InitialText, s.InitialEmoji)
				s.CurrentSong = Song{}
				s.Set = false
			}
		}
	}
}
```

A little bit of glue in main and it's ready!

```go
func main() {
	flag.Parse()

	api := NewSlack(os.Getenv("SLACK_TOKEN"))
	gpdmp := &GPDMPAPI{os.Getenv("GPDMPAPI_PATH")}

	api.Init()

	updates := make(chan Song)
	done := make(chan bool)

	go gpdmp.Watch(updates, done, 5*time.Second)
	go api.Sync(config.Emoji, updates, 15*time.Second)
	<-done
}
```

You can browse the source and grab your copy at [github.com/jmhobbs/gpdmp-to-slack](https://github.com/jmhobbs/gpdmp-to-slack)
