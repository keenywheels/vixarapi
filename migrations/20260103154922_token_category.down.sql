-- delete new mv
DROP INDEX IF EXISTS mv_token_search_pk;
DROP INDEX IF EXISTS mv_token_search_trgm_idx;
DROP INDEX IF EXISTS mv_token_search_interest_idx;
DROP INDEX IF EXISTS mv_token_search_category_idx;
DROP MATERIALIZED VIEW IF EXISTS mv_token_search;

-- delete category column
ALTER TABLE token_data
DROP COLUMN category;

-- recreate previous mv
CREATE MATERIALIZED VIEW mv_token_search AS
WITH
    aggr AS (SELECT token_name,
                    scrape_date,
                    SUM(interest)                   AS interest,
                    ROUND(AVG(sentiment))::SMALLINT AS sentiment
             FROM token_data
             GROUP BY (token_name, scrape_date)),
    medians AS (SELECT scrape_date,
                       PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY interest) AS median_interest
                FROM aggr
                GROUP BY scrape_date)
SELECT a.token_name,
       a.scrape_date,
       a.interest,
       a.sentiment,
       m.median_interest
FROM aggr a
         JOIN medians m ON a.scrape_date = m.scrape_date;

CREATE UNIQUE INDEX mv_token_search_pk ON mv_token_search (token_name, scrape_date);
CREATE INDEX mv_token_search_trgm_idx ON mv_token_search USING GIN (token_name gin_trgm_ops);
CREATE INDEX mv_token_search_interest_idx ON mv_token_search (interest DESC);
