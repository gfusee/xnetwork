# DB merger tool

- These tools are able to merge level-DBs directories by creating a new one containing all the data.

## The tools
### generalDBMerger tool

- This tool is a general purpose tool able to copy the data from any number of level-DBs into a new one.
It can contain any type of keys and values of any length and the processing & override rules are the following:
1. the tool copies the first source DB at the OS level (using provided raw-data copy functionality)
2. it then opens, in order, the next DBs provided as source and iterates over all existing keys and values, 
storing them in the destination DB.

How to use:

```
cd cmd/generalDBMerger
go build
```

after the compilation of the binary, the merge can be made by calling the binary with the following parameters:

```
mkdir destdb
./generalDBMerger -dest=./destdb -sources=./src1/db,./src2/db,./src3/db
```

for full flags list, launch the binary with the following parameter

```
./generalDBMerger -h
```

### trieMerger tool

< to be implemented >

## Audience

This tool should be as generic as possible, and it shouldn't have any custom code related to node instances.
