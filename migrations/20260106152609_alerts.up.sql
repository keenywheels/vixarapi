CREATE TABLE user_token_sub
(
    id            UUID                 DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL,
    token         TEXT        NOT NULL,
    category      TEXT        NOT NULL,
    curr_interest BIGINT      NOT NULL,
    prv_interest  BIGINT      NOT NULL DEFAULT 0,
    threshold     NUMERIC     NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL,

    CONSTRAINT user_token_sub_pk PRIMARY KEY (id),
    CONSTRAINT user_token_sub_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON COLUMN user_token_sub.id IS 'Идентификатор подписки пользователя на токен';
COMMENT ON COLUMN user_token_sub.user_id IS 'Идентификатор подписки пользователя';
COMMENT ON COLUMN user_token_sub.token IS 'Токен, на который подписываемся';
COMMENT ON COLUMN user_token_sub.category IS 'Категория токена';
COMMENT ON COLUMN user_token_sub.curr_interest IS 'Текущий интерес по токену';
COMMENT ON COLUMN user_token_sub.prv_interest IS 'Предыдущее значение интереса по токену';
COMMENT ON COLUMN user_token_sub.threshold IS 'Пороговое значение для уведомления';
COMMENT ON COLUMN user_token_sub.updated_at IS 'Дата и время последнего обновления';

CREATE UNIQUE INDEX user_token_sub_token_user_idx ON user_token_sub (user_id, category, token);
CREATE INDEX user_token_sub_user_id_idx ON user_token_sub (user_id);
CREATE INDEX user_token_sub_updated_at_idx ON user_token_sub (updated_at);
