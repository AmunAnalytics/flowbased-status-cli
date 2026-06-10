# Core Flowbased Status CLI
Create and supported by [Amun Analytics](https://amunanalytics.eu/)  
This is a terminal tool to interact and display the [Core Flowbased status page](https://status.coreflowbased.eu/).
It is provided for Windows, Mac (both intel and silicon) and Linux.
You can display the data in a friendly way in your terminal for each businessday you want. It autoscales the widths and colors the cells depending on the status.

Please note that if you want to get the data in your own application you can use the API with curl by appending /json or /xml to the status page url. 
This tool is meant for a viewing purpose in your terminal.

## Usage
There are two ways to use this tool. Either interactively or through direct command line arguments.  
### Interactive
Simply start the cli and you will get a prompt to enter a businessday. Enter it in the required format and press enter.
### Through Arguments
There are 2 commandline arguments available:  
```fbstatus version``` -> this prints the current version and git hash of the tool  
```fbstatus <businessday>``` -> displays the status directly of the chosen businessday. The format should be ```YYYY-MM-DD```

## Screenshot
![Screenshot](img/screenshot.png)

## Telemetry
This tool gathers some telemetry data and sends it to my server, namely your OS architecture and how often you use the tool. 
If you whish to disable this set the environment variable ```AMUN_DISABLE_TELEMETRY="1"```