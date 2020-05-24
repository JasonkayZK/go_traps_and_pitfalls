package main

import "os"

/* Fault: */
/* defer in circulation, run out of file descriptors */
/*
func main() {
	filenames := os.Args[1:]
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		// all file will close by the end of circulation!
		defer f.Close()
		// process of each file
	}
}
*/

/* Solution 1: extract file operation to another function */
func main() {
	filenames := os.Args[1:]
	for _, filename := range filenames {
		if err := func(filename string) error {
			f, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer f.Close()
			// process f
			// ...
			return nil
		}(filename); err != nil {
			panic(err)
		}
	}
}
