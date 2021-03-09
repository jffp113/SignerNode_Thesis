# Signer Node

## How to build?

Go to the SignerNode folder and execute the following command.
```bash
docker build -t signernode -f Docker/Dockerfile .
```


## How can I run?

After building the project and the dependencies:
- Crypto Provider
- SmartContract 
- ...

you can run by executing the following command:


```bash
docker-compose  -f Docker/docker-compose.yaml up
```