module main

go 1.26.6

require types v0.0.0 // indirect
require server v0.0.0
require loadConfig v0.0.0

replace types => ./Types
replace server => ./Server
replace loadConfig => ./LoadConfig
