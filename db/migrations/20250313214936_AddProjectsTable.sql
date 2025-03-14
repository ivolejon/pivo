-- migrate:up
CREATE TABLE projects (
  id UUID NOT NULL,
  client_id UUID NOT NULL,
  title TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (id, created_at)
);

-- migrate:down
DROP TABLE projects
