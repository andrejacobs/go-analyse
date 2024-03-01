# Usage

Use this file to keep track of how the CLI tools should be run. This will be used as the reference
point when it comes to writing the final end user documentation.

### Current commands I run - EXAMPLE

This is just to document what I run at the moment, so that I can write better documentation later.
Also a way to remember in a few months when I come back to this on how I even used this (before the docs)

-   Create a new database

        make build && clear && ./build/current/ajfs scan --db ~/temp/test.ajfs ~/temp

        or to create the database and calculate the hashes

        ajfs scan --hash --progress --db ~/temp/test.ajfs ~/temp

-   Continue to calculate the hashes (if previous hashing was stopped)

          ajfs update-pending --verbose --progress --db ~/temp/test.ajfs
