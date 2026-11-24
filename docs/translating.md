# Translating Quotes

Quotes ships six complete interface locales: English (`en`), Spanish (`es`), French (`fr`), Arabic (`ar`), German (`de`), and Italian (`it`). English is the source locale. All server and CLI binaries embed the same locale files from `src/common/i18n/locales/`.

| Code | Language | Native name | Direction |
|---|---|---|---|
| `en` | English | English | LTR |
| `es` | Spanish | Español | LTR |
| `fr` | French | Français | LTR |
| `ar` | Arabic | العربية | RTL |
| `de` | German | Deutsch | LTR |
| `it` | Italian | Italiano | LTR |

The language registry and embedded locale files are the source of supported-language names, direction, and plural categories. The browser resolves a one-time `?lang=` selection first, then the `lang` cookie when the device has accepted preference cookies, then `Accept-Language`, and finally English. A visitor can choose a language in `/server/preferences`; saving it stores a cookie only after the visitor accepts that cookie category. Arabic pages use right-to-left document direction.

## Contributing a translation

1. Edit one locale file in `src/common/i18n/locales/`. Keep its JSON structure and key order aligned with `en.json`.
2. Translate every value; never translate keys, product names, route names, API fields, configuration keys, code, URLs, or literal values such as `dark` and `auto`.
3. Preserve every `{placeholder}` exactly, including its spelling and braces. Use complete sentences so word order can vary naturally.
4. Use the language's required CLDR plural categories if translating a plural message. Do not copy English plural categories blindly.
5. Use a consistent, neutral register and the glossary below. Respect local punctuation and typography.
6. Update the locale's `meta.translators` credit when the contributor wants to be listed.
7. Review the page at 320, 768, and 1280 pixels. Check Arabic in a mirrored layout and confirm longer text does not clip or overflow.

The Docker-backed test suite checks that locale files are embedded, every locale has the same message keys as English, no translated value is empty, and interpolation placeholders match. Run the repository's documented `make test` target when asked to verify a translation change.

## Glossary

| English | Translation guidance |
|---|---|
| quote | A short quotation or saying; use one consistent term per language. |
| collection | A named group of quotes. |
| preference | A browser display choice, not an account setting. |
| sync | Transfer browser preferences between devices. |
| token | Keep the technical credential meaning; do not translate the token value. |
| MCP | Keep the protocol acronym unchanged. |
| API | Keep the acronym unchanged; translate surrounding text. |

## Locale file rules

- UTF-8 JSON, two-space indentation, no comments or trailing commas.
- English is the canonical key set. Never add a key only to another language.
- A missing translation falls back to English at runtime; missing keys are still defects and must be fixed in the locale file.
- Supported-language count in this page and the README must match the locale files and the runtime registry.
- New languages require a locale file, registry entry (native name, direction, and CLDR plural categories), matching CLI and browser support, and documentation updates in the same change.

## Credit

No individual translator credits are currently recorded. Contributors may add their preferred name or handle to the locale metadata and open a pull request with the affected locale file.
