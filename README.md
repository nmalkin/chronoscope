Chronoscope
===========
Chronoscope is a tool for timing the execution of commands applied to files.

Example use case: say you want to find out how long some queries take on your database of choice, YourSQL.

    # Put each query in its own file.
    mkdir queries && cd queries
    vim query1.sql query2.sql query3.sql

    # How to run each query 2 times using Chronoscope.
    /path/to/chronoscope -n=2 yoursql --input-file

    # What Chronoscope will run for you.
    yoursql --input-file query1.sql
    yoursql --input-file query1.sql
    yoursql --input-file query2.sql
    yoursql --input-file query2.sql
    yoursql --input-file query3.sql
    yoursql --input-file query3.sql


Command-line options
--------------------
The command you want to run, along with any flags, is the last argument to Chronoscope.

You can configure Chronoscope itself with the following flags.

    -threads    (default=1)     number of threads to use
                                If you use this option, set the GOMAXPROCS
                                environment variable to an appropriate value.

    -n                          how many times each thread will run the command

    -quiet      (default=false) suppress the output of the command

    -dir        (default=".")   directory whose files to pass to the program
