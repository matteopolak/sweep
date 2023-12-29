package main

// import bubbletea
import (
	"bufio"
	_ "embed"
	"fmt"
	"io/fs"
	"os"

	sw "github.com/matteopolak/sweep/pkg/api"
)

//go:embed sweep.toml
var defaultConfig string

// main is the entrypoint for the program
func main() {
	// first, get the current working directory
	cwd, err := os.Getwd()

	// if there's an error, print it and exit
	if err != nil {
		fmt.Printf("could not get current working directory: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("using root directory: ", cwd)

	// read sweep.toml, otherwise use the default config
	rawContent, err := fs.ReadFile(os.DirFS(cwd), "sweep.toml")
	content := ""

	if err != nil {
		content = defaultConfig

		fmt.Println("using default config")
	} else {
		content = string(rawContent)

		fmt.Println("using config from sweep.toml")
	}

	// parse the config
	config, err := sw.ParseConfig(content)

	// if there's an error, print it and exit
	if err != nil {
		fmt.Printf("could not parse config: %s\n", err)
		os.Exit(1)
	}

	// then, get a new iterator for the current working directory
	// with the default predicates
	iterator := sw.NewDirIterator(cwd, config.Predicates, []string{".git"})

	// then, iterate over the iterator
	for {
		// get the next folder
		folder, err := iterator.Next()

		// if there's an error, print it and exit
		if err != nil {
			fmt.Printf("could not get next folder: %s\n", err)
			os.Exit(1)
		}

		// if the folder is empty, we're done
		if folder == "" {
			break
		}

		// prompt user if they want to delete the folder
		// also provide the size of the folder
		size := sw.GetDirSize(folder)

		fmt.Printf("(%s) delete %s? [Y/n] ", FormatBytes(size), folder)

		// read the user's input
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		// if the user wants to delete the folder, delete it
		if text == "y\n" || text == "y\r\n" || text == "\n" || text == "\r\n" {
			err := os.RemoveAll(folder)

			if err != nil {
				fmt.Printf("could not delete folder: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("deleted %s\n", folder)
		}
	}
}

func FormatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	bytes = bytes / 1024

	if bytes < 1024 {
		return fmt.Sprintf("%d KB", bytes)
	}

	bytes = bytes / 1024

	if bytes < 1024 {
		return fmt.Sprintf("%d MB", bytes)
	}

	bytes = bytes / 1024

	if bytes < 1024 {
		return fmt.Sprintf("%d GB", bytes)
	}

	bytes = bytes / 1024

	if bytes < 1024 {
		return fmt.Sprintf("%d TB", bytes)
	}

	bytes = bytes / 1024

	return fmt.Sprintf("%d PB", bytes)
}
