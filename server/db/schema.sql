CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE todos2 (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0
	);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20250131222727');
