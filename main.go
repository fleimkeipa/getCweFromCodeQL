package main

func main(){
  	makeFilter()
}

type cweStruct struct {
	Name  string
	ID    string
	Count int
}

type language struct {
	name    string
	filter  string
	display string
}

var lang = []language{
	{
		name:    "cpp",
		filter:  "C++",
		display: "Cpp",
	},
	{
		name:    "csharp",
		filter:  "C#",
		display: "Cs",
	},
	{
		name:    "go",
		filter:  "Default",
		display: "Golang",
	},
	{
		name:    "java",
		filter:  "Java",
		display: "Java",
	},
	{
		name:    "javascript",
		filter:  "JavaScript",
		display: "JavaScript",
	},
	{
		name:    "python",
		filter:  "Python",
		display: "Python",
	},
	{
		name:    "ruby",
		filter:  "Default",
		display: "Ruby",
	},
}

func makeFilter() {
	f, err := os.Create("cwes.go")
	if err != nil {
		fmt.Println(err)
		return
	}

	f.WriteString("package main\nvar(")

	for _, v := range lang {

		blogTitles, err := GetAllCWE(v)
		f.WriteString(fmt.Sprintf("\n%s=map[string][]string{\n", v.display))

		if err != nil {
			log.Println(err)
		}

		all := []string{}
		for i, v := range blogTitles {
			for _, v2 := range blogTitles[i+1:] {
				if v.Name == v2.Name {
					blogTitles[i].ID = fmt.Sprintf("%s,%s", blogTitles[i].ID, v2.ID)
					all = append(all, v.ID)
					v2.ID = ""
					blogTitles[i].Count++
				}
			}
		}

		var allCwes string
		max := 0
		for i, v := range blogTitles {
			for _, v2 := range blogTitles {
				if v.Name == v2.Name {
					if v2.Count > max {
						max = v2.Count
					}
				}
			}

			for _, v2 := range blogTitles[i:] {
				if v.Name == v2.Name && v2.Count == max {
					v2.ID = strings.ReplaceAll(v2.ID, ",C", "\",\"C")
					v2.ID = strings.ReplaceAll(v2.ID, "‑", "-")
					allCwes += fmt.Sprintf("\"%s\" :{\"%s\"},\n", v2.Name, v2.ID)
					break
				}
				if v.Name == v2.Name && len(v2.ID) < 6 {
					allCwes += fmt.Sprintf("\"%s\" :{\"%s\"},\n", v2.Name, v2.ID)
					break
				}
			}

			max = 0
		}

		f.WriteString(allCwes)
		f.WriteString("\n}")
	}
	f.WriteString("\n)")
}

func GetAllCWE(lang language) ([]cweStruct, error) {
	url := fmt.Sprintf("https://codeql.github.com/codeql-query-help/%s-cwe/", lang.name)
	// Get the HTML
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var cwes []cweStruct
	var cwe cweStruct
	counter := 0
	doc.Find("tr td").Each(func(i int, s *goquery.Selection) {
		if s.Text() != lang.filter {
			if counter == 0 {
				cwe.ID = s.Text()
				counter++
			} else if counter == 1 {
				cwe.Name = s.Text()
				cwes = append(cwes, cwe)
				counter++
			} else if counter == 2 {
				counter = 0
			}
		}
	})
	return cwes, nil
}
