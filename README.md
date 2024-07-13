# SCA
Spy Cat Agency (SCA)

# installation

1. git clone git@github.com:blockseeker999th/SpyCat.git (SSH)

   git clone https://github.com/blockseeker999th/SpyCat.git (HTTPS)

2. Set up applications in docker
You can use make command from Makefile
``` make build ``` to start application and build all the necessary stuff
``` make restart ``` to restart the applications
``` make rebuild ``` to rebuilt the application in if you made some changes
or do it through standart docker compose commands

Note: if you use newer version of docker compose(v2), you need to remove dash(-) in commands like that ``` docker compose up ```

``` docker-compose up --build ```
sometimes the application can't wait to load DB-container, in that case do Ctrl+C and then
``` docker-compose up ``` again, we don't need to rebuild the project in that case

3. shutdown application -> Ctrl+C
``` docker-compose down ``` to stop docker-compose containers,networks and volumes