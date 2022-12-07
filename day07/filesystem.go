package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Visitor interface {
	visitFile(file *File)
	visitDirectory(directory *Directory)
}

type SizeVisitor struct {
	smallDirSizeSum int64
}

func newSizeVisitor() *SizeVisitor {
	return &SizeVisitor{0}
}

func (sv *SizeVisitor) visitFile(file *File) {}

func (sv *SizeVisitor) visitDirectory(directory *Directory) {
	var childrenSizeSum int64 = 0

	for _, child := range directory.children {
		child.accept(sv)
		childrenSizeSum += child.getSize()
	}

	directory.size = childrenSizeSum
	if directory.size <= 100000 {
		sv.smallDirSizeSum += directory.size
	}
}

type ToDeleteVisitor struct {
	freeNeeded int64
	toDelete   *Directory
}

func newToDeleteVisitor(freeNeeded int64) *ToDeleteVisitor {
	return &ToDeleteVisitor{freeNeeded, nil}
}

func (tdv *ToDeleteVisitor) visitFile(file *File) {}

func (tdv *ToDeleteVisitor) visitDirectory(directory *Directory) {
	if directory.size >= tdv.freeNeeded &&
		(tdv.toDelete == nil || directory.size < tdv.toDelete.size) {
		tdv.toDelete = directory
	}

	for _, child := range directory.children {
		child.accept(tdv)
	}
}

type FileSystemNode interface {
	accept(visitor Visitor)
	getSize() int64
}

type FileSystemNodeBase struct {
	name   string
	size   int64
	parent *Directory
}

func (n *FileSystemNodeBase) getSize() int64 {
	return n.size
}

type File struct {
	FileSystemNodeBase
}

func (f *File) accept(vistor Visitor) {
	vistor.visitFile(f)
}

func newFile(name string, size int64, parent *Directory) *File {
	return &File{
		FileSystemNodeBase: FileSystemNodeBase{
			name:   name,
			size:   size,
			parent: parent,
		},
	}
}

type Directory struct {
	FileSystemNodeBase
	children map[string]FileSystemNode
}

func (d *Directory) accept(vistor Visitor) {
	vistor.visitDirectory(d)
}

func newDirectory(name string, parent *Directory) *Directory {
	return &Directory{
		FileSystemNodeBase: FileSystemNodeBase{
			name:   name,
			size:   0,
			parent: parent,
		},
		children: map[string]FileSystemNode{},
	}
}

type Command interface {
	parse(scanner *bufio.Scanner) bool
	run(current *Directory, root *Directory) *Directory
}

type CdCommand struct {
	arg string
}

func (c *CdCommand) parse(scanner *bufio.Scanner) bool {
	c.arg = strings.Split(scanner.Text(), " ")[2]
	return scanner.Scan()
}

func (c *CdCommand) run(current *Directory, root *Directory) *Directory {
	switch c.arg {
	case "..":
		return current.parent
	case "/":
		return root
	default:
		return current.children[c.arg].(*Directory)
	}
}

type LsCommand struct {
	output []string
}

func (l *LsCommand) parse(scanner *bufio.Scanner) bool {
	for scanner.Scan() {
		if scanner.Text()[0] == '$' {
			return true
		}

		l.output = append(l.output, scanner.Text())
	}

	return false
}

func (l *LsCommand) run(current *Directory, root *Directory) *Directory {
	for _, line := range l.output {
		splitted := strings.Split(line, " ")
		name := splitted[1]
		var child FileSystemNode

		if splitted[0] == "dir" {
			child = newDirectory(name, current)
		} else {
			size, err := strconv.ParseInt(splitted[0], 10, 64)
			errorHandler(err)
			child = newFile(name, size, current)
		}

		current.children[name] = child
	}

	return current
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	commandsMap := map[string](func() Command){
		"cd": func() Command { return &CdCommand{} },
		"ls": func() Command { return &LsCommand{} },
	}

	rootDir := newDirectory("/", nil)
	currentDir := rootDir
	moreCommands := true
	for moreCommands {
		command := commandsMap[strings.Split(scanner.Text(), " ")[1]]()
		moreCommands = command.parse(scanner)
		currentDir = command.run(currentDir, rootDir)
	}

	sizeVisitor := newSizeVisitor()
	sizeVisitor.visitDirectory(rootDir)
	fmt.Println(sizeVisitor.smallDirSizeSum)

	freeNeeded := 30000000 - (70000000 - rootDir.size)
	toDeleteVisitor := newToDeleteVisitor(freeNeeded)
	toDeleteVisitor.visitDirectory(rootDir)
	fmt.Println(toDeleteVisitor.toDelete.size)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
