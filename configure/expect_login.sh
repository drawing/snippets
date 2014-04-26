#!/usr/bin/expect --

if { [llength $argv] < 3} {
	puts "usage: $argv0 host username password [command]"
	exit 1
}

set timeout 45
set host  [lindex $argv 0]
set username [lindex $argv 1]
set password [lindex $argv 2]
set command [lindex $argv 3]

trap {
	set rows [stty rows]
	set cols [stty columns]
	stty rows $rows columns $cols < $spawn_out(slave,name)
} WINCH

spawn zssh -l$username $host

expect {
	"password:" {
		send "$password\r"
		expect -re ".*(>|#|@|$).*"
		send "$command\r"
	}

	"(yes/no)?" {
		send "yes\r"
		expect "password:" {
			send "$password\r"
			expect -re ".*(>|#|@|$).*"
			send "$command\r"
		}
	}
	timeout {
		exit 1
	}
}

interact

