package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func genBasic(sb *strings.Builder) {
        sb.WriteString("package main\n")
	sb.WriteString("var bsa = map[string][]string{\n")
	files, err := os.ReadDir("atom_collections")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			// Read the contents of the file
			f, err := os.Open("atom_collections/" + fileName)
                        data, _ := io.ReadAll(f)

			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

			lines := strings.Split(string(data), "\n")
			var lineList []string
			for _, line := range lines {
				// Skip empty lines
				if line == "" {
					continue
				}
				lineList = append(lineList, line)
			}
			valueStr := "\"" + strings.Join(lineList, "\",\"") + "\""
			sb.WriteString(fmt.Sprintf("    \"%s\":{%s},\n", fileName, valueStr))
		}
	}
	sb.WriteString("}\n")
}

func main() {
	var sb strings.Builder
	genBasic(&sb)
        f, e := os.OpenFile("bs_atoms.go", os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)
        if e != nil {
                fmt.Print(e)
        }
        _, e = f.Write([]byte(sb.String()))
        if e != nil {
                fmt.Print(e)
        }
        f.Close()
}
