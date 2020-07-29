# Xray Microsoft Teams Integration
This integration will publish Xray violations to Microsoft Teams.

The integration utilizes Xray webhooks to publish violations to server built by the Go code in this project.

The integration server then routes the messages over to Microsoft Teams via the Microsoft Graph API in Azure.

# Setup

## Microsoft Teams Setup

You will need to be an administrator in Azure to setup the integration between Xray and Microsoft Teams.

In terms of permissions you will need the ability to create a new service principal (app) and grant it:

Microsoft Graph API Permissions

``` 
ChannelMessage.Send
Group.ReadWrite.All
```

It is recommended you review the Azure access documentation [here](https://docs.microsoft.com/en-us/graph/auth-v2-user)

## Microsoft Team Id

This is the tenant id of your Microsoft Team account. 

You can find this tenant id inside of Microsoft Teams by right clicking on your team name and clicking on "Get Team Link"

Inside of this team link there is a tenant id. This tenant id is the id you want for your team.

You need to set this as an environment variable MICROSOFT_TEAM_ID such as:

``` 
export MICROSOFT_TEAM_ID=<tenant_id>
```

## Microsoft Team Channel Id

This is the tenant id of your Microsoft Team Channel to post your Xray violations to as messages. 

You can find this by using the Graph API [here](https://docs.microsoft.com/en-us/graph/api/channel-list?view=graph-rest-1.0&tabs=http#code-try-1)

You will need to login first to your Azure account.

You will then need to update the {id} to the MICROSOFT_TEAM_ID from above.

Grant consent to the necessary Graph API endpoints (Azure admin permissions required).

Trigger the query to get back a payload of all channels in your Microsoft Teams. Look up the channel you desire and copy the tenant id.

You need to set this as an environment variable MICROSOFT_TEAM_ID such as:

``` 
export MICROSOFT_TEAM_CHANNEL_ID=<tenant_id>
```


## Microsoft Access Token

The integration will need an access token with Microsoft Graph API permissions for:

ChannelMessage.Send, Group.ReadWrite.All

Microsoft Graph documentation [here](https://docs.microsoft.com/en-us/graph/api/channel-post-messages?view=graph-rest-1.0&tabs=http)

Both of these are required to send a message to a channel in Microsoft Teams.

To create an access token you will need to make a new service principal for this integration.

Information on how to create a new App / service principal in Azure is available [here](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal)

Once you have a new App registered you can then create a client secret inside of the app registration. 

Save this value as we will use it for the app_client_secret_value in the curl below.

To get a new access token against the Microsoft Graph API (resource id 00000003-0000-0000-c000-000000000000) run the curl below replacing the relevant fields with your data:

``` 
curl -X GET -H 'Content-Type: application/x-www-form-urlencoded' \
-d 'grant_type=client_credentials&client_id=<app_tenant_id>&client_secret=<app_client_secret_value>&resource=00000003-0000-0000-c000-000000000000' \
https://login.microsoftonline.com/<subscription_tenant_id>/oauth2/token
```

subscription_tenant_id -> tenant id of your azure subscription

app_tenant_id -> tenant id of the app we registered

app_client_secret_value -> value of the client secret we created in the app

The response payload will contain an access token which can then be used for auth back into Microsoft Teams by our integration.

We then need to set this as an environment variable MICROSOFT_ACCESS_TOKEN such as:

``` 
export MICROSOFT_ACCESS_TOKEN=<access_token>
```

## Integration Server

First build the integration server by running inside of the project folder:

``` 
go build
```

This will create a new binary object 'xray_msteam' inside the project folder.

Ensure you have set the necessary environment variables MICROSOFT_TEAM_ID, MICROSOFT_TEAM_CHANNEL_ID, MICROSOFT_ACCESS_TOKEN for the integration to know auth into the Microsoft Teams account and channel to post the Xray SIEM violations messages.

Run the binary to start the integration server.

#### Assigning Different Port Number

By default the integration server runs on port 8080. If you wish to change this port you can update the value in routes.go on line 26.

Example using port 9090:

``` 
return &http.Server{Addr: ":9090", Handler: &serverHandler{}}
```

## Xray Webhook

Setup a webhook in Xray by opening the JFrog Unified Platform in a web browser.

You can then goto Admin -> Xray -> Configure Webhooks to create a new webhook.

You will need to supply the URL of the integration server that you deployed in the above step.

Once you have created the webhook inside of Xray we will need to add it as a new rule into a Policy.

Assuming this Policy is associated to a Watch that encounters an artifact in a repo that has violations it will then generate an outbound webhook event to our integration server.

The integration server built from this project will post the messages to your channel based upon your configurations to Microsoft Teams.

## Demo

To run this as a demo you will need a Microsoft Teams account that is associated to a school or organization.

Personal accounts are not supported due to limitations on the Microsoft APIs that support Microsoft Teams application.

Please follow the above steps to determine the correct values for MICROSOFT_TEAM_ID, MICROSOFT_TEAM_CHANNEL_ID, MICROSOFT_ACCESS_TOKEN as described above.

Once you have exported these to the correct environment variable you will need to run the server built from the go code in this integration project.

This will then be the URL to supply to Xray for the webhook.

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
