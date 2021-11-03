#!/usr/bin/expect -f

set username [lindex $argv 0]
set password [lindex $argv 1]

set timeout -1
spawn vmrest -C

expect "Username:"
send -- "$username\n"

expect "New password:"
send -- "$password\n"

expect "Retype new password:"
send -- "$password\n"

expect "Credential updated successfully"