# cwnotifier

Notifier for cherwell's SLA

To build use the following command:

    https://stackoverflow.com/a/36728885/1252947
    - go build -ldflags -H=windowsgui

The .exe file must be run having the cherwell.png file in the same folder, otherwise the image will not be displayed.

The log file will be created in the same folder as the .exe

In order for the program to connect to the database, there should be a "config.yaml" file in the same folder as the .exe file. Here is a basic template of config.yaml:

    job:
        start: "08:00" # A partir de qual horário o programa irá checar o cherwell
        end: "17:59" # Até qual horário o programa irá checar o cherwell
        sleepMinutes: 5 # De quanto em quanto tempo em minutos o programa deve checar o cherwell
        
    database: # Dados para a conexão com o banco de dados
        server: ""
        port: 1433
        user: ""
        password: ""
        databaseName: "" # Nome do banco de dados do cherwell