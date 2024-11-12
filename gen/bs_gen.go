package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func genBasic(sb *strings.Builder) {
        sb.WriteString("package main\n")
        sb.WriteString("import (\"math/rand\"; \"fmt\")\n")
	sb.WriteString("var bsc = map[string][]string{\n")
	files, err := os.ReadDir("collections")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			// Read the contents of the file
			f, err := os.Open("collections/" + fileName)
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
func compositeVariaablesTransform(s *[]string) []string {
        var r []string
        for i, v := range *s {
                if v[0] == '$' {
                        r = append(r, v)
                        (*s)[i] = "%s"
                }
                
        }
        return r 
}
func genComposite(sb *strings.Builder) {

	files, err := os.ReadDir("composite_collections")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	sb.WriteString("var bsc_composite = map[string]func() string{\n")
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			f, err := os.Open("composite_collections/" + fileName)
                        data, _ := io.ReadAll(f)

			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

                        sb.WriteString("\""+fileName+"\": func() string {\n")
                        s := strings.Fields(string(data))
                        cv := compositeVariaablesTransform(&s)
                        for i, v := range cv {
                                sb.WriteString(fmt.Sprintf("\tr_%d := rand.Intn(len(bsc[\"%s\"]))\n",i,v[1:]))
                        }
                        sb.WriteString(fmt.Sprintf("\tr := fmt.Sprintf("))
                        sb.WriteString(strconv.Quote(strings.Join(s, " ")))
                        for i, v := range cv {
                                sb.WriteString(fmt.Sprintf(",bsc[\"%s\"][r_%d]",v[1:],i))
                        }
                        sb.WriteString(")\nreturn r")
		}
                        sb.WriteString("}")
	}

	sb.WriteString("}\n")
}

func main() {
	var sb strings.Builder
	genBasic(&sb)
        genComposite(&sb)
        f, e := os.OpenFile("bsc.go", os.O_CREATE | os.O_WRONLY, 0644)
        if e != nil {
                fmt.Print(e)
        }
        _, e = f.Write([]byte(sb.String()))
        if e != nil {
                fmt.Print(e)
        }
        f.Close()
}
