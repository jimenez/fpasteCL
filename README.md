fpasteCL
========

[![Build Status](https://travis-ci.org/jimenez/fpasteCL.png?branch=master)](https://travis-ci.org/jimenez/fpasteCL)

A command line client for [fpaste.org](http://fpaste.org) written in Go

```
Usage: fpasteCL [-hP] [-e value] [-l value] [-p value] [-u value] [FILE...]
 -e, --expire=value  Seconds after which paste will be deleted from server
 -h, --help          Display this help
 -l, --lang=value    The development language used
 -p, --pass=value    Add a password
 -P, --private       Private paste flag
 -u, --user=value    An alphanumeric username of the paste author

```

### Example:

With stdin and a pipe,
```
$> echo "my logs to share" | fpasteCL -P -u devops
http://fpaste.org/#number#/#hash#
```

From multiple files,
```
$> fpasteCL -P -u devops file1 file2 file3
http://fpaste.org/#number#/#hash#
http://fpaste.org/#number#/#hash#
http://fpaste.org/#number#/#hash#
```