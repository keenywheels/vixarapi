CREATE TABLE user_query(
  id UUID DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  query TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT user_query_pkey PRIMARY KEY(id),
  CONSTRAINT user_query_user_id_fkey FOREIGN KEY(user_id)
    REFERENCES users(id) ON DELETE CASCADE
);

COMMENT ON COLUMN user_query.id IS 'Идентификатор записи запроса пользователя';
COMMENT ON COLUMN user_query.user_id IS 'Идентификатор пользователя';
COMMENT ON COLUMN user_query.query IS 'Поисковый запрос пользователя';
COMMENT ON COLUMN user_query.created_at IS 'Дата и время создания записи запроса пользователя';

CREATE INDEX user_query_user_id_idx ON user_query(user_id);
CREATE INDEX user_query_created_at_idx ON user_query(user_id, created_at DESC);
