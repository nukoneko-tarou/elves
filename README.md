# elves
Generate the initial project directory from the json file generated by tree

# Installing
command and is ready to use!
```shell
go install github.com/nukoneko-tarou/elves@latest
```

# Usage
Specify a json file after the create command.
```shell
cd <target directory>
elves create <json file path>
```
e.g.
```shell
elves create ../sample.json
.
├── api
├── assets
├── build
│   ├── ci
│   └── package
├── cmd
│   └── _your_app_
├── configs
├── deployments
├── docs
├── examples
├── githooks
├── init
├── internal
│   ├── app
│   │   └── _your_app_
│   └── pkg
│       └── _your_private_lib_
├── pkg
│   └── _your_public_lib_
├── scripts
├── test
├── third_party
├── tools
├── vendor
├── web
│   ├── app
│   ├── static
│   └── template
└── website

30 directories, 0 files
```

# Flags
### --sub, -s
Create the specified directory at the current location and create directories under it.  
e.g.
```shell
elves create ./sample.json --sub new-project
```
### --permission, -p
Specify the PERMISSION of the file to be generated.  
Default is 755.  
e.g.
```shell
elves create ./sample.json --permission 777
```
