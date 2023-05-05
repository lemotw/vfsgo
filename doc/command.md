# Command

## User Management

### register
```
register [username]
```

#### Response:
Add [username] successfully.

- Error: The [username] has already existed.
- Error: The [username] contain invalid chars.

### use
```
use [username]
```

#### Response:
Switch to [username] successfully, your path is [upath].
- Error: The [username] doesn't exist.

### cd
```
cd [foldername]
```
#### Response:
Change directory to [foldername] successfully.
- Error: The [foldername] doesn't exist.
- Error: The [foldername] is not a directory.
___

## Folder Management

### create-folder

```
create-folder [foldername] [description]?
```

#### Response:

Create [foldername] successfully.

- Error: You have to choose a user first.
- Error: The [foldername] contain invalid chars.
___

### delete-folder

```
delete-folder [foldername]
```

### Response:
Delete [foldername] successfully.

- Error: You have to choose a user first.
- Error: The [foldername] doesn't exist.

___

### list-folders

```
list-folders [--sort-name|--sort-created] [asc|desc]
```

#### Response:

List all the folders at your current scope in following formats:
```
[foldername] [description] [created at] [username]
```

Each field should be separated by whitespace or tab characters.
The [created at] is a human-readable date/time format.
The order of printed folder information is determined by the
--sort-name or --sort-created combined with asc or desc flags.

The --sort-name flag means sorting by [foldername] .
If neither --sort-name nor --sort-created is provided, sort the
list by [foldername] in ascending order.

- Warning: The [username] doesn't have any folders.

- Error: The [username] doesn't exist.

Prompt the user the usage of the command if there is an invalid flag.(should output to STDERR)

___

### rename-folder

```
rename-folder  [foldername] [new-folder-name]
```

#### Response:

Rename [foldername] to [new-folder-name] successfully.

- Error: You have to choose a user first.
- Error: The [foldername] doesn't exist.

___

## File Management

### create-file

```
create-file [foldername] [filename] [description]?
```

#### Response:

Create [filename] in [foldername] successfully.

- Error: You have to choose a user first.
- Error: The [foldername] doesn't exist.
- Error: The [filename] contains invalid chars.

___

### delete-file

```
delete-file [foldername] [filename]
```

#### Response:

Delete [filename] in [username] / [foldername] successfully.
- Error: You have to choose a user first.
- Error: The [foldername] doesn't exist.
- Error: The [filename] doesn't exist

___

### list-files

``` 
list-files [foldername] [--sort-name|--sort-created] [asc|desc]
```

#### Response:

List files with the following fields:

```
[filename] [description] [created at] [foldername]
```

Each field should be separated by whitespace or tab characters.
The [created at] is a human-readable date/time format.

The order of printed file information is determined by the
--sort-name or --sort-created combined with asc or desc flags.

The --sort-name means sorting by [filename] .
If neither --sort-name nor --sort-created is provided, sort the list by [filename] in ascending order.

- Warning: The folder is empty.
- Error: You have to choose a user first.
s- Error: The [foldername] doesn't exist.

Prompt the user the usage of the command if there is an invalid flag.(should output to STDERR)
Input Validation and Restriction
