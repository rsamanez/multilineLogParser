# Multiline Log Parser
 Basic parser, it needs to be improved with regular expressions    
 Output JSON format
 
 It works like a tail command   
    
 To be sure that the application is always running you can create a cron to relaunch the application, and always will be running only one instance of it
 
## Single
single provides a mechanism to ensure, that only one instance of a program is running.   
The package currently supports linux, solaris and windows.
## Requiere Packages to compile

```
go get github.com/hpcloud/tail/...
go get github.com/marcsauter/single
go build main.go
```

### Run Command
```
./main -f -n 1 my-file-to-tail.log
./main -f -n 1 my-file-to-tail.log >> /path/to/output_json.log
```

### Setup to use with New Relic infrastructure agent to send Logs
Requisite: Install the last New Relic infrastructure agent, version 1.11.4 or higher   
documentation: https://docs.newrelic.com/docs/logs/enable-log-management-new-relic/enable-log-monitoring-new-relic/forward-your-logs-using-infrastructure-agent

### Linux configuration sample
Navigate to the logging forwarder configuration folder:   
/etc/newrelic-infra/logging.d/
   
Create a configuration file (for example, logs.yml) with this content:   
```
# Remember to only use spaces for indentation
logs:
  - name: "test_log"
    file: /var/log/parsed_logs.log
```
Run the parser application
```
./main -f -n 1 /path/to/your_multiline_log_file.log >> /var/log/parsed_logs.log
```
TESTING   
To test it you can send line by line the content of the sample_multiline_log_to_parse.log file   
   
To do it you can use the next PHP sample code   
filename: readTest.php
```
<?php
$handle = @fopen("sample_multiline_log_to_parse.log", "r");
if ($handle) {
    while (($buffer = fgets($handle, 4096)) !== false) {
        usleep(10000);
        $f = @fopen("your_multiline_log_file.log","a+");
        fputs($f,$buffer);
        fclose($f);
    }
    if (!feof($handle)) {
        echo "Error: unexpected fgets() fail\n";
    }
    fclose($handle);
}
?>
```
```
./main -f -n 1 your_multiline_log_file.log >> /var/log/parsed_logs.log
```
Check your results in New Relic Insights
```
SELECT * FROM Log where hostname='YOUR-HOST-NAME'
```
