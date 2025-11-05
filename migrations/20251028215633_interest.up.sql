CREATE TABLE token_data (
  token_id BIGSERIAL,
  token_name TEXT NOT NULL,
  interest INT NOT NULL,
  site_name TEXT NOT NULL,
  scrape_date TIMESTAMP NOT NULL,
  context TEXT,

  CONSTRAINT token_id_pkey PRIMARY KEY(token_id)
);

COMMENT ON COLUMN token_data.token_id IS 'Идентификатор токена';
COMMENT ON COLUMN token_data.token_name IS 'Название токена';
COMMENT ON COLUMN token_data.interest IS 'Показатель интереса';
COMMENT ON COLUMN token_data.site_name IS 'Название сайта';
COMMENT ON COLUMN token_data.scrape_date IS 'Дата сбора данных';
COMMENT ON COLUMN token_data.context IS 'Контекст';

CREATE INDEX token_data_token_name_date_idx ON token_data USING btree (token_name, scrape_date);
CREATE INDEX token_data_token_name_site_name_date_idx ON token_data USING btree (token_name, site_name, scrape_date);
