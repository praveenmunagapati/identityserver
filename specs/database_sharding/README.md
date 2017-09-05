# Database sharding

This directory will contain the markdown documents for the database sharding so
they can easily be reffered to. The database backend is `Mongodb`. Next to the
architecture of the database itself, some backend changes will also be required.
Most notably, the collections which will be sharded will require an additional field
to be used as the shard key.

For a high level overview of the proposed changes, see
[high_level_design.md](high_level_design.md)
