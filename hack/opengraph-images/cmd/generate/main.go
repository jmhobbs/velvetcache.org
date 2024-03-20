package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jmhobbs/velvetcache.org/hack/opengraph-images/internal/opengraph"
)

func main() {

	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "usage: generate <root> <src> <output> <cache>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  root     the root directory that the other directories are based on")
		fmt.Fprintln(os.Stderr, "  src      the source directory to search for posts (markdown only)")
		fmt.Fprintln(os.Stderr, "  output   the directory to write opengraph images to")
		fmt.Fprintln(os.Stderr, "  cache    the directory to read/write cached opengraph images from")
		os.Exit(1)
	}

	root := os.Args[1]
	source := os.Args[2]
	output := filepath.Join(root, os.Args[3])
	cache := filepath.Join(root, os.Args[4])

	log.Printf("Collecting posts from %s", filepath.Join(root, source))
	log.Printf("Writing opengraph images to %s", output)
	log.Printf("Using cache at %s", cache)

	posts, err := collectPosts(root, source)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d posts, generating missing opengraph images", len(posts))

	skipped := 0
	copiedFromCache := 0
	generated := 0
	missingTitle := 0
	errors := 0

	last_len := 0

	var wg sync.WaitGroup
	postsChan := make(chan *post, len(posts))
	resultsChan := make(chan result, len(posts))
	done := make(chan struct{})

	go func() {
		ctr := 0
		for r := range resultsChan {
			ctr++
			extraSpaces := ""
			if last_len > len(r.Post.SourcePath) {
				extraSpaces = strings.Repeat(" ", last_len-len(r.Post.SourcePath))
			}
			fmt.Print("\r", ctr, " of ", len(posts), " (", skipped+copiedFromCache+generated+missingTitle+errors, ")", " ", r.Post.SourcePath, extraSpaces)
			last_len = len(r.Post.SourcePath)
			if r.Skipped {
				skipped++
			}
			if r.CopiedFromCache {
				copiedFromCache++
			}
			if r.Generated {
				generated++
			}
			if r.MissingTitle {
				missingTitle++
			}
			if r.Err != nil {
				log.Println(r.Err)
				last_len = 0
				errors++
			}
		}
		fmt.Println("")
		close(done)
	}()

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go Worker(cache, output, &wg, postsChan, resultsChan)
	}

	for _, post := range posts {
		postsChan <- post
	}
	close(postsChan)
	wg.Wait()
	close(resultsChan)
	<-done

	log.Printf("Copied from cache : %d", copiedFromCache)
	log.Printf("   Already exists : %d", skipped)
	log.Printf("          Created : %d", generated)
	log.Printf("    Missing title : %d", missingTitle)
	log.Printf("           Failed : %d", errors)
}

type result struct {
	Post            *post
	Err             error
	Skipped         bool
	CopiedFromCache bool
	Generated       bool
	MissingTitle    bool
}

func Worker(cache, output string, wg *sync.WaitGroup, posts <-chan *post, results chan result) {
	defer wg.Done()

	for post := range posts {
		if title, ok := post.Frontmatter["title"]; !ok || title == "" {
			results <- result{Post: post, MissingTitle: true}
			continue
		}

		outputPath := post.Path(output)
		cachePath := post.Path(cache)

		if _, err := os.Stat(outputPath); err == nil {
			results <- result{Post: post, Skipped: true}
			continue
		}

		if _, err := os.Stat(cachePath); err == nil {
			err = copyFile(cachePath, outputPath)
			if err != nil {
				results <- result{Post: post, Err: err}
			} else {
				results <- result{Post: post, CopiedFromCache: true}
				continue
			}
		}

		png, err := opengraph.Generate(post.Frontmatter["title"], "", "")
		if err != nil {
			results <- result{Post: post, Err: err}
			continue
		}

		out, err := os.Create(outputPath)
		if err != nil {
			results <- result{Post: post, Err: err}
			continue
		}
		defer out.Close()

		_, err = io.Copy(out, png)
		if err != nil {
			results <- result{Post: post, Err: err}
			continue
		}

		err = copyFile(outputPath, cachePath)
		if err != nil {
			results <- result{Post: post, Err: err}
		}
		results <- result{Post: post, Generated: true}
	}
}

func copyFile(src, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	return err
}

func readFrontmatter(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmOpen := false

	fm := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if "---" == scanner.Text() {
			if fmOpen {
				return fm, nil
			}
			fmOpen = true
			continue
		}

		if fmOpen {
			parts := strings.SplitN(scanner.Text(), ":", 2)
			if len(parts) != 2 {
				continue
			}
			fm[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return fm, err
	}

	return fm, nil
}
