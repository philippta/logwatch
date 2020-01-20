# Logwatch

logwatch watches a log file for changes and sends a notification to a specified email.

## Features

* Configure logwatch to search for a regular expression before sending emails.
* It keeps track of newly created files, so it handles log rotation gracefully.
* Does not parse the entire file upon start, so it only sends emails on new log file actions.
* Customize sender, subject, reciepent and search pattern

## Installation

Download the `bin/logwatch` binary from this repository to your server and mark it as executable with `chmod +x logwatch`.

## Usage

```
$ ./logwatch -help

Usage of ./logwatch:
  -file string
        Log file to watch
  -from string
        Email FROM (default "logwatch@localhost")
  -regex string
        RegEx pattern for matching lines
  -subject string
        Email subject (default "New log entry")
  -to string
        Email TO
```

## Example usage

Watch Apache log files for PHP fatal errors:
```
./logwatch \
    -file /var/log/apache2/error.log \
    -regex "PHP Fatal error" \
    -from logwatch@mycompany.com \
    -to bob@mycompany.com \
    -subject "A PHP fatal error occured"
