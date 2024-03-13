package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmhobbs/velvetcache.org/hack/opengraph-images/internal/opengraph"
)

func main() {

	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: generate <src> <output> <cache>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  src      the source directory to search for posts (markdown only)")
		fmt.Fprintln(os.Stderr, "  output   the directory to write opengraph images to")
		fmt.Fprintln(os.Stderr, "  cache    the directory to read/write cached opengraph images from")
		os.Exit(1)
	}

	root := os.Args[1]
	output := os.Args[2]
	cache := os.Args[3]

	log.Printf("Collecting posts from %s", root)

	posts, err := collectPosts(root)
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

	for idx, post := range posts {
		extraSpaces := ""
		if last_len > len(post.SourcePath) {
			extraSpaces = strings.Repeat(" ", last_len-len(post.SourcePath))
		}
		fmt.Print("\r", idx+1, " of ", len(posts), " (", skipped+copiedFromCache+generated+missingTitle+errors, ")", " ", post.SourcePath, extraSpaces)
		last_len = len(post.SourcePath)

		if title, ok := post.Frontmatter["title"]; !ok || title == "" {
			log.Println("No title found in frontmatter for", post.SourcePath)
			missingTitle++
			continue
		}

		outputPath := post.Path(output)
		cachePath := post.Path(cache)

		if _, err := os.Stat(outputPath); err == nil {
			// skip it, it already exists in place
			skipped++
			continue
		}

		if _, err := os.Stat(cachePath); err == nil {
			err = copyFile(cachePath, outputPath)
			if err != nil {
				log.Printf("error copying file from cache: %v", err)
			} else {
				copiedFromCache++
				continue
			}
		}

		png, err := opengraph.Generate(post.Frontmatter["title"], "", "")
		if err != nil {
			log.Printf("unable to generate opengraph image: %v", err)
			errors++
			continue
		}

		out, err := os.Create(outputPath)
		if err != nil {
			log.Printf("unable to create file: %v", err)
			errors++
			continue
		}
		defer out.Close()

		_, err = io.Copy(out, png)
		if err != nil {
			log.Printf("unable to write file: %v", err)
			errors++
			continue
		}

		err = copyFile(outputPath, cachePath)
		if err != nil {
			log.Printf("error copying file to cache: %v", err)
		}
		generated++
	}
	fmt.Println("")

	log.Printf("Copied from cache : %d", copiedFromCache)
	log.Printf("   Already exists : %d", skipped)
	log.Printf("          Created : %d", generated)
	log.Printf("    Missing title : %d", missingTitle)
	log.Printf("           Failed : %d", errors)
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
