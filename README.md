# SCP Interactive File Copy Utility

This is a command-line utility written in Go that enables you to interactively choose a file from your local directory and copy it to a remote host using the SCP protocol. It utilizes the `github.com/AlecAivazis/survey/v2` library to provide a user-friendly interface for selecting files and configuring the copy process.

## Features

- Interactive file and directory selection.
- Customizable configuration options for the copy process.
- Copy files to a remote host using SCP.
- Support for specifying remote port, username, host, and destination path.
- Simplified and streamlined user experience.

## Prerequisites

- Go (Golang) installed on your system.
- Basic familiarity with the command-line interface.

## How to Use

1. Clone or download the repository to your local machine.

2. Open a terminal and navigate to the directory containing the cloned repository.

3. Run the following command to build and execute the program:
   
   ```bash
   go run main.go
   ```
    The program will guide you through the process of selecting a file, configuring the remote host details, and initiating the file copy.

4. Follow the on-screen prompts to navigate through directories, choose a file, and configure the copy settings.

5. Once the configuration is complete, the program will use SCP to copy the selected file to the specified remote host.

## License
This project is licensed under the [MIT License](LICENSE).