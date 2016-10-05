------------------------------- DAILY DOWNLOADS --------------------------------

DROP TABLE IF EXISTS daily_downloads;

--------------------------- ALL TIME INSTALL TOTALS ----------------------------

CREATE TABLE IF NOT EXISTS all_time_install_totals (
  package_author text,
  package_repo text,
  package_description text,
  total bigint,
  PRIMARY KEY (package_author, package_repo)
);

CREATE CUSTOM INDEX IF NOT EXISTS all_time_install_totals_index
  ON all_time_install_totals ()
  USING 'com.stratio.cassandra.lucene.Index'
  WITH OPTIONS = {
  'schema': '{
    fields: {
      total: { type: "bigint", sorted: true }
    }
  }'
};

---------------------------- RANGED INSTALL TOTALS -----------------------------

CREATE TABLE IF NOT EXISTS ranged_install_totals (
  date timestamp,
  package_author text,
  package_repo text,
  package_description text,
  daily_total bigint,
  monthly_total bigint,
  weekly_total bigint,
  PRIMARY KEY (date, package_author, package_repo)
) WITH CLUSTERING ORDER BY (package_author ASC, package_repo ASC);

CREATE CUSTOM INDEX IF NOT EXISTS ranged_install_totals_index
  ON ranged_install_totals ()
  USING 'com.stratio.cassandra.lucene.Index'
  WITH OPTIONS = {
  'schema': '{
    fields: {
      daily_total:    { type: "bigint", sorted: true },
      weekly_total:   { type: "bigint", sorted: true },
      monthly_total:  { type: "bigint", sorted: true }
    }
  }'
};

----------------------- RANGED VERSIONED INSTALL TOTALS ------------------------

CREATE TABLE IF NOT EXISTS ranged_versioned_install_totals (
  date timestamp,
  package_author text,
  package_repo text,
  package_version text,
  daily_total bigint,
  monthly_total bigint,
  weekly_total bigint,
  PRIMARY KEY (date, package_author, package_repo, package_version)
) WITH CLUSTERING ORDER BY (package_author ASC, package_repo ASC, package_version ASC);

CREATE CUSTOM INDEX IF NOT EXISTS ranged_versioned_install_totals_index
  ON ranged_versioned_install_totals ()
  USING 'com.stratio.cassandra.lucene.Index'
  WITH OPTIONS = {
  'schema': '{
    fields: {
      daily_total:    { type: "bigint", sorted: true },
      weekly_total:   { type: "bigint", sorted: true },
      monthly_total:  { type: "bigint", sorted: true }
    }
  }'
};
