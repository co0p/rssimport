package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/xml"
	"html/template"
)

type Rss2 struct {
	XMLName		xml.Name	`xml:"rss"`
	Version		string		`xml:"version,attr"`
	// Required
	Title		string		`xml:"channel>title"`
	Link		string		`xml:"channel>link"`
	Description	string		`xml:"channel>description"`
	// Optional
	PubDate		string		`xml:"channel>pubDate"`
	ItemList	[]Item		`xml:"channel>item"`
}

type Item struct {
	// Required
	Title		string		`xml:"title"`
	Link		string		`xml:"link"`
	Description	template.HTML	`xml:"description"`
	// Optional
	Content		template.HTML	`xml:"encoded"`
	PubDate		string		`xml:"pubDate"`
	Comments	string		`xml:"comments"`
}

func (i Item) String() string {
	return fmt.Sprintf("**%s**\n- %s\n- %s\n- %s\n\n", i.Title, i.PubDate, i.Description, i.Link)
}

var feedUrl string

// init initializes the command line arguments
func init() {
	flag.StringVar(&feedUrl, "input", "", "The path to the rss feed")
	flag.StringVar(&feedUrl, "i", "", "The path to the rss feed (shorthand)")
	flag.StringVar(&feedUrl, "output", "", "The path to the output directory")
	flag.StringVar(&feedUrl, "o", "", "The path to the output directory (shorthand)")
}

// FetchFeed fetches the given feedUrl and returns the response body
func FetchFeed(feedUrl string) (error, []byte) {

	response, err := http.Get(feedUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err, nil
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "reading %s: %v\n", feedUrl, err)
		return err, nil
	}

	return err, body
}

// ParseFeed parses the given input feed and returns the rss items found
func ParseFeed(input []byte) (error, items []Item){
	var data Rss2
	xml.Unmarshal(input, &data)
	return nil, data.ItemList
}

func main() {
	flag.Parse()

	if len(feedUrl) < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var _, result = FetchFeed(feedUrl)
	var _, items = ParseFeed(result)

	for _, item := range items {
		fmt.Printf("%s", item)
	}
}
