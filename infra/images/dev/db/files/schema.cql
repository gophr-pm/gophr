
CREATE KEYSPACE IF NOT EXISTS gophr WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  AND durable_writes = true;

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
) WITH CLUSTERING ORDER BY (repo ASC)
    AND bloom_filter_fp_chance = 0.01
    AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'}
    AND comment = ''
    AND compaction = {'class': 'org.apache.cassandra.db.compaction.SizeTieredCompactionStrategy', 'max_threshold': '32', 'min_threshold': '4'}
    AND compression = {'chunk_length_in_kb': '64', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor'}
    AND crc_check_chance = 1.0
    AND dclocal_read_repair_chance = 0.1
    AND default_time_to_live = 0
    AND gc_grace_seconds = 864000
    AND max_index_interval = 2048
    AND memtable_flush_period_in_ms = 0
    AND min_index_interval = 128
    AND read_repair_chance = 0.0
    AND speculative_retry = '99PERCENTILE';
CREATE CUSTOM INDEX packages_index ON gophr.packages () USING 'com.stratio.cassandra.lucene.Index' WITH OPTIONS = {
  'schema': '{
    fields : {
      search_blob : { type : "text" }
    }
  }'
};

CREATE TABLE gophr.all_time_install_totals (
    package_author text,
    package_repo text,
    total bigint,
    PRIMARY KEY (package_author, package_repo)
) WITH bloom_filter_fp_chance = 0.01
    AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'}
    AND comment = ''
    AND compaction = {'class': 'org.apache.cassandra.db.compaction.SizeTieredCompactionStrategy', 'max_threshold': '32', 'min_threshold': '4'}
    AND compression = {'chunk_length_in_kb': '64', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor'}
    AND crc_check_chance = 1.0
    AND dclocal_read_repair_chance = 0.1
    AND default_time_to_live = 0
    AND gc_grace_seconds = 864000
    AND max_index_interval = 2048
    AND memtable_flush_period_in_ms = 0
    AND min_index_interval = 128
    AND read_repair_chance = 0.0
    AND speculative_retry = '99PERCENTILE';
CREATE CUSTOM INDEX all_time_install_totals_index ON gophr.all_time_install_totals () USING 'com.stratio.cassandra.lucene.Index' WITH OPTIONS = {
  'schema': '{
    fields: {
      total: { type: "bigint", sorted: true }
    }
  }'
};

CREATE TABLE gophr.ranged_install_totals (
    date timestamp,
    package_author text,
    package_repo text,
    daily_total bigint,
    monthly_total bigint,
    weekly_total bigint,
    PRIMARY KEY (date, package_author, package_repo)
) WITH CLUSTERING ORDER BY (package_author ASC, package_repo ASC)
    AND bloom_filter_fp_chance = 0.01
    AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'}
    AND comment = ''
    AND compaction = {'class': 'org.apache.cassandra.db.compaction.SizeTieredCompactionStrategy', 'max_threshold': '32', 'min_threshold': '4'}
    AND compression = {'chunk_length_in_kb': '64', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor'}
    AND crc_check_chance = 1.0
    AND dclocal_read_repair_chance = 0.1
    AND default_time_to_live = 0
    AND gc_grace_seconds = 864000
    AND max_index_interval = 2048
    AND memtable_flush_period_in_ms = 0
    AND min_index_interval = 128
    AND read_repair_chance = 0.0
    AND speculative_retry = '99PERCENTILE';
CREATE CUSTOM INDEX ranged_install_totals_index ON gophr.ranged_install_totals () USING 'com.stratio.cassandra.lucene.Index' WITH OPTIONS = {
  'schema': '{
    fields: {
      daily_total: { type: "bigint", sorted: true },
      weekly_total: { type: "bigint", sorted: true },
      monthly_total: { type: "bigint", sorted: true }
    }
  }'
};

CREATE TABLE gophr.ranged_versioned_install_totals (
    date timestamp,
    package_author text,
    package_repo text,
    package_version text,
    daily_total bigint,
    monthly_total bigint,
    weekly_total bigint,
    PRIMARY KEY (date, package_author, package_repo, package_version)
) WITH CLUSTERING ORDER BY (package_author ASC, package_repo ASC, package_version ASC)
    AND bloom_filter_fp_chance = 0.01
    AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'}
    AND comment = ''
    AND compaction = {'class': 'org.apache.cassandra.db.compaction.SizeTieredCompactionStrategy', 'max_threshold': '32', 'min_threshold': '4'}
    AND compression = {'chunk_length_in_kb': '64', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor'}
    AND crc_check_chance = 1.0
    AND dclocal_read_repair_chance = 0.1
    AND default_time_to_live = 0
    AND gc_grace_seconds = 864000
    AND max_index_interval = 2048
    AND memtable_flush_period_in_ms = 0
    AND min_index_interval = 128
    AND read_repair_chance = 0.0
    AND speculative_retry = '99PERCENTILE';
CREATE CUSTOM INDEX ranged_versioned_install_totals_index ON gophr.ranged_versioned_install_totals () USING 'com.stratio.cassandra.lucene.Index' WITH OPTIONS = {
  'schema': '{
    fields: {
      daily_total: { type: "bigint", sorted: true },
      weekly_total: { type: "bigint", sorted: true },
      monthly_total: { type: "bigint", sorted: true }
    }
  }'
};
