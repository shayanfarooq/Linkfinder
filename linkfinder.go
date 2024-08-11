package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// ANSI escape codes for colors
const (
	Reset      = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BrightBlue = "\033[94m"
)

// Regex patterns to extract URLs and paths from JavaScript files
const (
	relUrlRegexStr   = `(?:url\(['"]?([^'"\)]+)['"]?\))|(?:href\s*=\s*['"]([^'"]+)['"])|(?:src\s*=\s*['"]([^'"]+)['"])`
	fetchApiRegexStr = `fetch\s*\(\s*['"]([^'"]+)['"]`
	absoluteUrlRegexStr = `(?:https?:\/\/[^\s"'()<>]+)`
	newRegexStr      = `(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`
)

// Extract URLs from the given content
func extractFromContent(content string) []string {
	var results []string

	// Extract relative URLs and paths
	relRegExp, err := regexp.Compile(relUrlRegexStr)
	if err != nil {
		log.Fatal(err)
	}
	relMatches := relRegExp.FindAllStringSubmatch(content, -1)
	for _, match := range relMatches {
		if match[1] != "" {
			results = append(results, match[1])
		} else if match[2] != "" {
			results = append(results, match[2])
		} else if match[3] != "" {
			results = append(results, match[3])
		}
	}

	// Extract fetch API URLs
	fetchRegExp, err := regexp.Compile(fetchApiRegexStr)
	if err != nil {
		log.Fatal(err)
	}
	fetchMatches := fetchRegExp.FindAllStringSubmatch(content, -1)
	for _, match := range fetchMatches {
		results = append(results, match[1])
	}

	// Extract absolute URLs
	absRegExp, err := regexp.Compile(absoluteUrlRegexStr)
	if err != nil {
		log.Fatal(err)
	}
	absMatches := absRegExp.FindAllString(content, -1)
	for _, match := range absMatches {
		results = append(results, match)
	}

	// Extract additional URLs using new regex pattern
	newRegExp, err := regexp.Compile(newRegexStr)
	if err != nil {
		log.Fatal(err)
	}
	newMatches := newRegExp.FindAllStringSubmatch(content, -1)
	for _, match := range newMatches {
		for _, group := range match[1:] {
			if group != "" {
				results = append(results, group)
			}
		}
	}

	// Remove duplicates
	uniqueResults := unique(results)
	return uniqueResults
}

// Remove duplicates from the slice of strings
func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Process a single URL
func processURL(jsURL string) {
	resp, err := http.Get(jsURL)
	if err != nil {
		log.Printf("Failed to fetch URL %s: %v", jsURL, err)
		return
	}
	defer resp.Body.Close()

	// Read content
	var content strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			content.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Failed to read response body from %s: %v", jsURL, err)
			return
		}
	}

	results := extractFromContent(content.String())

	// Print results
	fmt.Println(Green + "[URL] " + Reset + jsURL)
	for _, result := range results {
		fmt.Println(Cyan + "  " + result + Reset)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		jsURL := strings.TrimSpace(scanner.Text())
		if jsURL != "" {
			processURL(jsURL)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read from standard input: %v", err)
	}
}
