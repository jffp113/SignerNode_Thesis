# Signer Node

## How to build?

You don't need to build to program to run it.
This process is only done by project developers, if you only want to run, go to the next section.
Go to the SignerNode folder and execute the following command.

```bash
make build-docker
```

If you have auth to publish in the repo run: (this command will build and publish the docker image)

```bash
make push-docker
```

## How can I run?

You can run the project with 5 nodes and TBLS256 n = 5 , t = 3 by running the following command

For permissioned:

```bash
make run-permissioned
```

For permissionless:

```bash
make run-permissionless
```

To stop running the previous configuration, you only need to execute:

For permissioned:
```bash
make stop-permissioned
```

For permissionless:
```bash
make stop-permissionless
```
