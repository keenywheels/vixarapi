CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE MATERIALIZED VIEW mv_token_search AS
SELECT
	token_name,
  scrape_date,
  SUM(interest) AS interest,
  ROUND(AVG(sentiment))::SMALLINT AS sentiment,
  MAX(SUM(interest)) OVER () AS max_interest
FROM token_data
GROUP BY (token_name, scrape_date);

CREATE UNIQUE INDEX mv_token_search_pk ON mv_token_search (token_name, scrape_date);
CREATE INDEX mv_token_search_trgm_idx ON mv_token_search USING GIN (token_name gin_trgm_ops);
CREATE INDEX mv_token_search_interest_idx ON mv_token_search (interest DESC);