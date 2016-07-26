-- Create the keyspace for the database if it doesn't exist.
CREATE KEYSPACE IF NOT EXISTS gophr
  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}
  AND durable_writes = true;

----------------------------------- PACKAGES -----------------------------------

CREATE TABLE gophr.packages (
  author text,
  repo text,
  description text,
  awesome_go boolean,
  exists boolean,
  godoc_url text,
  search_blob text,
  index_time timestamp,
  versions list<text>,
  PRIMARY KEY (author, repo)
) WITH CLUSTERING ORDER BY (
  repo ASC
);

CREATE CUSTOM INDEX packages_index
  ON gophr.packages ()
  USING 'com.stratio.cassandra.lucene.Index'
  WITH OPTIONS = {
    'schema': '{
      fields : {
        search_blob : { type : "text" }
      }
    }'
  };

---------------------------- RANGED INSTALL TOTALS -----------------------------

CREATE TABLE gophr.ranged_install_totals (
  date timestamp,
  package_author text,
  package_repo text,
  package_description text,
  daily_total bigint,
  monthly_total bigint,
  weekly_total bigint,
  PRIMARY KEY (date, package_author, package_repo)
) WITH CLUSTERING ORDER BY (package_author ASC, package_repo ASC);

CREATE CUSTOM INDEX ranged_install_totals_index
  ON gophr.ranged_install_totals ()
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
