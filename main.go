package main

func main(){
  makeFilter()
}

func makeFilter() {
	blogTitles, err := GetSans25()
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
				fmt.Printf("\"%s\" :{\"%s\"},\n", v2.Name, v2.ID)
				break
			}
			if v.Name == v2.Name && len(v2.ID) < 6 {
				fmt.Printf("\"%s\" :{\"%s\"},\n", v2.Name, v2.ID)
				break
			}
		}
		max = 0
	}
}


func GetSans25() ([]cweStruct, error) {
	url := "https://codeql.github.com/codeql-query-help/cpp-cwe/"
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
		if s.Text() != "C++" {
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
