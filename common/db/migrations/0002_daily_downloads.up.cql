------------------------------- DAILY DOWNLOADS --------------------------------

CREATE TABLE daily_downloads (
  day timestamp,
  author text,
  repo text,
  version text,
  total counter,
  PRIMARY KEY ((author, repo), day, version)
) WITH CLUSTERING ORDER BY (day DESC, version ASC);

--------------------------- ALL TIME INSTALL TOTALS ----------------------------

DROP INDEX IF EXISTS all_time_install_totals_index;
DROP TABLE IF EXISTS all_time_install_totals;

---------------------------- RANGED INSTALL TOTALS -----------------------------

DROP INDEX IF EXISTS ranged_install_totals_index;
DROP TABLE IF EXISTS ranged_install_totals;

----------------------- RANGED VERSIONED INSTALL TOTALS ------------------------

DROP INDEX IF EXISTS ranged_versioned_install_totals_index;
DROP TABLE IF EXISTS ranged_versioned_install_totals;
