package grepo

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/pratikjethe/grepo/cmd"
)

type searchResult struct {
	lineNumber    int
	lineText      string
	startPosition int
	endPosition   int
	fileName      string
}

type output struct {
	searchResult     []searchResult
	parsedInput      cmd.ParsedInput
	matchCount       int
	stopOutput       bool
	checkIfFileExist bool
}

/*
GrepSearch function is responsible to call appropriate search function based on user flags.
It taks ParsedInput as input param.
Based on the flags it calls searchUserInput, searchDirectory or searchFile.
results from above functions are passed to handleOutput function
*/
func GrepSearch(parsedInput cmd.ParsedInput) {

	var err error
	var results []searchResult
	var wg sync.WaitGroup

	outputChannel := make(chan output)

	go handleOutputWithChannel(outputChannel, &wg)
	if len(parsedInput.UserGivenSearchList) > 0 {
		results = searchUserInput(parsedInput)
		wg.Add(1)
		outputChannel <- output{searchResult: results, parsedInput: parsedInput, checkIfFileExist: true}
		outputChannel <- output{stopOutput: true}
	} else {

		fileInfo, err := os.Stat(parsedInput.Filename)
		if err != nil {
			log.Fatal(err)
		}
		//check if provided filepath is dir or file
		if fileInfo.IsDir() {
			wg.Add(1)
			searchDirectory(parsedInput, outputChannel)

		} else {
			results, err = searchFile(parsedInput)

			if err != nil {
				log.Fatal(err)
			}
			wg.Add(1)
			outputChannel <- output{searchResult: results, parsedInput: parsedInput, checkIfFileExist: true}
			outputChannel <- output{stopOutput: true}

		}

	}
	// err = handleOutput(results, parsedInput)
	wg.Wait()
	if err != nil {
		log.Fatal(err)
	}

}

/*
searchDirectory function is responsible to find all txt files in a given directory.
It calls searchFile function to get search results, combine them and returns it
*/
func searchDirectory(parsedInput cmd.ParsedInput, outputchannel chan output) {
	var wg sync.WaitGroup
	var matchCount int = 0
	var mu sync.Mutex

	var checkIfFileExist bool = true
	err := filepath.Walk(parsedInput.Filename, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {

			wg.Add(1)
			go func(filePath string, fileCheck bool) {
				newParsedInput := parsedInput
				newParsedInput.Filename = filePath
				results, err := searchFile(newParsedInput)

				if err != nil {
					wg.Done()
					return

				}
				if parsedInput.OnlyCount {
					mu.Lock()
					matchCount = matchCount + len(results)
					mu.Unlock()
				} else {

					outputchannel <- output{searchResult: results, parsedInput: newParsedInput, checkIfFileExist: fileCheck}
				}

				wg.Done()

			}(path, checkIfFileExist)
			if checkIfFileExist {
				checkIfFileExist = false
			}

		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	if parsedInput.OnlyCount {
		outputchannel <- output{parsedInput: parsedInput, matchCount: matchCount}
	}
	outputchannel <- output{stopOutput: true}

}

/*
searchFile function is responsible to return all the matches in a given file.
It open the given file, loops through each line and passes it to searchText function
*/
func searchFile(parsedInput cmd.ParsedInput) ([]searchResult, error) {
	results := []searchResult{}
	file, err := os.Open(parsedInput.Filename)
	if err != nil {
		return results, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		lineText := scanner.Text()

		result := searchText(lineText, parsedInput, lineNum)
		results = append(results, result...)

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return results, err
	}

	return results, nil
}

/*
searchUserInput function is responsible to return all the matches found in user provided input.
It loops through user provided list of words and passes them to searchText function
*/
func searchUserInput(parsedInput cmd.ParsedInput) []searchResult {
	results := []searchResult{}

	for i, word := range parsedInput.UserGivenSearchList {

		result := searchText(word, parsedInput, i+1)
		results = append(results, result...)

	}
	return results

}

/*
searchText function is responsible to return all the matches found in provided text.
It constructs regex based on user flags and checks the match against given text
*/
func searchText(text string, parsedInput cmd.ParsedInput, lineNumber int) []searchResult {
	results := []searchResult{}
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
			result := searchResult{
				lineNumber: lineNumber,
				lineText:   text, startPosition: match[0], endPosition: match[1],
				fileName: parsedInput.Filename,
			}
			results = append(results, result)
		}

	}

	return results
}

