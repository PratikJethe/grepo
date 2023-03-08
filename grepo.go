package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type ParsedInput struct {
	Filename            string
	OutputFilename      string
	SerachQuery         string
	IsCaseInsensitive   bool
	IsExactMatch        bool
	UserGivenSearchList []string
	IsInputProvided     bool
}

type SearchResult struct {
	LineNumber    int
	LineText      string
	StartPosition int
	EndPosition   int
}

func GetParsedInput() ParsedInput {
	filename := flag.String("f", "", "filename of file to be searched")
	searchword := flag.String("s", "", "search word to be searched")
	isCaseInsensitive := flag.Bool("i", false, "makes an case-insensitive search")
	isExactMatch := flag.Bool("e", false, "makes an exact search")
	isUserInputProvided := flag.Bool("input", false, "lets user enter list of words")
	flag.Parse()

	if *searchword == "" {
		flag.Usage()
		log.Fatal("-s (search word) is required")
	}
	if *filename == "" && !*isUserInputProvided {
		flag.Usage()

		log.Fatal("-f (filename) or -input (standard input) is required")
	}
	if *filename != "" && *isUserInputProvided {
		flag.Usage()
		log.Fatal("-f (filename) and -input (standard input) cannot be provided at same time")
	}

	userProvidedSearchlist := []string{}

	if *isUserInputProvided {
		reader := bufio.NewReader(os.Stdin)
		searchString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		userProvidedSearchlist = strings.Fields(searchString)

	}

	return ParsedInput{
		Filename:            *filename,
		SerachQuery:         *searchword,
		IsCaseInsensitive:   *isCaseInsensitive,
		IsExactMatch:        *isExactMatch,
		UserGivenSearchList: userProvidedSearchlist,
		IsInputProvided:     *isUserInputProvided,
	}
}

func GrepSearch(parsedInput ParsedInput) {

	if parsedInput.IsInputProvided {
		fileSearchResult := searchUserInput(parsedInput)

		printMessages(constructRedableMessages(fileSearchResult))
	} else {
		userInputserchresult, err := searchFile(parsedInput)
		if err != nil {
			log.Fatal(err)
		}
		printMessages(constructRedableMessages(userInputserchresult))

	}

}

func searchFile(parsedInput ParsedInput) ([]SearchResult, error) {
	results := []SearchResult{}
	file, err := os.Open(parsedInput.Filename)
	if err != nil {
		return results, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNum := 1

	// Create the regular expression to match the search word

	for scanner.Scan() {
		lineText := scanner.Text()

		result := SearchText(lineText, parsedInput, lineNum)
		results = append(results, result...)

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return results, err
	}

	return results, nil
}

func searchUserInput(parsedInput ParsedInput) []SearchResult {
	results := []SearchResult{}

	for i, word := range parsedInput.UserGivenSearchList {

		result := SearchText(word, parsedInput, i+1)
		results = append(results, result...)

	}
	return results

}

func SearchText(text string, parsedInput ParsedInput, lineNumber int) []SearchResult {

	results := []SearchResult{}
	regexPattern := regexp.QuoteMeta(parsedInput.SerachQuery)
	if parsedInput.IsCaseInsensitive {
		regexPattern = "(?i)" + regexPattern
	}
	if parsedInput.IsExactMatch {
		regexPattern = "\\b" + regexPattern + "\\b"
	}

	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindAllStringIndex(text, -1)

	if len(matches) > 0 {
		for _, match := range matches {
			result := SearchResult{
				LineNumber: lineNumber,
				LineText:   text, StartPosition: match[0], EndPosition: match[1]}
			results = append(results, result)
		}

	}

	return results
}

func constructRedableMessages(results []SearchResult) []string {
	var redableMessages = []string{}

	for _, result := range results {
		message := fmt.Sprintf("Match on  line  %v:%v \"%v\"", result.LineNumber,result.StartPosition, result.LineText)
		redableMessages = append(redableMessages, message)

	}

	return redableMessages
}

func printMessages(messages []string) {
	for _, message := range messages {
		fmt.Println(message)
	}
}
