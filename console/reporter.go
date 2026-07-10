package console

import (
	"fmt"
)

var REPORTERS = map[string]Reporter{
	"human": HumanReporter{},
	"json": JsonReporter{},
}

func GetReporter(reporter_id string) (Reporter, error) {
	reporter, ok := REPORTERS[reporter_id]
	if !ok{
		return  nil, fmt.Errorf("Unknown reporter: %q", reporter_id)
	}
	return reporter, nil
}


/*
Reporters are meant to deliver operations in different print format
*/
type Reporter interface {
	// reported when a page add been rendered
	PageRendered(pageNum int, path string)
	// reported when antas finished successfully
	Done(outDir string, pageCount int)
	// reported when antas encountered an error
	Error(err error)
}

/*
Base reporter for human readable output.
*/
type HumanReporter struct{}

/*
Base reporter for json parsable output, 
meant for other programs calling antas as another process.
*/

type JsonReporter struct{}

