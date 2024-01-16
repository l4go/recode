package recode_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/l4go/recode"
)

type TextTextFile struct {
	File     string
	TextFile TextFile `json:"-"`
}

func (ttf *TextTextFile) RebuildByType(fsys fs.FS) error {
	f, err := fsys.Open(ttf.File)
	if err != nil {
		return &fs.PathError{Op: "RebuildByType", Path: ttf.File, Err: err} 
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&ttf.TextFile); err != nil {
		return &fs.PathError{Op: "RebuildByType", Path: ttf.File, Err: err} 
	}

	return nil
}

type TextFile struct {
	File string
	Text string `json:"-"`
}

func (tf *TextFile) RebuildByType(fsys fs.FS) error {
	f, err := fsys.Open(tf.File)
	if err != nil {
		return &fs.PathError{Op: "RebuildByType", Path: tf.File, Err: err} 
	}
	defer f.Close()

	text, err := io.ReadAll(f)
	if err != nil {
		return &fs.PathError{Op: "RebuildByType", Path: tf.File, Err: err} 
	}

	tf.Text = string(text)
	return nil
}

func Example() {
	fsys := os.DirFS("testfs")
	f, err := fsys.Open("testfile.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var ttf TextTextFile
	if err := dec.Decode(&ttf); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Before: %q\n", ttf.TextFile.Text)

	if err := recode.RecursiveRebuild(&ttf, fsys); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("After: %q\n", ttf.TextFile.Text)

	// Output:
	// Before: ""
	// After: "test text"
}
