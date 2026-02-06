SET duckdb.force_execution = true;

SELECT count(*) AS trend_doc_count
FROM trends;
