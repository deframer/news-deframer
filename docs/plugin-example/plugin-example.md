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