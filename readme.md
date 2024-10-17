# Go budget sharing app 

Go real time budget sharing app for [CSE 40842](https://www3.nd.edu/~pbui/teaching/cse.40842.fa24/project02.html). The app allow

the app is deployed through [Railway](https://railway.app/). you can see a live version [here](https://go-budget-sharing-app-production.up.railway.app/).

## 1. install Go

to run the app, you'll need to have Go installed.

- **download Go**: visit the [Go offcial download page](https://golang.org/dl/) and download the latest stable version for your OS.
  
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
```bash
http://localhost:8080
```