package main

import (
	"bytes"
	getopt "code.google.com/p/getopt"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestIniConfig(t *testing.T) {
	if config := initConfig([]string{"a.out", "-h"}); !(*config.help) {
		t.Fatal("Help flag expected but not found")
	}
	if config := initConfig([]string{"a.out", "-P"}); !(*config.priv) {
		t.Fatal("Private flag expected but not found")
	}
	if config := initConfig([]string{"a.out", "-u", "user", "-lGo"}); *config.user != "user" {
		t.Fatalf("User 'user' expected, found: [%s]", *config.user)
	} else if *config.lang != "Go" {
		t.Fatalf("Language 'Go' expected, found: [%s]", *config.lang)
	}
}

func createTempFile(t *testing.T, prefix string) *os.File {
	file, err := ioutil.TempFile("", prefix)
	if err != nil {
		t.Fatalf("Could not create temporary file: %s for testing: %s", prefix, err)
	}
	if _, err := file.WriteString(fmt.Sprintf("Testing fpaste command line client with temp file: %s", prefix)); err != nil {
		t.Fatalf("Could not write on temporary file: %s for testing: %s", prefix, err)
	}
	return file
}

type fakeReader struct {
}

func (*fakeReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("Fake reader error")
}

func TestHandleArgs(t *testing.T) {
	/*Testing 3 invalid files*/
	testCommand := getopt.New()
	testCommand.Parse([]string{"a.out", "faketextfile", "fakelogfile", "fakecodefile"})
	if _, errs := handleArgs(nil, testCommand); len(errs) != 3 {
		t.Fatalf("Expected 3 errors, got: %d", len(errs))
	}
	/*Testing 3 valid files*/
	testCommand = getopt.New()
	textfile := createTempFile(t, "textfile")
	defer os.Remove(textfile.Name())
	logfile := createTempFile(t, "logfile")
	defer os.Remove(logfile.Name())
	codefile := createTempFile(t, "codefile")
	defer os.Remove(codefile.Name())
	testCommand.Parse([]string{"a.out", textfile.Name(), logfile.Name(), codefile.Name()})
	if files, errs := handleArgs(nil, testCommand); len(errs) != 0 {
		t.Fatalf("Errors encountered on treatment of valid files")
	} else if len(files) != 3 {
		t.Fatalf("Expected 3 files, got: %d", len(files))
	} else if string(files[0]) != "Testing fpaste command line client with temp file: textfile" {
		t.Fatalf("Unexpected content for textfile, got: %s", string(files[0]))
	} else if string(files[1]) != "Testing fpaste command line client with temp file: logfile" {
		t.Fatalf("Unexpected content for logfile, got: %s", string(files[1]))
	} else if string(files[2]) != "Testing fpaste command line client with temp file: codefile" {
		t.Fatalf("Unexpected content for codefile, got: %s", string(files[2]))
	}
	/*Testing invalid stdin*/
	testCommand = getopt.New()
	testCommand.Parse([]string{"a.out"})
	if _, errs := handleArgs(&fakeReader{}, testCommand); len(errs) != 1 {
		t.Fatalf("Error expected on treatment of invalid stdin")
	}
	/*Testing stdin*/
	testCommand = getopt.New()
	testCommand.Parse([]string{"a.out"})
	testReader := bytes.NewBufferString("Testing fpaste command line client stdin option")
	if files, errs := handleArgs(testReader, testCommand); len(errs) != 0 {
		t.Fatalf("Errors encountered on treatment of stdin")
	} else if len(files) != 1 {
		t.Fatalf("Expected 1 file, got: %d", len(files))
	} else if string(files[0]) != "Testing fpaste command line client stdin option" {
		t.Fatalf("Unexpected content for stdin, got: %s", string(files[0]))
	}
	/*Do not test nil as stdin because os.Stdin is the only option*/
}
