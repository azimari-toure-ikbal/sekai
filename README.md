# Sekai 🌏

**Sekai** (meaning "world" in Japanese) is a powerful localization tool that simplifies generating translation files for your app. With just one command, Sekai crawls your project to find text strings and automatically creates both the original and translated files.

## Features

- **Automatic Crawling**: No need to feed it any information. Sekai will explore your project, locating all text strings within your files.
- **Localization Generation**: Generates both original and translated localization files effortlessly.
- **Current Framework**: Tested with Next.js `.jsx` and `.tsx`. Future support for additional frameworks and platforms, starting with Flutter and WordPress, is planned.
- **Current Languages**: English, French, German, Spanish, and Portuguese.

## Output Structure

Sekai generates a JSON file where the data is organized as key-value pairs in the following format:

```json
{
  "path.to.component/page": "string value {{variableX}}"
}
```

Each key represents the path to a component or page, and each value is the translated string, supporting dynamic variables wrapped in double curly braces (e.g., `{{variableX}}`).

## Prerequisites

Sekai relies on **Ollama** to handle translation requests through **Llama3.2**. Ensure both are installed on your system before running Sekai.

## Installation

1. Clone this repository.
2. Build the project using `make build`.

## Usage

1. Navigate to the root of your project.
2. In your CLI, run:

   ```bash
   <path_to_sekai_bin> -env yourenv -i inputLang -o outputLang
   ```

   - `yourenv`: The framework environment in which Sekai is operating.
   - `inputLang`: The language of the original text.
   - `outputLang`: The target language for translation.

### Example

```bash
sekai -env nextjs -i en -o fr
```

This command translates your NextJS app from English to French.

## Roadmap

- **Bring your own model**: Allow users to specify their own models to be used for translation.
- **Additional Frameworks**: Support for frameworks like Flutter and WordPress is on the way.
- **More Platforms**: Expansion to other popular app and web platforms in the future.

---
