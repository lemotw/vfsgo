# File Model
There are many infomation that we should note in file or memory. This document will describe the file model in this project.

The following graph is how the file model looks like.

// `tree` result future

1. file: `RootInode` (keep all user information)
2. dir: []`{username}_pool` (each user has a pool to keep all file information, you can think it as a home directory)
    1. file: `UserInode` (keep all file information)
    2. dir: []`{block_id}` (each file has a block to keep all block information)
        1. file: `BlockInode` (keep all **file hash map** and **current block id** and **previous block id**)
        2. file: []`{filehash}` (keep file header)
