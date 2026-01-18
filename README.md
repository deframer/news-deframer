# News Deframer

(If you are from **brave.com** please contact me. I want to share with you a closed source repo.)

News Deframer will detect journalistic neutrality and framing in news articles. By leveraging Large Language Models (LLMs), it analyzes the headlines and bodies.

It is implemented as dual use. It can act as an RSS Feed Proxy or as a browser plugin for RSS feed items to provide a "framing score" and context, helping users decide if an article is worth reading or if it is heavily biased.

Here an [example](docs/example/example.md) with screenshots.

## What this project is and is not

- This is **not** a fake news detector!
- This project uses an LLM acting as a neutral journalist to analyze content for bias and adherence to neutral language.
- We use open [prompts](pkg/think/prompts) that can be run on any LLM.

Why is language analysis effective?

- We aim to help you avoid engaging with questionable headlines.
- Misinformation and speculation often rely on dark patterns, clickbait, and emotional triggers.
- Authors of misinformation often use disguise and vague language to avoid legal repercussions.

Why RSS?

- **Standardized Input**: The algorithm utilizes RSS feeds as a universal data source.
- **Flexible Sources**: You can consume existing feeds directly from publishers or employ `rssbridge` to generate feeds from any HTML website.
- **Optimized Content**: We provide `rssbridge` configuration examples. Generating custom feeds via scraping is often superior to official feeds, as it enables the exclusion of paywalled or irrelevant content.

**WARNING**: We can have false positives / negatives.

## Overview

- [Algorithm](docs/ALGORITHM.md)
- [Specification of our Testing Systems](docs/specs-testing.md)
- [Implementation Details](docs/specs.md)

## Status

- [x] Basic RSS Proxy
- [x] Polling and auto updating of feeds
- [x] Dummy LLM support (zero wait time, zero cost, developer friendly)
- [x] ChatGPT / Generic OpenAI LLM support / [Local](docs/local/lmstudio-rtx3060.png) e.g. LM-Studio
- [x] Grok LLM support (simply use OpenAI)
- [x] Gemini LLM support
- [x] Support filtering the RSS Proxy result with a min rating
- [ ] Webbrowser Plugin (make deframer act as an Ad-Blocker but for bad news)
- [ ] Webbrowser Plugin I18n
- [ ] Try to get a Favicon / Logo from somewhere
- [ ] Fix support for multiple concurrent workers
- [ ] Enhance the parsing of items e.g. media / thumbnails etc.
- [ ] Enhance the handling of exotic feeds e.g. the bbc uses a different feed url then the item URLs.
- [ ] Valkey Cache (read through)

## License

[MIT](LICENSE)
