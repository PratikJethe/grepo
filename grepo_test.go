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

func TestSearchText(t *testing.T) {

	testCases := []SearchTextTestCase{
		{
			Description:    "test case for 0 occurence in file",
			ParsedInput:    ParsedInput{SerachQuery: "test"},
			TestTextList:   []string{"text with no occurences"},
			ExpectedOutput: []SearchResult{},
			ErrMsg:         "Test failed for 0 occurence",
		},
		{
			Description:  "test case for 1 occurence in file",
			ParsedInput:  ParsedInput{SerachQuery: "test"},
			TestTextList: []string{"text with 1 occurence test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text with 1 occurence test", StartPosition: 22, EndPosition: 26},
			},
			ErrMsg: "Test failed for 1 occurence",
		},
		{
			Description:  "test case for multiple occurence on same line",
			ParsedInput:  ParsedInput{SerachQuery: "test"},
			TestTextList: []string{"text with multiple occurences on single line test test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text with multiple occurences on single line test test", StartPosition: 45, EndPosition: 49},
				{LineNumber: 1, LineText: "text with multiple occurences on single line test test", StartPosition: 50, EndPosition: 54},
			},
			ErrMsg: "Test failed for multiple occurence on same line",
		},
		{
			Description:  "test case for multiple occurence on different lines",
			ParsedInput:  ParsedInput{SerachQuery: "test"},
			TestTextList: []string{"test on line one", "on line two test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "test on line one", StartPosition: 0, EndPosition: 4},
				{LineNumber: 2, LineText: "on line two test", StartPosition: 12, EndPosition: 16},
			},
			ErrMsg: "Test failed for multiple occurence on different lines",
		},
		{
			Description:  "test case for case sensitive search",
			ParsedInput:  ParsedInput{SerachQuery: "test"},
			TestTextList: []string{"text for case sensitive test", "text for case sensitive Test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text for case sensitive test", StartPosition: 24, EndPosition: 28},
			},
			ErrMsg: "Test failed for case sensitive search",
		},
		{
			Description:  "test case for case insensitive search",
			ParsedInput:  ParsedInput{SerachQuery: "test", IsCaseInsensitive: true},
			TestTextList: []string{"text for case insensitive test", "text for case insensitive TEST"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text for case insensitive test", StartPosition: 26, EndPosition: 30},
				{LineNumber: 2, LineText: "text for case insensitive TEST", StartPosition: 26, EndPosition: 30},
			},
			ErrMsg: "Test failed for case insensitive search",
		},
		{
			Description:  "test case for non exact match (substring matching)",
			ParsedInput:  ParsedInput{SerachQuery: "test"},
			TestTextList: []string{"text for non exact match testing", "text for non exact match test"},
			ExpectedOutput: []SearchResult{
				{LineNumber: 1, LineText: "text for non exact match testing", StartPosition: 25, EndPosition: 29},
				{LineNumber: 2, LineText: "text for non exact match test", StartPosition: 25, EndPosition: 29},
			},
			ErrMsg: "Test failed for non exact match (substring matching)",
		},
		{
			Description:  "test case for exact match",
			ParsedInput:  ParsedInput{SerachQuery: "test", IsExactMatch: true},
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
