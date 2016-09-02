---------------------------- PACKAGE ARCHIVE RECORD ----------------------------

CREATE TABLE IF NOT EXISTS package_archive_records (
  author text,
  repo text,
  ref text,
  PRIMARY KEY (author, repo, ref)
);
