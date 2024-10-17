# go Budget Sharing App 

This guide walks you through how to download, compile, and run the Go Budget Sharing App on your local machine.

## 1. install Go

to run the app, you'll need to have Go installed.

- **download Go**: vsit the [Go offcial download page](https://golang.org/dl/) and download the latest stable version for your OS.
  
 After installation, open your terminal check for Go version to see if it's successfully installed

    ```bash
    go version
    ```

## 2. clone the repository

to get the source code for the app, clone the repository from GitHub:

- **clone the repo**: open you terminal and run:

    ```bash
    git clone <repo-url>
    ```

    replace `<repo-url>` with the URL of the repository (EX: `https://github.com/your-username/go-budget-sharing-app.git`).

- **navigate to the project directory**:

    ```bash
    cd go-budget-sharing-app
    ```

## 3. install dependencies

If your app has external dependencies, you'll need to download them. Run the following command in the terminal:

```bash
go mod tidy
```

## 4. run the app

open a terminal and type:

    ```bash
go run main.go

    ```


## 5. access the app locally

open a browser window and type:
```
    http://localhost:8080

```