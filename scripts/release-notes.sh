#!/bin/bash

QUERY="$(cat <<EOF
SELECT
  FORMAT("At %t, the product %s got the following update: %s",
    published_at,
    product_name,
    REPLACE(description, "\n", ". ")) as line
FROM \`bigquery-public-data.google_cloud_release_notes.release_notes\`
WHERE published_at >= DATE_ADD(CURRENT_DATE(), INTERVAL -30 DAY)
  AND lower(product_name) LIKE lower("%${1}%")
ORDER BY published_at ASC
LIMIT 10
EOF
)"
bq query --format=json --use_legacy_sql=false "${QUERY}" | jq -r .[].line > notes.txt

PROMPT="Summarize the following text in a few topics:

$(cat notes.txt)"

textbison "${PROMPT}" | jq -r .predictions[].content | tee notes.md
