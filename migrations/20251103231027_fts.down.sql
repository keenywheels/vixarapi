DROP INDEX IF EXISTS mv_token_search_pk;
DROP INDEX IF EXISTS mv_token_search_trgm_idx;
DROP INDEX IF EXISTS mv_token_search_interest_idx;

DROP MATERIALIZED VIEW IF EXISTS mv_token_search;

DROP EXTENSION IF EXISTS pg_trgm;
