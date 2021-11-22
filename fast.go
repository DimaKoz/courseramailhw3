package main

import (
	"bufio"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*

go test -bench . -benchmem

Before changes on my machine:

BenchmarkSlow-8               40          28750363 ns/op        19575028 B/op     189817 allocs/op
BenchmarkFast-8               40          28746761 ns/op        19543444 B/op     189812 allocs/op
PASS
ok      mailhw3 3.306s

After changes on my machine:

BenchmarkSlow-8               37          28363908 ns/op        19542307 B/op     189813 allocs/op
BenchmarkFast-8              464           2542303 ns/op          438153 B/op       4489 allocs/op
PASS
ok      mailhw3 3.036s

Original benchmark to compare:

BenchmarkSlow-8 			  10 		 142703250 ns/op 	   336887900 B/op 	  284175 allocs/op
BenchmarkSolution-8 		 500 	       2782432 ns/op 		  559910 B/op 	   10422 allocs/op

*/

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		errCl := file.Close()
		if errCl != nil {
			panic(errCl)
		}
	}()

	seenBrowsers := make(map[string]bool, 150)

	scanner := bufio.NewScanner(file)
	var maxBufSize = 1024
	buf := make([]byte, maxBufSize)
	scanner.Buffer(buf, maxBufSize)
	isAndroid := false
	isMSIE := false

	userCounter := -1
	_, err = fmt.Fprintln(out, "found users:")
	if err != nil {
		panic(err)
	}
	for scanner.Scan() {
		isAndroid = false
		isMSIE = false
		var bytes = scanner.Bytes()
		userCounter++

		_, err := jsonparser.ArrayEach(bytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if err != nil {
				panic(err)
			}
			browser := string(value)
			if strings.Contains(browser, "Android") {
				isAndroid = true
				countUniqueBrowser(seenBrowsers, browser)
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				countUniqueBrowser(seenBrowsers, browser)
			}
		}, "browsers")
		if err != nil {
			panic(err)
		}
		if !isAndroid || !isMSIE {
			continue
		}
		userName, err := jsonparser.GetString(bytes, "name")
		if err != nil {
			panic(err)
		}
		userEmail, err := jsonparser.GetString(bytes, "email")
		if err != nil {
			panic(err)
		}
		userEmail = strings.Replace(userEmail, "@", " [at] ", 1)

		_, err = fmt.Fprint(out, "["+strconv.Itoa(userCounter)+"] "+userName+" <"+userEmail+">\n")
		if err != nil {
			panic(err)
		}

	}

	_, err = fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
	if err != nil {
		panic(err)
	}

}

func countUniqueBrowser(browserNames map[string]bool, browser string) {
	if _, found := (browserNames)[browser]; !found {
		(browserNames)[browser] = true
	}
}

func main() {
	FastSearch(ioutil.Discard)
}
