# Diplomat

For usage and documentation, see [the wiki](https://github.com/schoonology/diplomat/wiki).

## Development

### Get the Source Code

This project uses a version of Go supporting [modules](https://blog.golang.org/modules2019), which means you can clone the repository anywhere you want. It also handles dependencies without need of secondary tooling such as Glide.

```sh
git clone git@github.com:schoonology/diplomat.git
cd diplomat
```

### Setup

Install [`go`](https://golang.org/doc/install):

```sh
# Mac OSX with Homebrew
brew install go
```

Install [`bats-core`](https://github.com/bats-core/bats-core):

```sh
# Mac OSX with Homebrew
brew install bats-core
```

In order to develop on this repository, you'll also need to install [Docker](https://docs.docker.com/install/).

By default, we use a Docker image version of [`httpbin`](https://httpbin.org/) to make requests for the `bats` tests. If you prefer not to use Docker, change the value of `TEST_HOST` in `test/helpers/helpers.bash` to point to `https://httpbin.org`.

### Watch Code and Run Tests

The core development script in `Makefile` is `watch`:

```sh
make watch
```

This script uses `rg --files | entr -rc` to watch the code, and run several other `make` scripts on update. Those scripts can also be run on a one-off basis, if desired (see `Makefile`). The `watch` script provides a general-purpose, fast feedback cycle for development, which encourages making small changes and keeping tests green.

### Using Mocks

The unit tests in this repository use [`mockery`](http://docs.mockery.io/en/latest/) and [`testify`](https://github.com/stretchr/testify) to generate mocks for each `type` in the project. These mock files live in `mocks`, are checked into source control, and should _not_ be manually updated.

If you add or modify any `type`s, you will need to regenerate the mocks.

```sh
make generate
```

## Code of Conduct

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/) for all community interactions, including (but not limited to) one-on-one communications, public posts/comments, code reviews, pull requests, and GitHub issues. If violations occur, the maintainers and will take any action they deem appropriate for the infraction, up to and including blocking a user from repository access.
