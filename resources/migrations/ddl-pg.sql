CREATE TABLE trino_clusters
(
    name    varchar(128) primary key,
    url     varchar(256),
    tags    json    default '{}',
    enabled boolean default true
);


CREATE TABLE trino_queries
(
    id                                 varchar primary key,
    cluster                            varchar(128) primary key,

    query                              text,

    user                               varchar,
    principal                          varchar,
    catalog                            varchar,
    schema                             varchar,
    client_address                     varchar,

    resource_group                     varchar,
    submission_time                    timestamp,
    completion_time                    timestamp,

    elapsed_time                       interval,
    queued_time                        interval,
    analysis_time                      interval,
    planning_time                      interval,
    execution_time                     interval,

    resources_cpu_time                 interval,
    resources_scheduled_time           interval,
    resources_input_rows               bigint,
    resources_input_data               bigint,
    resources_physical_input_rows      bigint,
    resources_physical_input_data      bigint,
    resources_physical_input_read_time bigint,
    resources_internal_network_rows    bigint,
    resources_internal_network_data    bigint,
    resources_peak_user_memory         bigint,
    resources_peak_total_memory        bigint,
    resources_cumulative_user_memory   bigint,
    resources_output_rows              bigint,
    resources_output_data              bigint,
    resources_written_rows             bigint,
    resources_logical_written_data     bigint,
    resources_physical_written_data    bigint

)
