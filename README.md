# Core Flowbased Status CLI
Create and supported by [Amun Analytics](https://amunanalytics.eu/)  
This is a terminal tool to interact and display the [Core Flowbased status page](https://status.coreflowbased.eu/).
It is provided for Windows, Mac (both intel and silicon) and Linux.
You can display the data in a friendly way in your terminal for each businessday you want. It autoscales the widths and colors the cells depending on the status.

Please note that if you want to get the data in your own application you can use the API with curl by appending /json or /xml to the status page url. 
This tool is meant for a viewing purpose in your terminal.

## Usage
There are the following ways to use the tool:
```fbstatus version``` -> this prints the current version and git hash of the tool  
```fbstatus``` -> prompts the user for a business day and display the full status table  
```fbstatus <businessday>``` -> displays the full status page directly of the specified businessday.  
```fbstatus short``` -> prompts the user for a business day and shows a short table with trafic lights  
```fbstatus short <businessday>``` -> displays the short table directly for the specified businessday  

Any business day input should have the format ```YYYY-MM-DD```

## Screenshot
![Screenshot](img/screenshot.png)

## Telemetry
This tool gathers some telemetry data and sends it to my server, namely your OS architecture and how often you use the tool. 
If you whish to disable this set the environment variable ```AMUN_DISABLE_TELEMETRY="1"```