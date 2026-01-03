-- add category column to token_data table
ALTER TABLE token_data
ADD COLUMN category TEXT NOT NULL DEFAULT 'other';

-- drop previous search mv (copy fts down migration)
DROP INDEX IF EXISTS mv_token_search_pk;
DROP INDEX IF EXISTS mv_token_search_trgm_idx;
DROP INDEX IF EXISTS mv_token_search_interest_idx;
DROP MATERIALIZED VIEW IF EXISTS mv_token_search;

-- create new mv (based on fts up)
CREATE MATERIALIZED VIEW mv_token_search AS
WITH
    aggr AS (SELECT token_name,
                    scrape_date,
                    category,
                    SUM(interest)                   AS interest,
                    ROUND(AVG(sentiment))::SMALLINT AS sentiment
             FROM token_data
             GROUP BY (token_name, scrape_date, category)),
    global_medians AS (SELECT scrape_date,
                              PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY interest) AS median_interest
                       FROM aggr
                       GROUP BY scrape_date),
    category_medians AS (SELECT scrape_date,
                                category,
                                PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY interest) AS median_interest
                         FROM aggr
                         GROUP BY (scrape_date, category))
SELECT a.token_name,
       a.scrape_date,
       a.interest,
       a.sentiment,
       a.category,
       gm.median_interest AS global_median,
       cm.median_interest AS category_median
FROM aggr a
         JOIN global_medians gm ON a.scrape_date = gm.scrape_date
         JOIN category_medians cm ON a.scrape_date = cm.scrape_date AND a.category = cm.category;

CREATE UNIQUE INDEX mv_token_search_pk ON mv_token_search (token_name, scrape_date, category);
CREATE INDEX mv_token_search_trgm_idx ON mv_token_search USING GIN (token_name gin_trgm_ops);
CREATE INDEX mv_token_search_interest_idx ON mv_token_search (interest DESC);
CREATE INDEX mv_token_search_category_idx ON mv_token_search (category);

-- NOTE: after this migration, you also should replace default category values based on site name:
-- UPDATE token_data
-- SET category = 'category_name'
-- WHERE site_name IN ('site1', 'site2');
--
-- after replacing default values, update mv
