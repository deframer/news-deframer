# Webbrowser Plugin Example

## AI used for these tests

- Mac Mini M1 (16GB RAM), Tahoe 26.2
- [LM Studio](https://lmstudio.ai/) 0.3.39
- [meta-llama-3.1-8b-instruct](https://huggingface.co/lmstudio-community/Meta-Llama-3.1-8B-Instruct-GGUF)

Using these [Prompts](../../pkg/think/prompts/)

**Quote:** *"You are a strictly objective, neutral media analyst and news editor. Your task is to analyze texts (title and description) for journalistic quality, bias, and sensationalism, and to create a neutral version. You must provide specific, concise reasons for every assessment and correction, culminating in an overall summary. [...]"*

The prompt instructs the AI to be objective, neutral, and strictly politically unbiased.

The LLM is slow on the M1 (but also free, efficient and cheap).

```txt
level=DEBUG msg="openai request duration" duration=26.833331235s
level=DEBUG msg="openai token usage" prompt_tokens=1470 completion_tokens=303 thoughts_tokens=0 total_tokens=1773
```

## Examples

### BBC

- RSS feed `https://feeds.bbci.co.uk/news/world/rss.xml`


**Portal**

We replace the entire web page. Press "Hide" to see the original page.

Please note - you never see the title or description of the original page. You always see neutral titles and descriptions converted by the AI to >not< spike emotional reactions.

The percentage shows, how much someone is messing with you. Red = bad.

<img src="bbc-portal.jpg" alt="BBC Portal" width="600"/>

Hovering over the percentage bar displays the overall news deframer result for that news item. Usually at this point you don't have to click on the article.

<img src="portal-hover.png" alt="Portal Hover" width="200"/>


**Article**

Hide closes the deframer overlay and shows you the article. We don't show the original title or description to protect your mental health. You have the freedom to expand the view and see all the details. "Hide" closes the deframer view and opens the original article.

**NOTE**: The article itself was never touched by the deframer. We still use title and description of the RSS feed. This might change in the future by using trusted flaggers to trigger a full scan.

Article 1

| Basic | Expanded |
| :---: | :---: |
| <img src="bbc-article1-basic.png" alt="BBC Article 1 - Basic" width="600"/> | <img src="bbc-article1-expanded.png" alt="BBC Article 1 - Expanded" width="600"/> |

Article 2

| Basic | Expanded |
| :---: | :---: |
| <img src="bbc-article2-basic.png" alt="BBC Article 1 - Basic" width="600"/> | <img src="bbc-article2-expanded.png" alt="BBC Article 1 - Expanded" width="600"/> |

