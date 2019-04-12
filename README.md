# BufferFramework

RPC Bi-Directional Streaming Client and Server which returns if number is max at that moment

# Server
Bi Directional server which return number if it is max, and data is duly signed 

## Running Server
```
git clone https://github.com/niktrix/bufferframework
cd bufferframework/server
go build
./server.go
```

# Client
Bi Directional server which send streams of request, Request contains data, public key and signed data

```
type  Req  struct {
	Num int32   
	Key []byte   
	SignedData []byte  
}
```


## Running Client
```
git clone https://github.com/niktrix/bufferframework
cd bufferframework/client
go build
./client.go
```

Client generates random keypair, sign data and send to stream


## Code Walkthrough

### Server
Server starts listening on port mentioned in config.json
For all request server checks for public key and signed 
Decode key to Block which can be used by 

PEMBlock, _  := pem.Decode(req.Key)

and unmarshal to to key rsa.Public key


### Client
crypt package contains functions  for

1) SignData Requires data and rsa private key
2) MarshalPublicKey public key to bytes which can be wired
3)UnMarshalPublicKey publi to get back rsaPublic 
4) GetCerts generates keypairs

