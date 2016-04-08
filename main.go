package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
)

var feedUrl string

// init initializes the command line arguments
func init() {
  flag.StringVar(&feedUrl, "input", "", "The path to the rss feed")
  flag.StringVar(&feedUrl, "i", "", "The path to the rss feed (shorthand)")
  flag.StringVar(&feedUrl, "output", "", "The path to the output directory")
  flag.StringVar(&feedUrl, "o", "", "The path to the output directory (shorthand)")
}

// FetchFeed fetches the given feedUrl and returns the result as string
func FetchFeed(feedUrl string) (error, string) {

  response, err := http.Get(feedUrl)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err)
    return err, ""
  }

  body, err := ioutil.ReadAll(response.Body)
  response.Body.Close()

  if err != nil {
    fmt.Fprintf(os.Stderr, "reading %s: %v\n", feedUrl, err)
    return err, ""
  }

  return err, string(body)
}

func main() {
  flag.Parse()

  if len(feedUrl) < 1 {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
  }

  var _, result = FetchFeed(feedUrl)
  fmt.Println(result)
}
