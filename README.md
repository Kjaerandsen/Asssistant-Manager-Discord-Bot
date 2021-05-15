# Personal Assistant discord bot

A discord bot you can place in your personal discord server to help you with tasks like:
Making dinner plans
Reminding you of tasks
Handling bills
Retrieving weather forecast from a saved location or a specified location
All information to the service will be sent through discord in your discord-server to the bot, and all information will be sent to your discord-server from the bot. For backend storage google firestore is used.

## Endpoints / discord commands and functionality

Descriptive instructions on how to use these functionalities are found using the: "**@Bot help**" command in discord after inviting it to your server
 
Weather:
- Get weather with a specified location
- Get weather with a default / saved location

Reminders:

- Add reminder
- Get reminder
- Remove reminder (name)

Fridge and cooking:

- Add to fridge (ingredient)
- Check fridge
- Remove from fridge (ingredient)
- Get recipe based on ingredients in the fridge
- Get recipe based on specified ingredients

News
- Get news with or without parameters

# Adding the bot to your own server

The requirements for you to add the bot to the server is that you have Manage Server privileges on the server.
Follow this [link](https://discord.com/oauth2/authorize?client_id=834015714200649758&permissions=8&scope=bot), you might need to log in to your discord account before the interface shows up. Once logged in, it will ask to connect to your server, scroll through the list of servers under add server and pick the server you want to add the bot in.
To add the bot to several servers, simply repeat this process for every server.

# Running it on your own computer

## Setup
Download the git repo.
Add the bot to your server as described above.
Set up the database as shown under "Database setup" further down.
Compile the program using golang.
Alternatively see the docker instructions below.
Run the program from within the base folder of the project..

## Docker setup:
Download the git repo.
Run "docker build -t discobot ." from within the project folder
Run the bot by using "docker run discobot"

### Database setup:
This project uses google's firebase firestore as the database. First set up a firebase account and project. Go to project settings, under the "Service accounts" tab click "Generate new private key". Accept the prompt and download the json key file. This file will be used by the application to connect to the database. Rename this file to "service-account.json" and move it to the DB folder of the git repo.

### Dependencies:
Discordgo https://github.com/bwmarrin/discordgo

## Original Project Plan
The original plan was to create the structure for a bot in which we could easily add "services'' through implementing many different external APIs. After implementing this structure we wanted to test it by adding four services, Bills, Meal Planner, Weather forecast and News fetcher. We wanted this structure to have a high level of scalability, only needing to do the bare minimum to implement new services.
We were successful in creating an environment in which it was easy to implement external APIs, reformat their responses into discord messages and post them to the user. However, it was not perfect.
Ideally we would have wanted to refactor the request handler and the message handler, move this message handler onto its own file such that it would be much more organized. There is also room for further refactoring of functions used within the services. With good refactoring, we still believe that all the services can be a part of the same package. But it has become apparent that without further refactoring, there is too much clutter, and each service should be their own package.
## Hard Aspects and Unimplemented Features
It was hard to find reliable API’s that did not have strict limits. One example is the meal API, spoonacular. On the free plan you are only allowed 100 API calls per day, where each recipe received also was worth 0.01 points (1 point = one call). The ideal would have been to both receive the recipe, and instructions. 

But to do this you would have to use two of their endpoints. The first endpoint has a cost of 1 call + 0.01 calls * amount of recipes. A call for 10 recipes would then cost 1.1 points, allowing for 90 calls per day. If we were to implement it with instructions it would cost 
1 call + 0.5 calls * amount of recipes  which would be 6 points. In total that would be 7.1 points per message, resulting in 14 possible calls per day. 

We planned to overcome this by caching many recipes over time, and only allowing the users to get random recipes. If we wanted to then introduce the customization we could introduce our own algorithms to find compatible recipes, but that was beyond the scope of the project.

Another solution to this could have been to make the customizable recipes a part of a paid membership. If we had a source of income from the bot, we could then reinvest some of that income in more API’s and easily expand the amount of services the bot provides.

Another time consuming part of the process was learning how to use the discord bot library in go. Some of our members already had experience in using the discord libraries for other languages, but the go library for discord did not provide the same amount of functionalities. 
The documentation for discordgo was ok, however a lot of time was spent trying to figure out the different functionalities as there was a lack of examples and well written descriptions.

Another hard aspect was time management. There were some easy fixes that we could have made, if we had the time. During the project all the project members had a lot of work in their respective courses.

This made us prioritize poorly, choosing to work in between assignments and projects deliveries, and resulted in worse code quality and structure. With better time management and organizing we could have organized work into small manageable chunks that we could do during the periods of the other course assignments. Even though these small chunks would most likely be done in a rush, reflecting over the project for a longer time would have increased the chances of the group coming up with new ideas and improvements.


## What have we learned

The group members had different strengths and weaknesses. Some were proficient in using firestore, others had great skill in coding and some were good at docker, etc.

We learned a lot from each other. In areas where we had low knowledge, other project members would teach us what we were missing. And our skill sets became more well rounded as we drew from each other's knowledge. A list of some are:

- Discord API
- Dockerfiles & deployment
- Firestore, setting up and using it
- The value of good time management
- Better code modularity and methods of structuring the program
- Use of branches in Git





## Total work hours

143
