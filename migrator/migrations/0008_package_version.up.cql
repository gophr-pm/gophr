----------------------------------- PACKAGES -----------------------------------

ALTER TABLE gophr.packages DROP versions;

------------------------------- PACKAGE VERSIONS -------------------------------

CREATE TABLE IF NOT EXISTS package_versions (
  author text,
  repo text,
  sha text,
  version text,
  date_discovered timestamp,
  PRIMARY KEY ((author, repo), date_discovered)
) WITH CLUSTERING ORDER BY (date_discovered DESC);
