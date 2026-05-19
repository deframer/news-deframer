# Languages

## RSS Feed Languages

We came up with this convention because some feeds provide improper or no language metadata.

We treat RSS feed language data as follows:

- If `feed.language` is already set in the database, we use that value.
- If `feed.language` is empty and the XML feed declares a language, we store the XML language in the database and use it.
- If both the database value and the XML language are empty, we store `en` as the default and use it.
- If the database value and XML language both exist but differ, we ignore the XML value and keep the database value.
- Empty imported language values are treated as unset.
