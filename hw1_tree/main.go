package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type ByAlphabetFilename []os.FileInfo

func (af ByAlphabetFilename) Len() int {
	return len(af)
}

func (af ByAlphabetFilename) Swap(i, j int) {
	af[i], af[j] = af[j], af[i]
}

func (af ByAlphabetFilename) Less(i, j int) bool {
	iRunes := []rune(af[i].Name())
	jRunes := []rune(af[j].Name())

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		if ir != jr {
			return ir < jr
		}
	}

	return len(iRunes) < len(jRunes)
}

func main() {
	out := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}

	_, err = fmt.Fprintln(os.Stdout, out)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(output *bytes.Buffer, path string, viewFile bool) error {
	return getDirTree(output, path, "", viewFile)
}

func getDirTree(output *bytes.Buffer, path string, prefix string, viewFile bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	if !viewFile {
		files = removeFiles(files)
	}

	sort.Sort(ByAlphabetFilename(files))

	newPrefix := ""

	for index, file := range files {

		if len(files) > index+1 {
			addIntermediateBranch(output, file, prefix)
			newPrefix = addParentDirPrefix(prefix)
		} else {
			addFinalBranch(output, file, prefix)
			newPrefix = addParentDirPrefixFinal(prefix)
		}

		if file.IsDir() {
			err = getDirTree(output, addPath(path, file.Name()), newPrefix, viewFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addIntermediateBranch(output *bytes.Buffer, info os.FileInfo, prefix string) {
	addBranch(output, info, prefix, "├───")
}

func addFinalBranch(output *bytes.Buffer, info os.FileInfo, prefix string) {
	addBranch(output, info, prefix, "└───")
}

func addBranch(output *bytes.Buffer, info os.FileInfo, prefix string, graphic string) {
	if !info.IsDir() {
		output.Write([]byte(prefix + graphic + info.Name()))
		addFileSize(output, info)
		addNewLine(output)
		return
	}

	output.Write([]byte(prefix + graphic + info.Name()))
	addNewLine(output)

}

func addNewLine(output *bytes.Buffer) {
	output.Write([]byte("\n"))
}

func removeFiles(files []os.FileInfo) []os.FileInfo {
	dirs := make([]os.FileInfo, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}

	return dirs
}

func addFileSize(output *bytes.Buffer, info os.FileInfo) {
	fileSize := ""
	size := int(info.Size())

	if size == 0 {
		fileSize = " (empty)"
	} else {
		fileSize = " (" + strconv.Itoa(size) + "b)"
	}

	output.Write([]byte(fileSize))
}

func addPath(currentPathName string, targetPathName string) string {
	return currentPathName + "/" + targetPathName
}

func addParentDirPrefix(output string) string {
	output += "│\t"

	return output
}

func addParentDirPrefixFinal(output string) string {
	output += "\t"

	return output
}
