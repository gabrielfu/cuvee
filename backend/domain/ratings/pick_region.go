package ratings

import (
	"cuvee/external/llm"
	"cuvee/external/search"
	"fmt"
	"regexp"
	"strconv"
)

type WineInfo struct {
	Name    string
	Vintage string
	Country string
	Region  string
}

var re = regexp.MustCompile(`Region: ([0-9]+)`)

func searchAboutWine(
	searchEngine search.SearchEngine,
	target WineInfo,
) string {
	query := fmt.Sprintf("%s wine region location", target.Name)
	searchResponse, err := searchEngine.WebSearch(query, nil)
	if err != nil {
		return ""
	}

	if len(searchResponse.Items) == 0 {
		return ""
	}

	snippets := ""
	for _, item := range searchResponse.Items[:min(len(searchResponse.Items), 5)] {
		snippet := item.Snippet
		if item.Desc != "" {
			snippet = item.Desc
		}
		snippets += fmt.Sprintf("- %s\n", snippet)
	}

	return snippets
}

func PickRegion(
	llmObj llm.LLM,
	searchEngine search.SearchEngine,
	target WineInfo,
	regions []string,
) (string, error) {
	messages := []llm.Message{
		llm.SystemMessage("You are a wine expert. You are very familiar with different types of wines and different wine regions around the world."),
	}

	prompt := fmt.Sprintf(`Consider this wine: 
Name: %s %s %s %s
`, target.Name, target.Vintage, target.Country, target.Region)

	if extra := searchAboutWine(searchEngine, target); extra != "" {
		prompt += fmt.Sprintf(`
Here is some information that I found online. They may or may not be useful.
%s
`, extra)
	}

	regionPrompt := ""
	regionMapping := make(map[string]string)
	for i, r := range regions {
		regionPrompt += fmt.Sprintf("(%d) %s\n", i+1, r)
		regionMapping[strconv.Itoa(i+1)] = r
	}
	prompt += fmt.Sprintf(`
You must choose one of the following regions as the producing region of the wine.
If you are not sure, please choose OTHERS.

(0) OTHERS
%s

IMPORTANT
---------
You must answer in the following format, without any additional information.
Region: number
---------
You must not answer a region that is not in the list.
If none of the above regions is correct, please answer
Region: 0
`, regionPrompt)
	messages = append(messages, llm.UserMessage(prompt))

	resp, err := llmObj.Chat(messages)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(resp.Content)
	if matches == nil || len(matches) < 2 || matches[1] == "0" {
		return "", nil
	}

	choice := matches[1]
	return regionMapping[choice], nil
}
