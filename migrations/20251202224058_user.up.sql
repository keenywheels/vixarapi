CREATE TABLE users (
  id UUID DEFAULT gen_random_uuid(),
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  tguser TEXT UNIQUE,
  vkid INTEGER UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT id_pkey PRIMARY KEY(id)
);

COMMENT ON COLUMN users.id IS 'Идентификатор пользователя';
COMMENT ON COLUMN users.username IS 'Имя пользователя';
COMMENT ON COLUMN users.email IS 'Электронная почта пользователя';
COMMENT ON COLUMN users.vkid IS 'Идентификатор пользователя в vk';
COMMENT ON COLUMN users.created_at IS 'Дата и время создания пользователя';

CREATE INDEX users_email_idx ON users(email);
CREATE INDEX users_vkid_idx ON users(vkid) WHERE vkid IS NOT NULL;
