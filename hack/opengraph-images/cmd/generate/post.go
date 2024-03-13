package main

import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type post struct {
	SourcePath  string
	Frontmatter map[string]string
}

func (p *post) Path(root string) string {
	return filepath.Join(root, p.Name())
}

func (p *post) Name() string {
	h := sha1.New()
	h.Write([]byte(p.SourcePath))
	bs := h.Sum(nil)
	return fmt.Sprintf("%s-%x.png", slugify(p.Frontmatter["title"]), bs)
}

var nonAlphaNum *regexp.Regexp = regexp.MustCompile("[^a-z0-9]+")
var dasher *regexp.Regexp = regexp.MustCompile("-+")

func slugify(s string) string {
	return dasher.ReplaceAllString(nonAlphaNum.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
func collectPosts(root string) ([]*post, error) {
	posts := []*post{}

	fileSystem := os.DirFS(root)

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasSuffix(path, ".md") {
			fm, err := readFrontmatter(filepath.Join(root, path))
			if err != nil {
				log.Println(err)
				return nil
			}
			posts = append(posts, &post{Frontmatter: fm, SourcePath: path})
		}
		return nil
	})

	return posts, err
}
