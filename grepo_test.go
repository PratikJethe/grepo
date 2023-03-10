package main

import (
	"reflect"
	"testing"
)

// {1 text test with no occurences 5 9 test}
type SearchTextTestCase struct {
	Description    string
	ParsedInput    ParsedInput
	TestTextList   []string
	ExpectedOutput []SearchResult
	ErrMsg         string
}
type SearchFileTestCase struct {
	Description    string
	ParsedInput    ParsedInput
	ExpectedOutput []SearchResult
	ErrMsg         string
}
type SearchUserInputTestCase struct {
	Description    string
	ParsedInput    ParsedInput
	ExpectedOutput []SearchResult
	ErrMsg         string
}
type SearchDirectoryTestCase struct {
	Description    string
	ParsedInput    ParsedInput
	ExpectedOutput []SearchResult
	ErrMsg         string
}

func TestSearchText(t *testing.T) {

	testCases := []SearchTextTestCase{
		{
			Description:    "test case for 0 occurence in file",
			ParsedInput:    ParsedInput{SearchQuery: "test"},
			TestTextList:   []string{"text with no occurences"},
			ExpectedOutput: []SearchResult{},
			ErrMsg:         "Test failed for 0 occurence",
		},
		{
			Description:  "test case for 1 occurence in file",
			ParsedInput:  ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text with 1 occurence test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text with 1 occurence test", StartPosition: 22, EndPosition: 26},
			},
			ErrMsg: "Test failed for 1 occurence",
		},
		{
			Description:  "test case for multiple occurence on same line",
			ParsedInput:  ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text with multiple occurences on single line test test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text with multiple occurences on single line test test", StartPosition: 45, EndPosition: 49},
				{LineNumber: 1, LineText: "text with multiple occurences on single line test test", StartPosition: 50, EndPosition: 54},
			},
			ErrMsg: "Test failed for multiple occurence on same line",
		},
		{
			Description:  "test case for multiple occurence on different lines",
			ParsedInput:  ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"test on line one", "on line two test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "test on line one", StartPosition: 0, EndPosition: 4},
				{LineNumber: 2, LineText: "on line two test", StartPosition: 12, EndPosition: 16},
			},
			ErrMsg: "Test failed for multiple occurence on different lines",
		},
		{
			Description:  "test case for case sensitive search",
			ParsedInput:  ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text for case sensitive test", "text for case sensitive Test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text for case sensitive test", StartPosition: 24, EndPosition: 28},
			},
			ErrMsg: "Test failed for case sensitive search",
		},
		{
			Description:  "test case for case insensitive search",
			ParsedInput:  ParsedInput{SearchQuery: "test", IsCaseInsensitive: true},
			TestTextList: []string{"text for case insensitive test", "text for case insensitive TEST"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text for case insensitive test", StartPosition: 26, EndPosition: 30},
				{LineNumber: 2, LineText: "text for case insensitive TEST", StartPosition: 26, EndPosition: 30},
			},
			ErrMsg: "Test failed for case insensitive search",
		},
		{
			Description:  "test case for non exact match (substring matching)",
			ParsedInput:  ParsedInput{SearchQuery: "test"},
			TestTextList: []string{"text for non exact match testing", "text for non exact match test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text for non exact match testing", StartPosition: 25, EndPosition: 29},
				{LineNumber: 2, LineText: "text for non exact match test", StartPosition: 25, EndPosition: 29},
			},
			ErrMsg: "Test failed for non exact match (substring matching)",
		},
		{
			Description:  "test case for exact match",
			ParsedInput:  ParsedInput{SearchQuery: "test", IsExactMatch: true},
			TestTextList: []string{"text for exact match testing", "text for exact match test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 2, LineText: "text for exact match test", StartPosition: 21, EndPosition: 25},
			},
			ErrMsg: "Test failed for for exact match",
		},
	}

	for _, testCase := range testCases {

		results := []SearchResult{}

		for i, text := range testCase.TestTextList {

			result := SearchText(text, testCase.ParsedInput, i+1)

			results = append(results, result...)

		}
		if !reflect.DeepEqual(results, testCase.ExpectedOutput) {
			t.Fatal(testCase.ErrMsg)
		}

	}
}

