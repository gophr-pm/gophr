----------------------------------- PACKAGES -----------------------------------

ALTER TABLE gophr.packages ADD versions list<text>;

------------------------------- PACKAGE VERSIONS -------------------------------

DROP TABLE IF EXISTS gophr.package_versions;