/*
constructOutputFromResults function is responsible to return a list of formatted messages based on searchResult array.
It constructs readable messages based usecase (file search or user input search)
*/
func constructOutputFromResults(results []searchResult) []string {
	var readableMessages = []string{}
	var message string
	for _, result := range results {

		if result.fileName != "" {
			message = fmt.Sprintf("Match in file: %v  line  %v:%v \"%v\"", result.fileName, result.lineNumber, result.startPosition, result.lineText)
		} else {
			message = fmt.Sprintf("Match found: %v", result.lineText)

		}
		readableMessages = append(readableMessages, message)

	}
	return readableMessages
}

/*
handleOutput function is responsible to output search result.
Based on flags output is displayed on console or stored in output file.
*/
// func handleOutput(results []searchResult, parsedInput cmd.ParsedInput) error {
// 	var formattedMessages []string

// 	if parsedInput.OnlyCount {
// 		formattedMessages = append(formattedMessages, fmt.Sprintf("Number of matches : %v", len(results)))
// 	} else if (parsedInput.ShowLineBeforeMatch || parsedInput.ShowLinesAfterMatch) && len(results) > 0 {
// 		var err error
// 		formattedMessages, err = getLinesFromFileAroundLineNumber(results[0].fileName, results[0].lineNumber, parsedInput.ShowLinesAfterMatch, parsedInput.ShowLineBeforeMatch)

// 		if err != nil {
// 			return err
// 		}

// 	} else {

// 		formattedMessages = constructOutputFromResults(results)
// 	}

// 	if parsedInput.OutputFilename != "" {
// 		err := writeOutputToFile(parsedInput.OutputFilename, formattedMessages,)

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Println("output stored into " + parsedInput.OutputFilename)
// 	} else {
// 		printMessages(formattedMessages)
// 	}

// 	return nil
// }
func handleOutputWithChannel(outputChannel <-chan output, wg *sync.WaitGroup) {
	for {
		output := <-outputChannel
		if output.stopOutput {
			wg.Done()
			break
		}

		var formattedMessages []string

		if output.parsedInput.OnlyCount {
			formattedMessages = append(formattedMessages, fmt.Sprintf("Number of matches : %v", output.matchCount))
		} else if (output.parsedInput.ShowLineBeforeMatch || output.parsedInput.ShowLinesAfterMatch) && len(output.searchResult) > 0 {
			var err error
			formattedMessages, err = getLinesFromFileAroundLineNumber(output.searchResult[0].fileName, output.searchResult[0].lineNumber, output.parsedInput.ShowLinesAfterMatch, output.parsedInput.ShowLineBeforeMatch)

			if err != nil {
				log.Fatal(err)
			}

		} else {

			formattedMessages = constructOutputFromResults(output.searchResult)
		}

		if output.parsedInput.OutputFilename != "" {
			err := writeOutputToFile(output.parsedInput.OutputFilename, formattedMessages, output.checkIfFileExist)

			if err != nil {
				log.Fatal(err)
			}

		} else {
			printMessages(formattedMessages)
		}

	}

}

/*
writeOutputToFile function is responsible to store output in output file when -o flag is provided.
*/
func writeOutputToFile(filename string, messages []string, checkIfFileExist bool) error {
	var file *os.File
	var err error

	_, err = os.Stat(filename)
	if checkIfFileExist {

		if err == nil {
			return errors.New(filename + " file already exists")
		}
	}
	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	// if err != nil {

	// 	file, err = os.Create(filename)
	// } else {
	// 	file, err = os.OpenFile(filename)
	// }
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

/*
printMessages function is responsible to get line after or before a particular line number.
*/
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
