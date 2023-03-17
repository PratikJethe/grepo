package grepo

import (
	"reflect"
	"testing"

	"github.com/pratikjethe/grepo/cmd"
)

// {1 text test with no occurences 5 9 test}
type searchTextTestCase struct {
	Description    string
	ParsedInput    cmd.ParsedInput
	TestTextList   []string
	ExpectedOutput []searchResult
	ErrMsg         string
}
type searchFileTestCase struct {
	Description    string
	ParsedInput    cmd.ParsedInput
	ExpectedOutput []searchResult
	ErrMsg         string
}
type searchUserInputTestCase struct {
	Description    string
	ParsedInput    cmd.ParsedInput
	ExpectedOutput []searchResult
	ErrMsg         string
}
type searchDirectoryTestCase struct {
	Description    string
	ParsedInput    cmd.ParsedInput
	ExpectedOutput []searchResult
	ErrMsg         string
}

func TestSearchText(t *testing.T) {

	testCases := []searchTextTestCase{
		{
			Description:    "test case for 0 occurrencein file",
			ParsedInput:    cmd.ParsedInput{SearchQuery: "test"},
			TestTextList:   []string{"text with no occurences"},
			ExpectedOutput: []searchResult{},
			ErrMsg:         "Test failed for 0 occurence",
		},
		{
			Description:  "test case for 1 occurrencein file",
			ParsedInput:  cmd.ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text with 1 occurrencetest"},
			ExpectedOutput: []searchResult{
				{lineNumber: 1, lineText: "text with 1 occurrencetest", startPosition: 22, endPosition: 26},
			},
			ErrMsg: "Test failed for 1 occurence",
		},
		{
			Description:  "test case for multiple occurrenceon same line",
			ParsedInput:  cmd.ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text with multiple occurences on single line test test"},
			ExpectedOutput: []searchResult{
				{lineNumber: 1, lineText: "text with multiple occurences on single line test test", startPosition: 45, endPosition: 49},
				{lineNumber: 1, lineText: "text with multiple occurences on single line test test", startPosition: 50, endPosition: 54},
			},
			ErrMsg: "Test failed for multiple occurrenceon same line",
		},
		{
			Description:  "test case for multiple occurrenceon different lines",
			ParsedInput:  cmd.ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"test on line one", "on line two test"},
			ExpectedOutput: []searchResult{
				{lineNumber: 1, lineText: "test on line one", startPosition: 0, endPosition: 4},
				{lineNumber: 2, lineText: "on line two test", startPosition: 12, endPosition: 16},
			},
			ErrMsg: "Test failed for multiple occurrenceon different lines",
		},
		{
			Description:  "test case for case sensitive search",
			ParsedInput:  cmd.ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text for case sensitive test", "text for case sensitive Test"},
			ExpectedOutput: []searchResult{
				{lineNumber: 1, lineText: "text for case sensitive test", startPosition: 24, endPosition: 28},
			},
			ErrMsg: "Test failed for case sensitive search",
		},
		{
			Description:  "test case for case insensitive search",
			ParsedInput:  cmd.ParsedInput{SearchQuery: "test", IsCaseInsensitive: true},
			TestTextList: []string{"text for case insensitive test", "text for case insensitive TEST"},
			ExpectedOutput: []searchResult{
				{lineNumber: 1, lineText: "text for case insensitive test", startPosition: 26, endPosition: 30},
				{lineNumber: 2, lineText: "text for case insensitive TEST", startPosition: 26, endPosition: 30},
			},
			ErrMsg: "Test failed for case insensitive search",
		},
		{
			Description:  "test case for non exact match (substring matching)",
			ParsedInput:  cmd.ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text for non exact match testing", "text for non exact match test"},
			ExpectedOutput: []searchResult{
				{lineNumber: 1, lineText: "text for non exact match testing", startPosition: 25, endPosition: 29},
				{lineNumber: 2, lineText: "text for non exact match test", startPosition: 25, endPosition: 29},
			},
			ErrMsg: "Test failed for non exact match (substring matching)",
		},
		{
			Description:  "test case for exact match",
			ParsedInput:  cmd.ParsedInput{SearchQuery: "test", IsExactMatch: true},
			TestTextList: []string{"text for exact match testing", "text for exact match test"},
			ExpectedOutput: []searchResult{
				{lineNumber: 2, lineText: "text for exact match test", startPosition: 21, endPosition: 25},
			},
			ErrMsg: "Test failed for for exact match",
		},
	}

	for _, testCase := range testCases {

		results := []searchResult{}

		for i, text := range testCase.TestTextList {

			result := searchText(text, testCase.ParsedInput, i+1)

			results = append(results, result...)

		}
		if !reflect.DeepEqual(results, testCase.ExpectedOutput) {
			t.Fatal(testCase.ErrMsg)
		}

	}
}

func TestSearchFile(t *testing.T) {
	testCases1 := searchFileTestCase{
		Description: "testcase for file not found",
		ParsedInput: cmd.ParsedInput{Filename: "wrongpath.txt", SearchQuery: "test"},
		ErrMsg:      "Test failed for file not found",
	}
	_, err := searchFile(testCases1.ParsedInput)
	if err == nil {
		t.Fatal(testCases1.ErrMsg)
	}

	testCases2 := searchFileTestCase{
		Description: "testcase for match found in file",
		ParsedInput: cmd.ParsedInput{Filename: "../test/directory_one/test_data_dir_1.txt", SearchQuery: "test"},
		ExpectedOutput: []searchResult{
			{lineNumber: 1, lineText: "test data in directory one", startPosition: 0, endPosition: 4, fileName: "../test/directory_one/test_data_dir_1.txt"},
		},
		ErrMsg: "Test failed for match found in file",
	}

	resluts, err := searchFile(testCases2.ParsedInput)

	if err != nil {
		t.Fatal(testCases2.ErrMsg)
	}

	if !reflect.DeepEqual(resluts, testCases2.ExpectedOutput) {
		t.Fatal(testCases2.ErrMsg)
	}

}
func TestSearchUserInput(t *testing.T) {

	testcases := []searchUserInputTestCase{
		{
			Description: "testcase for match found in user input",
			ParsedInput: cmd.ParsedInput{UserGivenSearchList: []string{"test", "testing", "lorem"}, SearchQuery: "test"},
			ExpectedOutput: []searchResult{
				{lineNumber: 1, lineText: "test", startPosition: 0, endPosition: 4},
				{lineNumber: 2, lineText: "testing", startPosition: 0, endPosition: 4},
			},
			ErrMsg: "Test failed for match found in user input",
		},
		{
			Description:    "testcase for match not found in user input",
			ParsedInput:    cmd.ParsedInput{UserGivenSearchList: []string{"test", "testing", "lorem"}, SearchQuery: "nomatch"},
			ExpectedOutput: []searchResult{},
			ErrMsg:         "Test failed for match not found in user input",
		},
	}

	for _, testcase := range testcases {
		results := searchUserInput(testcase.ParsedInput)

		if !reflect.DeepEqual(results, testcase.ExpectedOutput) {
			t.Fatal(testcase.ErrMsg)
		}
	}

}

func TestWriteOutputToFile(t *testing.T) {

	filename := "test/test_data.txt"
	messages := []string{"test message 1", "test message 2"}
	err := writeOutputToFile(filename, messages, true)
	if err == nil {
		t.Fatal("Test failed for file already exist")
	}
}
