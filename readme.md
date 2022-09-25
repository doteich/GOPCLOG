# gOPCLOG

`GOPCLOG` is a GO based OPC-UA Logger, which tracks tag value changes based on OPC-UA subscriptions and forwards them through a various channels (check roadmap) 

### Quickstart ⚙️
- Install and Download the Binary from Github
- Mount the Config-JSON to the instances config folder `./config`
- Start the Binary file

### Setting up the config file 🗒️
`GOPCLOG` works nicely with the [OPC-UA-Browser](https://github.com/doteich/OPC-UA-Browser), as it generates the config file through the UI. If you choose to set up the config by yourself you can use the sample down below
```json
{
    "opcConfig": {
        "url": "opc.tcp://IP_or_URL",
        "securityPolicy": "None",
        "securityMode": "None",
        "authType": "User & Password, Anonymous or Certificate are supported",
        "username": "*",
        "password": "*",
        "node": "ns=3;s=NODE"
    },
    "selectedTags": [{
        "nodeId": "ns=3;s=XYZ",
        "name": "TESTTAG1"
    }],
    "methodConfig": {
        "subInterval": 10,
        "name": "TestLogger",
        "description": null
    }
}
```
### Roadmap 🚀
- [ ] Docker/Kubernetes Integration
- [ ] Cert generation (:small_blue_diamond: works but needs further testing)
- [ ] Support for Websockets
- [ ] Support for gRPC
- [ ] Logging to a metrics endpoint `/metrics`
- [ ] Enhanced error handling and logging
- [ ] Support for different config types like YAML, TOML, ENV


