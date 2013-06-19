APNS binding for Go
====================

This is a basic wrapper around the Apple Push Notification Service. It is built on the client-server module that is inspired by [pynaps](https://github.com/samuraisam/pyapns).


Pushing notifications
------------------
The client needs to be initialized, configured, and provisioned before sending notifications to APNS.
```go  
c := &apns.Client{}
c.Configure(8080) // set the port
c.Provision("Foo", "path_to_cert", "sandbox") // set the individual appId, path to .pem certificate, and environment
```    
Once this setup is complete, the server will need to be running in another process. This allows for the use of multiple clients pushing to the server for the application while maintaining an open SSL connection to the APNS servers.
```go
apns.StartServer("sandbox", 8080)
```    
Now that the client and server are ready to talk, we can build a notification client side with the hexlified device token, an APNS payload, and notification identifier. After that is initialized, send it off!

```go
notification := &apns.Notification{
  Token:      "hexlified token",
  Payload:    &apns.Payload{Type: "alert", Message: "FLATBUSH"},
  Identifier: "0",
}
c.Notify("Foo", notification)
```
