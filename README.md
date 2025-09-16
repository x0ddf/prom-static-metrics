# prom-static-metrics
this is a simple Prometheus server which exposes list of static metrics configured through config file

## server config

server configured through the env vars:
 - `CONFIG_FILE` - path to the config file
 - `PORT` - server port

## config format

```yaml
metrics:
  - name: release_version # name of the metric
    value: 2025.10.192 # metric value, 
    description: prints release information
  - name: is_some_feature_flag_enabled
    value: true #  could be whatever -> converted to the 
    description: Prints state of the feature flag # metric description
```

As a result, you'll see on the /metrics endpoint of the server the following set of metrics:
```
# HELP is_some_feature_flag_enabled Prints state of the feature flag
# TYPE is_some_feature_flag_enabled gauge
is_some_feature_flag_enabled{is_some_feature_flag_enabled="true"} 1
# HELP release_version prints release information
# TYPE release_version gauge
release_version{release_version="2025.10.192"} 1```