CREATE TABLE user_token_sub
(
    id            UUID                 DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL,
    token         TEXT        NOT NULL,
    category      TEXT        NOT NULL,
    curr_interest NUMERIC     NOT NULL,
    prv_interest  NUMERIC     NOT NULL DEFAULT 0,
    threshold     NUMERIC     NOT NULL,
    method        TEXT        NOT NULL DEFAULT 'denormalized',
    scan_date     TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT user_token_sub_pk PRIMARY KEY (id),
    CONSTRAINT user_token_sub_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT user_token_sub_method_check
        CHECK (method IN ('denormalized', 'global_median', 'category_median'))
);

COMMENT ON COLUMN user_token_sub.id IS 'Идентификатор подписки пользователя на токен';
COMMENT ON COLUMN user_token_sub.user_id IS 'Идентификатор подписки пользователя';
COMMENT ON COLUMN user_token_sub.token IS 'Токен, на который подписываемся';
COMMENT ON COLUMN user_token_sub.category IS 'Категория токена';
COMMENT ON COLUMN user_token_sub.curr_interest IS 'Текущий интерес по токену';
COMMENT ON COLUMN user_token_sub.prv_interest IS 'Предыдущее значение интереса по токену';
COMMENT ON COLUMN user_token_sub.threshold IS 'Пороговое значение для уведомления';
COMMENT ON COLUMN user_token_sub.method IS 'Метод для сравнение значений интереса';
COMMENT ON COLUMN user_token_sub.scan_date IS 'Дата и время последнего обновления';
COMMENT ON COLUMN user_token_sub.created_at IS 'Дата и время создания записи';

CREATE UNIQUE INDEX user_token_sub_unique_idx ON user_token_sub (user_id, category, token, method);
CREATE INDEX user_token_sub_user_id_idx ON user_token_sub (user_id);
CREATE INDEX user_token_sub_scan_date_idx ON user_token_sub (scan_date);
