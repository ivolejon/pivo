-- migrate:up
CREATE TABLE documents (
  id UUID NOT NULL,
  embeddings_ids UUID[] NOT NULL,
  filename TEXT NOT NULL,
  title TEXT,
  project_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (id)
)
-- migrate:down
DROP TABLE documents
