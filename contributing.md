# Contributing

## Limit the number of external libraries

While we are not going to reinvent the wheel all the time, there must be a good reason to add an external library. This software values security very high so all external libraries have to be vendored and checked in to prevent an attack by manipulating external git repositories and slipping in our codebase without us noticing. The less dependencies, the lower the risk and the higher the likelihood of us noticing a suspicious change.
