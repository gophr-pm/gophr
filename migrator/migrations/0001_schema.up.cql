
-------------------------------- PACKAGES TABLE --------------------------------

CREATE TABLE IF NOT EXISTS packages (
  repo text,
  stars int,
  author text,
  awesome boolean,
  description text,
  trend_score float,
  search_blob text,
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
  'refresh_seconds': '3600',
  'schema': '{
    default_analyzer: "english",
    fields: {
      trend_score: {type: "float"},
      search_blob: {type: "text"},
      daily_downloads: {type: "bigint"},
      weekly_downloads: {type: "bigint"},
      monthly_downloads: {type: "bigint"},
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
  sha text,
  total counter,
  PRIMARY KEY ((author, repo), day, sha)
) WITH CLUSTERING ORDER BY (day DESC, sha ASC);
