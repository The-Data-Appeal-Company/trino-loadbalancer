CREATE TABLE presto_clusters
(
    name         varchar(128) primary key,
    url          varchar(256),
    tags         json    default '{}',
    enabled      boolean default true,
    distribution varchar(64)
);


