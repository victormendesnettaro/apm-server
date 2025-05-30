user_name = "USER"

# APM bench

worker_instance_type = "c6i.2xlarge"

# Elastic Cloud

# The number of AZs the APM Server should span.
apm_server_zone_count = 1
# The Elasticsearch cluster node size.
elasticsearch_size = "64g"
# The number of AZs the Elasticsearch cluster should have.
elasticsearch_zone_count = 2
# APM server instance size
apm_server_size = "8g"

# Standalone

standalone_apm_server_instance_size = "c6i.xlarge"
standalone_apm_server_volume_type   = "io2"
standalone_apm_server_iops          = 6000
standalone_moxy_instance_size       = "c6i.2xlarge"
