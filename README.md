# Xray Microsoft Teams Integration
This integration will publish Xray violations to a channel in Microsoft Teams.

The integration utilizes Xray webhooks to publish violations to the integration server built by the code in this project.

The integration server then routes the messages over to Microsoft Teams.

# Setup

## Microsoft Teams Channel Setup

You will need to be a member of the channel you wish to send messages to in Microsoft Teams.

You will need to setup a new Connector by click on the icon "..." to load more options next to the channel name.

In the more options menu select Connectors for the Connectors popup to appear.

In the search textbox type "Incoming Webhooks" and click Add.

Click Configure and supply the follow options:

Name: Xray Webhook

Image: [Download](images/xray.png)

Save the image and upload the file when creating the new incoming webhook in Microsoft Teams for the picture.

Scroll down after you submit to see the URL of the incoming webhook. Save this value.

You need to set this as an environment variable MICROSOFT_TEAM_WEBHOOK:

``` 
export MICROSOFT_TEAM_WEBHOOK=<url>
```

## Integration Server

Bbuild the integration server by running go build inside the project folder:

``` 
go build
```

This will create a new binary object 'xray_msteam' inside the same folder.

Ensure you have set the necessary environment variable MICROSOFT_TEAM_WEBHOOK for the integration to know where to post the Xray violations message.

Run the output binary from go to start the integration server. 

Take note of the IP/hostname of the machine you ran this one.

You will need it to build the below URL for the Xray Webhook configuration.

``` 
http://{host}:{port}/api/send
```

Example:
```
http://localhost:8080/api/send
```

#### Assigning Different Port Number

By default the integration server runs on port 8080. If you wish to change this port you can update the value in routes.go on line 26.

Example using port 9090:

``` 
return &http.Server{Addr: ":9090", Handler: &serverHandler{}}
```

#### Connectivity Test

The integration server has a Ping endpoint you can hit to test connectivity from your Xray node.

Using the Host or IP address of the integration server you can test it via Ping by invoking:

``` 
curl http://{host}:{port}/api/ping
```

Example:

``` 
curl -v http://localhost:8080/api/ping
```

Result:

``` 
 HTTP/1.1 200 OK
 Date: Thu, 30 Jul 2020 02:54:06 GMT
 Content-Length: 0
```

## Xray Webhook

Setup a webhook in Xray by opening the JFrog Unified Platform in a web browser.

You can then goto Admin -> Xray -> Configure Webhooks to create a new webhook.

You will need to supply the URL of the integration server that you deployed in the above step.

Once you have created the webhook inside of Xray we will need to add it as a new rule into a Policy.

Assuming this Policy is associated to a Watch that encounters an artifact in a repo that has violations it will then generate an outbound webhook event to our integration server.

The integration server built from this project will post the messages to your channel based upon your configurations to Microsoft Teams.

## Demo

Follow the guide above to setup the webhook in your Microsoft Teams channel you wish to deliver the violation messages to.

Once you have exported the URL to the correct environment variable you will need to run the server built from the go code in this integration project.

You will need to use the url of the integration server to supply to Xray for the webhook.

To run this as a demo please create a new Orbitera trial of JFrog Xray available [here](https://jfrog.orbitera.com/c2m/trials/signup?testDrive=1500&goto=%2Fc2m%2Ftrial%2F1500)

Once your environment has been created you will receive an email with the URL to JFrog Unified Platform & admin account password.

You can then use this to login to access Xray and setup the webhook as described above.

## Tools
* [JFrog Xray](https://jfrog.com/xray/) - JFrog Xray Security Scanner
* [Microsoft Teams](https://www.microsoft.com/en-us/microsoft-365/microsoft-teams/group-chat-software) - Microsoft Teams

# Contributing
Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests to us.

# Versioning
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).

# Contact
* Github
