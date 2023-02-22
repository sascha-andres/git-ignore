# git ignore

Manage your local and global ignore file

## Operations

### list

Print out the git ignore file

### add

Add an entry from git ignore file

### remove

Remove an entry from git ignore file

## Options

- -unique: save unique lines only
- -global: edit global git ignore file

Options may be provided by environment variables.

Rule: prefix capitalized flag name with GIT_IGNORE_. Example: `GIT_IGNORE_UNIQUE`. Bollean values must be set by providing `true` as the value.
