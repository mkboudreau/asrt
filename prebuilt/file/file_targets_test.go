package file

import (
	"strings"
	"testing"
	"time"
)

var testTimeoutDurationOneMinute = (1 * time.Minute)

func TestFileReaderWithEmptyLines(t *testing.T) {
	var tc targetFromFile
	testLine := `              
       

	`
	testReader := strings.NewReader(testLine)

	targets, err := tc.targetsFromReaderWithTimeout(testReader, testTimeoutDurationOneMinute, false)

	if err != nil {
		t.Fail()
		t.Logf("Error is not nil %v", err)
	}

	if len(targets) > 0 {
		t.Fail()
		t.Logf("Targets should have zero length. Found %d, %+v", len(targets), targets)
	}
}

func TestFileReaderWithComment(t *testing.T) {
	var tc targetFromFile
	testLine := `# some line with stuff
	# another line
	# another line
# another line
     # another line
	`
	testReader := strings.NewReader(testLine)

	targets, err := tc.targetsFromReaderWithTimeout(testReader, testTimeoutDurationOneMinute, false)

	if err != nil {
		t.Fail()
		t.Logf("Error is not nil %v", err)
	}

	if len(targets) > 0 {
		t.Fail()
		t.Logf("Targets should have zero length. Found %d, %+v", len(targets), targets)
	}
}

//func (tc *targetFromFile) targetsFromReaderWithTimeout(reader io.Reader, timeout time.Duration) ([]*config.Target, error) {
