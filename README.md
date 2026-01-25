# News Deframer

~~(If you are from **brave.com** please contact me. I want to share with you a closed source repo.)~~ Thanks everybody for helping 💕.

News Deframer will detect journalistic neutrality and framing in news articles. By leveraging Large Language Models (LLMs), it analyzes the headlines and bodies.

It is implemented as dual use. It can act as an **RSS Feed Proxy** or as a **browser plugin** to provide a "framing score" and context, helping users decide if an article is worth reading or if it is heavily biased.

Here are examples with screenshots for our [browser plugin](https://deframer.github.io/screenshots/) and [RSS feed](https://deframer.github.io/screenshots/rss-feed/) for very quick use on your mobile phone.

## What this project is and is not

- This is **not** a fake news detector!
- This project uses an LLM acting as a neutral journalist to analyze content for bias and adherence to neutral language.
- We use open transparent [prompts](pkg/think/prompts) that can be run on any LLM.

Why is language analysis effective?

- We aim to help you avoid engaging with questionable headlines.
- Misinformation and speculation often rely on dark patterns, clickbait, and emotional triggers.
- Authors of misinformation often use disguise and vague language to avoid legal repercussions.

**WARNING**: We can have false positives / negatives.

## Installation

Set up your own [instance](docs/SETUP.md).

## Overview

- [Algorithm](docs/ALGORITHM.md)
- [Specification of our Testing Systems](docs/specs-testing.md)
- [Implementation Details](docs/specs.md)
- [TODO List / Status](https://deframer.github.io/de/todo/)

## License

[MIT](LICENSE)
