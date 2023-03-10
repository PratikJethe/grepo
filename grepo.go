package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ParsedInput struct {
	Filename            string
	OutputFilename      string
	SearchQuery         string
	IsCaseInsensitive   bool
	IsExactMatch        bool
	UserGivenSearchList []string
	IsInputProvided     bool
	SearchDirectory     string
	OnlyCount           bool
	ShowLinesAfterMatch bool
	ShowLineBeforeMatch  bool
}

type SearchResult struct {
	LineNumber    int
	LineText      string
	StartPosition int
	EndPosition   int
	FileName      string
}

/*
GetParsedInput function is responsible to register flags.
It returns ParsedInput structure which contains all the flags.
It also performs basic validation on given flags.
*/
func GetParsedInput() ParsedInput {
	filename := flag.String("f", "", "filename of file to be searched")
	searchword := flag.String("s", "", "search word to be searched")
	isCaseInsensitive := flag.Bool("i", false, "makes an case-insensitive search")
	isExactMatch := flag.Bool("e", false, "makes an exact search")
	isUserInputProvided := flag.Bool("input", false, "lets user enter list of words")
	outputFilename := flag.String("o", "", "stores output in given file")
	serahDirectory := flag.String("dir", "", "search word in all files of provided directory")
	onlyCount := flag.Bool("c", false, "output only count of matches")
	showLinesAfterMatch := flag.Bool("a", false, "display lines after match")
	showLinesBeforeMatch := flag.Bool("b", false, "display lines before match")
	flag.Parse()

	if *searchword == "" {
		flag.Usage()
		log.Fatal("-s (search word) is required")
	}
	if *filename == "" && !*isUserInputProvided && *serahDirectory == "" {
		flag.Usage()

		log.Fatal("-f (filename) , -input (standard input) or -dir (directory) is required ")
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
		SearchQuery:         *searchword,
		IsCaseInsensitive:   *isCaseInsensitive,
		IsExactMatch:        *isExactMatch,
		UserGivenSearchList: userProvidedSearchlist,
		IsInputProvided:     *isUserInputProvided,
		OutputFilename:      *outputFilename,
		SearchDirectory:     *serahDirectory,
		OnlyCount:           *onlyCount,
		ShowLinesAfterMatch: *showLinesAfterMatch,
		ShowLineBeforeMatch:  *showLinesBeforeMatch,
	}
}

/*
GrepSearch function is responsible to call appropriate search function based on user flags.
It taks ParsedInput as input param.
Based on the flags it calls SearchUserInput, SearchDirectory or SearchFile.
results from above functions are passed to handleOutput function
*/
func GrepSearch(parsedInput ParsedInput) {

	var err error
	var results []SearchResult
	if parsedInput.IsInputProvided {
		results = SearchUserInput(parsedInput)

	} else {
		if parsedInput.SearchDirectory != "" {
			results, err = SearchDirectory(parsedInput)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			results, err = SearchFile(parsedInput)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	err = handleOutput(results, parsedInput)

	if err != nil {
		log.Fatal(err)
	}
}

/*
SearchDirectory function is responsible to find all txt files in a given directory.
It calls SearchFile function to get search results, combine them and returns it
*/
func SearchDirectory(parsedInput ParsedInput) ([]SearchResult, error) {
	combinedResults := []SearchResult{}
	err := filepath.Walk(parsedInput.SearchDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			newParsedInput := parsedInput
			newParsedInput.Filename = path
			results, err := SearchFile(newParsedInput)

			combinedResults = append(combinedResults, results...)
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return combinedResults, nil

}

/*
SearchFile function is responsible to return all the matches in a given file.
It open the given file, loops through each line and passes it to SearchText function
*/
func SearchFile(parsedInput ParsedInput) ([]SearchResult, error) {
	results := []SearchResult{}
	file, err := os.Open(parsedInput.Filename)
	if err != nil {
		return results, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNum := 1

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

/*
SearchUserInput function is responsible to return all the matches found in user provided input.
It loops through user provided list of words and passes them to SearchText function
*/
func SearchUserInput(parsedInput ParsedInput) []SearchResult {
	results := []SearchResult{}

	for i, word := range parsedInput.UserGivenSearchList {

		result := SearchText(word, parsedInput, i+1)
		results = append(results, result...)

	}
	return results

}

/*
SearchText function is responsible to return all the matches found in provided text.
It constructs regex based on user flags and checks the match against given text
*/
func SearchText(text string, parsedInput ParsedInput, lineNumber int) []SearchResult {

	results := []SearchResult{}
	regexPattern := regexp.QuoteMeta(parsedInput.SearchQuery)
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
				LineText:   text, StartPosition: match[0], EndPosition: match[1],
				FileName: parsedInput.Filename,
			}
			results = append(results, result)
		}

	}

	return results
}

/*
constructOutputFromResults function is responsible to return a list of formatted messages based on SearchResult array.
It constructs redable messages based usecase (file search or user input search)
*/
func constructOutputFromResults(results []SearchResult) []string {
	var redableMessages = []string{}
	var message string
	for _, result := range results {

		if result.FileName != "" {
			message = fmt.Sprintf("Match in file: %v  line  %v:%v \"%v\"", result.FileName, result.LineNumber, result.StartPosition, result.LineText)
		} else {
			message = fmt.Sprintf("Match found: %v", result.LineText)

		}
		redableMessages = append(redableMessages, message)

	}
	return redableMessages
}

/*
handleOutput function is responsible to output search result.
Based on flags output is displayed on console or stored in output file.
*/
func handleOutput(results []SearchResult, parsedInput ParsedInput) error {
	var formattedMessages []string

	if parsedInput.OnlyCount {
		formattedMessages = append(formattedMessages, fmt.Sprintf("Number of matches : %v", len(results)))
	} else if (parsedInput.ShowLineBeforeMatch || parsedInput.ShowLinesAfterMatch) && len(results) > 0 {
		var err error
		formattedMessages, err = getLinesFromFileAroundLineNumber(results[0].FileName, results[0].LineNumber, parsedInput.ShowLinesAfterMatch, parsedInput.ShowLineBeforeMatch)

		if err != nil {
			return err
		}

	} else {

		formattedMessages = constructOutputFromResults(results)
	}

	if parsedInput.OutputFilename != "" {
		err := WriteOutputToFile(parsedInput.OutputFilename, formattedMessages)

		if err != nil {
			return err
		}

		fmt.Println("output stored into " + parsedInput.OutputFilename)
	} else {
		printMessages(formattedMessages)
	}

	return nil
}

/*
WriteOutputToFile function is responsible to store output in output file when -o flag is provided.
*/
func WriteOutputToFile(filename string, messages []string) error {
	_, err := os.Stat(filename)

	if err == nil {
		return errors.New(filename + " file already exists")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, message := range messages {
		_, err := file.WriteString(message + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

/*
printMessages function is responsible to print output on console.
*/
func printMessages(messages []string) {
	for _, message := range messages {
		fmt.Println(message)
	}
}

func getLinesFromFileAroundLineNumber(filename string, linenumber int, after bool, before bool) ([]string, error) {
	results := []string{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	currntLineNumber := 1

	for scanner.Scan() {
		lineText := scanner.Text()

		if after && currntLineNumber > linenumber {
			results = append(results, lineText)
		}
		if before && currntLineNumber < linenumber {
			results = append(results, lineText)
		}
		currntLineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
