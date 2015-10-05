# instainer
###Run any Docker container on the cloud instantly

[![Join the chat at https://gitter.im/instainer/instainer](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/instainer/instainer?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

#What is instainer?

Instainer is a Docker container hosting service which allows run instantly any Docker container on the cloud with Heroku-style Git deployment.

We believe in [KISS](https://en.wikipedia.org/wiki/KISS_principle). When we started migration to Docker in our company, we felt that something was still missing. Docker brought amazing capabilities to our DevOps team, but still there wasn't any service to click and run any Docker containers instantly. We are doers. We developed Instainer for engineers who want to run Docker containers on the cloud instantly. Your feedbacks & thoughts are really welcome. Thanks in advance!


#Quick Start
####Open [www.instainer.com](www.instainer.com), click any container to run.

![Click any container](http://beta.instainer.com/docs/instainer.png)

####After a few seconds you can access your container's information

![Container Details](http://beta.instainer.com/docs/container-details.png)


#Alternative methods to running container on Instainer


####From command line  (Mac & Linux & Windows supported)

You can find Instainer CLI Installation Guide from [here](https://github.com/instainer/instainer/wiki/Installation-Instainer-Client).

After installing instainer-cli and configuring it you can type on command line to run an nginx container.

    ./instainer run nginx 



For more information about instainer-cli commands you can find Instainer CLI Commands documentation from [here](https://github.com/instainer/instainer/wiki#instainer-cli-commands).

####REST API
Make GET request to REST API 

    curl -i -H "API-KEY: XXXXXXXXX" \
    -X POST http://www.instainer.com/backend/api/container/run?image=nginx

You can find Instainer REST API documentation from [here](https://github.com/instainer/instainer/wiki#instainer-rest-api-documentation)

####From your favourite CI tool

travis-ci, bamboo or jenkins -- We are working on amazing plugins, follow us for updates.


#How can I access to my container’s files?

instainer provides Heroku-style Git deployment for your containers. After running your container; instainer automatically creates Git repository for you and pushes your container’s data into this repository. You can easily clone and change your data using Git. For more information you can look Access Container Data topic from [here](https://github.com/instainer/instainer/wiki#accessing-container-data).



#Limits
The only limitation is that non-permanent containers keep alive for 15 minutes! If you want to create permanent containers please sign in with GitHub. Permanent containers will also terminate if you don't use instainer in 48 hours.

#How can I bash to my container?
If you want to access your container’s bash you have two different ways;

- Web browser; 
    To access your container's bash you should Sign in with GitHub to Instainer. After login Click My Containers from top menu then use *Open Terminal* button on My Containers page.

![Bash Container](http://beta.instainer.com/docs/bash-container.png)
![Terminal](http://beta.instainer.com/docs/terminal2.png)

####It can also run nyancat :)


![Terminal](http://beta.instainer.com/docs/terminal.png)

- Command line;

    ./instainer bash CONTAINER_ID

#How can I access the container’s logs?

You can use Show Logs link on My Containers panel; Instainer REST API or Instainer CLI

####After signing in with Github you can find My Containers menu top of the page

![My Containers](http://beta.instainer.com/docs/my-containers.png)
![Show logs](http://beta.instainer.com/docs/redis-my-containers.png)
![Show logs](http://beta.instainer.com/docs/redis.png)


#### REST API
[REST API Get Logs](https://github.com/instainer/instainer/wiki#container-logs)

#### Command Line
[Instainer CLI tool for getting logs](https://github.com/instainer/instainer/wiki#accessing-logs). 

#Can I use docker-compose.yml?

Yes! You can use docker-compose section on www.instainer.com to run your compose file or you can use [Instainer CLI](https://github.com/instainer/instainer/wiki#instainer-cli-commands) or [Instainer REST API](https://github.com/instainer/instainer/wiki#instainer-rest-api-documentation).  

#Documentation
You can check documentation from [https://github.com/instainer/instainer/wiki](https://github.com/instainer/instainer/wiki)

#”I have a question.”
Please feel free to use GITTER room your questions. 
[![Join the chat at https://gitter.im/instainer/instainer](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/instainer/instainer?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


#Contribution
Please feel free to create issue and open pull request. Also mention us on Twitter and we will tell you to follow which rabbit! [@instainer](http://twitter.com/instainer)
