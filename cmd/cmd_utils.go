package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type ParsedInput struct {
	Filename            string
	OutputFilename      string
	SearchQuery         string
	IsCaseInsensitive   bool
	IsExactMatch        bool
	UserGivenSearchList []string
	OnlyCount           bool
	ShowLinesAfterMatch bool
	ShowLineBeforeMatch bool
}

/*
GetParsedInput function is responsible to register flags.
It returns ParsedInput structure which contains all the flags.
It also performs basic validation on given flags.
*/
func GetParsedInput() ParsedInput {
	filename := flag.String("f", "", "accepts file path where search is to be done")
	searchword := flag.String("s", "", "saccepts search query")
	isCaseInsensitive := flag.Bool("i", false, "performs case insensitive search")
	isExactMatch := flag.Bool("e", false, " performs exact matching search")
	outputFilename := flag.String("o", "", "accepts output file path to store output")
	onlyCount := flag.Bool("c", false, "show only count of matches")
	showLinesAfterMatch := flag.Bool("a", false, "display lines after match")
	showLinesBeforeMatch := flag.Bool("b", false, "display lines before match")
	flag.Parse()
	userProvidedSearchlist := []string{}
	if len(*filename) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter input:")
		searchString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		userProvidedSearchlist = strings.Fields(searchString)

	}

	if *searchword == "" {
		flag.Usage()
		log.Fatal("-s (search word) is required")
	}


	return ParsedInput{
		Filename:            *filename,
		SearchQuery:         *searchword,
		IsCaseInsensitive:   *isCaseInsensitive,
		IsExactMatch:        *isExactMatch,
		UserGivenSearchList: userProvidedSearchlist,
		OutputFilename:      *outputFilename,
		OnlyCount:           *onlyCount,
		ShowLinesAfterMatch: *showLinesAfterMatch,
		ShowLineBeforeMatch: *showLinesBeforeMatch,
	}
}