func TestSearchFile(t *testing.T) {
	testCases1 := SearchFileTestCase{
		Description: "testcase for file not found",
		ParsedInput: ParsedInput{Filename: "wrongpath.txt", SearchQuery: "test"},
		ErrMsg:      "Test failed for file not found",
	}
	_, err := SearchFile(testCases1.ParsedInput)
	if err == nil {
		t.Fatal(testCases1.ErrMsg)
	}

	testCases2 := SearchFileTestCase{
		Description: "testcase for match found in file",
		ParsedInput: ParsedInput{Filename: "test/directory_one/test_data_dir_1.txt", SearchQuery: "test"},
		ExpectedOutput: []SearchResult{
			{LineNumber: 1, LineText: "test data in directory one", StartPosition: 0, EndPosition: 4, FileName: "test/directory_one/test_data_dir_1.txt"},
		},
		ErrMsg: "Test failed for match found in file",
	}

	resluts, err := SearchFile(testCases2.ParsedInput)
	if err != nil {
		t.Fatal(testCases2.ErrMsg)
	}

	if !reflect.DeepEqual(resluts, testCases2.ExpectedOutput) {
		t.Fatal(testCases2.ErrMsg)
	}

}
func TestSearchUserInput(t *testing.T) {

	testcases := []SearchUserInputTestCase{
		{
			Description: "testcase for match found in user input",
			ParsedInput: ParsedInput{UserGivenSearchList: []string{"test", "testing", "lorem"}, SearchQuery: "test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "test", StartPosition: 0, EndPosition: 4},
				{LineNumber: 2, LineText: "testing", StartPosition: 0, EndPosition: 4},
			},
			ErrMsg: "Test failed for match found in user input",
		},
		{
			Description:    "testcase for match not found in user input",
			ParsedInput:    ParsedInput{UserGivenSearchList: []string{"test", "testing", "lorem"}, SearchQuery: "nomatch"},
			ExpectedOutput: []SearchResult{},
			ErrMsg:         "Test failed for match not found in user input",
		},
	}

	for _, testcase := range testcases {
		results := SearchUserInput(testcase.ParsedInput)

		if !reflect.DeepEqual(results, testcase.ExpectedOutput) {
			t.Fatal(testcase.ErrMsg)
		}
	}

}

func TestSearchDirectory(t *testing.T) {

	testcases := []SearchDirectoryTestCase{
		{
			Description: "testcase for match found in files of directory",
			ParsedInput: ParsedInput{
				SearchDirectory: "test",
				SearchQuery:     "test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "test data in directory one", StartPosition: 0, EndPosition: 4, FileName: "test\\directory_one\\test_data_dir_1.txt"},
				{LineNumber: 1, LineText: "test data in directory one", StartPosition: 0, EndPosition: 4, FileName: "test\\directory_two\\test_data_dir_2.txt"},
				{LineNumber: 1, LineText: "test line for one occurence in single line", StartPosition: 0, EndPosition: 4, FileName: "test\\test_data.txt"},
				{LineNumber: 2, LineText: "test for multiple occurences in single line test test", StartPosition: 0, EndPosition: 4, FileName: "test\\test_data.txt"},
				{LineNumber: 2, LineText: "test for multiple occurences in single line test test", StartPosition: 44, EndPosition: 48, FileName: "test\\test_data.txt"},
				{LineNumber: 2, LineText: "test for multiple occurences in single line test test", StartPosition: 49, EndPosition: 53, FileName: "test\\test_data.txt"},
			},
			ErrMsg: "Test failed for match found in files of directory",
		},
		{
			Description: "testcase for match not found in files of directory",
			ParsedInput: ParsedInput{
				SearchDirectory: "test",
				SearchQuery:     "nomatch"},
			ExpectedOutput: []SearchResult{},
			ErrMsg:         "Test failed for match not found in files of directory",
		},
	}
	for _, testcase := range testcases {

		results, err := SearchDirectory(testcase.ParsedInput)
		if err != nil {
			t.Fatal(testcase.ErrMsg)
		}
		if !reflect.DeepEqual(results, testcase.ExpectedOutput) {
			t.Fatal(testcase.ErrMsg)
		}
	}

	testCase := SearchDirectoryTestCase{
		Description: "testcase for directory not found",
		ParsedInput: ParsedInput{
			SearchDirectory: "wrongdir",
			SearchQuery:     "test"},
		ErrMsg: "Test failed for directory not found",
	}

	_, err := SearchDirectory(testCase.ParsedInput)

	if err == nil {
		t.Fatal(testCase.ErrMsg)

	}

}

func TestWriteOutputToFile(t *testing.T) {

	filename := "test/test_data.txt"
	messages := []string{"test message 1", "test message 2"}
	err := WriteOutputToFile(filename, messages)
	if err == nil {
		t.Fatal("Test failed for file already exist")
	}
}
