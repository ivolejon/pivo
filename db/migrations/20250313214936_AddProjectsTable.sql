-- migrate:up
CREATE TABLE projects (
  id UUID NOT NULL,
  client_id UUID,
  title TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (id, created_at)
);

-- migrate:down
DROP TABLE projects
