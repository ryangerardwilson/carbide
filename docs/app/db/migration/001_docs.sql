CREATE TABLE IF NOT EXISTS docs_checks (
  id bigserial PRIMARY KEY,
  checked_at timestamptz NOT NULL DEFAULT now(),
  service text NOT NULL DEFAULT 'carbide-docs'
);
