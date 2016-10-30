
-------------------------------- PACKAGES TABLE --------------------------------

CREATE TABLE IF NOT EXISTS packages (
  repo text,
  stars int,
  author text,
  awesome boolean,
  description text,
  trend_score float,
  search_blob text,
  date_discovered timestamp,
  daily_downloads bigint,
  weekly_downloads bigint,
  monthly_downloads bigint,
  date_last_indexed timestamp,
  all_time_downloads bigint,
  all_time_version_downloads map<text, bigint>,
  PRIMARY KEY (author, repo)
) WITH CLUSTERING ORDER BY (repo ASC);

CREATE CUSTOM INDEX IF NOT EXISTS packages_index
  ON packages ()
  USING 'com.stratio.cassandra.lucene.Index'
  WITH OPTIONS = {
    'refresh_seconds': '60',
    'schema': '{
      default_analyzer: "english",
      fields: {
        stars: {type: "integer"},
        awesome: {type: "boolean"},
        trend_score: {type: "float"},
        search_blob: {type: "text"},
        date_discovered: {type: "date"},
        daily_downloads: {type: "bigint"},
        weekly_downloads: {type: "bigint"},
        monthly_downloads: {type: "bigint"},
        all_time_downloads: {type: "bigint"}
      }
    }'
  };

---------------------------- AWESOME PACKAGE TABLE -----------------------------

CREATE TABLE IF NOT EXISTS awesome_packages (
  author text,
  repo text,
  PRIMARY KEY (author, repo)
);

------------------------- PACKAGE ARCHIVE RECORD TABLE -------------------------

CREATE TABLE IF NOT EXISTS package_archive_records (
  author text,
  repo text,
  sha text,
  PRIMARY KEY (author, repo, sha)
);

----------------------------- GITHUB API KEY TABLE -----------------------------

CREATE TABLE IF NOT EXISTS github_api_key (
  for_indexer boolean,
  key text PRIMARY KEY,
);

--------------------------- ALL-TIME DOWNLOADS TABLE ---------------------------

CREATE TABLE IF NOT EXISTS all_time_downloads (
  author text,
  repo text,
  sha text,
  total counter,
  PRIMARY KEY ((author, repo), sha)
) WITH CLUSTERING ORDER BY (sha ASC);

---------------------------- DAILY DOWNLOADS TABLE -----------------------------

CREATE TABLE IF NOT EXISTS daily_downloads (
  day timestamp,
  author text,
  repo text,
  total counter,
  PRIMARY KEY ((author, repo), day)
) WITH CLUSTERING ORDER BY (day DESC);
