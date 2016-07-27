# Telegraf Input Plugin: Fleet

The plugin will gather names of running units from [fleet](https://github.com/coreos/fleet) and the sum total of each running unit. It uses the fleet v1 API to gather data. 

### Configuration:

```toml
# Description
[[inputs.fleet]]
## Works with Fleet HTTP API
## Multiple Hosts from which to read Fleet stats:
	hosts = ["http://localhost:49153/fleet/v1/state"]
```

### Measurements & Fields:

The fields are dynamically generated from the output of the fleet API. Using the ```name``` value.. The values of those fields are the number of containers  with the ```systemdSubState``` value of "running".   
<insert example json output here of both running and not to show what is included and not included>

The unit names will have their instanced id and the @ symbol stripped off.  
For example if you had a unit named ```nginx-1.10.1@35``` the field name would be ```nginx-1.10.1```. 

- fleet
    - ```<dynamic unit name>``` (int)

### Tags:

- All measurements have the following tags:
    - container_name (unit or container name found running in fleet cluster)
    - server (name of the fleet node this unit or container resides on)

### Example Output:

```
$ ./telegraf -config telegraf.conf -input-filter example -test
measurement1,tag1=foo,tag2=bar field1=1i,field2=2.1 1453831884664956455
measurement2,tag1=foo,tag2=bar,tag3=baz field3=1i 1453831884664956455
```