CREATE TABLE IF NOT EXISTS deploy_checks (
  id bigserial PRIMARY KEY,
  checked_at timestamptz NOT NULL DEFAULT now(),
  service text NOT NULL DEFAULT 'carbide-docs'
);
