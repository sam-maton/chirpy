-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens(
  token TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP,
  user_id UUID NOT NULL,
  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
      REFERENCES users(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd