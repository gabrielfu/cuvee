package vintagecharts

import (
	"cuvee/domain/wines"
	"cuvee/external/llm"
	"cuvee/external/search"
	"fmt"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile(`Region: ([0-9]+)`)

func searchAboutWine(
	name string,
	searchEngine search.SearchEngine,
) string {
	query := fmt.Sprintf("%s wine region location", name)
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

// PickRegion picks a region for the given wine within the list of regions.
// Output is either an empty string or a value in the regions list.
func PickRegion(
	wine *wines.Wine,
	regions []string,
	llmObj llm.LLM,
	searchEngine search.SearchEngine,
) (string, error) {
	messages := []llm.Message{
		llm.SystemMessage("You are a wine expert. You are very familiar with different types of wines and different wine regions around the world."),
	}

	prompt := fmt.Sprintf(`Consider this wine: 
Name: %s %s %s %s
`, wine.Name, wine.Vintage, wine.Country, wine.Region)

	if extra := searchAboutWine(wine.Name, searchEngine); extra != "" {
		prompt += fmt.Sprintf(
			"\nHere is some information that I found online. They may or may not be useful.\n```\n%s\n```\n",
			extra,
		)
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
