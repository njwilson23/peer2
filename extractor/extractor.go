package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	pdfcontent "github.com/njwilson23/unidoc/pdf/contentstream"
	pdf "github.com/njwilson23/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_extract_text.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	counts, err := processContentStreams(inputPath, stopWordsMongoDB())
	for k, v := range counts {
		fmt.Println(k, v)
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func mergeCounts(counts ...map[string]int) map[string]int {
	merged := make(map[string]int)
	for _, count := range counts {
		for word, n := range count {
			if prev, ok := merged[word]; ok {
				merged[word] = prev + n
			} else {
				merged[word] = n
			}
		}
	}
	return merged
}

func wordCount(text string, stopWords map[string]bool) (map[string]int, error) {

	counts := make(map[string]int)

	for _, substr := range strings.Split(strings.Replace(wordNormalize(text), "\n", " ", -1), " ") {

		if len(substr) < 3 {
			continue
		}

		if n, ok := counts[substr]; ok {
			n++
			counts[substr] = n
		} else {
			counts[substr] = 1
		}
	}
	return counts, nil
}

func wordNormalize(word string) string {
	buffer := bytes.NewBuffer([]byte{})

	// remove numbers and punctuation
	for _, char := range strings.ToLower(word) {
		switch char {
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', ' ', '\n':
			buffer.WriteByte(byte(char))
		default:
			continue
		}
	}
	return buffer.String()
}

func processContentStreams(inputPath string, stopWords map[string]bool) (map[string]int, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return nil, err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	var counts map[string]int

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return counts, err
		}

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return counts, err
		}

		// If the value is an array, the effect shall be as if all of the streams in the array were concatenated,
		// in order, to form a single stream.
		pageContentStr := ""
		for _, cstream := range contentStreams {
			//fmt.Println(cstream)
			pageContentStr += cstream
		}

		cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
		txt, err := cstreamParser.ExtractText()
		if err != nil {
			return counts, err
		}

		pageCounts, err := wordCount(txt, stopWords)
		if err != nil {
			return counts, err
		}
		counts = mergeCounts(counts, pageCounts)
	}

	return counts, nil
}
