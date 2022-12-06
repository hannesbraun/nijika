# Nijika

*Definitely not the command player*

This is a tool to run a bunch of commands in series or in parallel. You decide.

## Building

`go build nijika.go`

Granted, that was quite obvious...

## Running

Use the option `-id` to set an ID for the instance. Nijika will then read the file named `commands<id>.txt`
from the current working directory. Within that file, one line equals one command/task. Those will be
processed in order: commands at the top of the file are executed first. They will be executed in the current
working directory where Nijika itself also runs. At the start of the execution, the command will be logged
together with a timestamp into a file named `commandLog<id>.txt`. If a command does not complete
successfully, you'll also find an error message for this in the log file.

Use the option `-p` to allow more than one task to be executed at a time. With `-p 4` for example, four
tasks at most can run at the same time.

## License

This software is licensed under the terms of the [GNU General Public License, Version 3](LICENSE).
