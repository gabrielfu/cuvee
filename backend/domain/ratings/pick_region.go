package ratings

import (
	"cuvee/external/llm"
	"cuvee/external/search"
	"fmt"
	"log"
)

type WineInfo struct {
	Name    string
	Vintage string
	Country string
	Region  string
}

func searchAboutWine(
	searchEngine search.SearchEngine,
	target WineInfo,
) string {
	query := fmt.Sprintf("%s wine region location", target.Name)
	log.Println("Searching for:", query)
	searchResponse, err := search.Search[search.GoogleWebSearchResponse](searchEngine, query, search.WebSearchGoogleSearchParam)
	if err != nil {
		return ""
	}

	if len(searchResponse.Items) == 0 {
		return ""
	}

	snippets := ""
	for _, item := range searchResponse.Items[:min(len(searchResponse.Items), 5)] {
		var snippet string
		if len(item.PageMap.Metatags) > 0 {
			snippet = item.PageMap.Metatags[0].OGDescription
		}
		if snippet == "" {
			snippet = item.Snippet
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
	for _, r := range regions {
		regionPrompt += fmt.Sprintf("- %s\n", r)
	}
	prompt += fmt.Sprintf(`
You must choose one of the following regions as the producing region of the wine.
If you are not sure, please choose OTHERS.

- OTHERS
%s

IMPORTANT
---------
You must answer in the following format, without any additional information.
Region: <choice>
---------
You must not answer a region that is not in the list.
If none of the above regions is correct, please answer
Region: OTHERS
`, regionPrompt)
	messages = append(messages, llm.UserMessage(prompt))

	log.Println(messages)

	resp, err := llmObj.Chat(messages)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}
