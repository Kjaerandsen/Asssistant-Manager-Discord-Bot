# Personal Assistant discord bot

A discord bot you can place in your personal discord server to help you with tasks like:
* Making dinner plans
* Reminding you of tasks
* Handling and reminding you of bills
* Retrieving weather forecast from a saved location or a specified location

All information to the service will be sent through discord in your discord-server to the bot, and all information will be sent to your discord-server from the bot. For backend storage google firestore is used.

## Planned endpoints / discord commands and functionality

Weather:
* Get weather with a specified location
* Get weather with a default / saved location


Bills:
* Add bill (reocurring or not)
* Get bill
* Remove bill (name)

Reminders:
* Add reminder
* Get reminder
* Remove reminder (name)

Fridge and cooking:
* Add to fridge (ingredient)
* Check fridge
* Remove from fridge (ingredient)
* Get recipe based on ingredients in the fridge
* Get recipe based on specified ingredients

# Adding the bot to your own server

The requirements for you to add the bot to the server is that you have **Manage Server** privileges on the server. 

Follow this [link](https://discord.com/api/oauth2/authorize?client_id=834015714200649758&permissions=8&scope=bot), you might need to log in to your discord account before the interface shows up. Once logged in, it will ask to connect to your server, scroll through the list of servers under _add server_ and pick the server you want to add the bot in. 

To add the bot to several servers, simply repeat this process for every server.

# Running it on your own computer

### Setup

- Download the git repo.
- Add the bot to your server as described above.
- Set up the database as shown under "Database setup" further down.
- Compile the program using go and run it from the base folder of the project. Alternatively see the docker instructions below.
- Run the program.

### Docker setup:
- Download the git repo.
- Run "docker build -t discobot ." from the project folder
- Run the bot by using "docker run discobot"

#### Database setup: 
This project uses google's firebase firestore as the database. First set up a firebase account and project.
Go to project settings, under the "Service accounts" tab click "Generate new private key".
Accept the prompt and download the json key file.  This file will be used by the application to connect to the database.
Rename this file to "service-account.json" and move it to the DB folder of the git repo.

### Dependencies:

Discordgo https://github.com/bwmarrin/discordgo

## Original Project Plan

The original plan was to create the structure for a bot in which we could easily add "services" through implementing many different external APIs. After implementing this structure we wanted to test it by adding three services, Bills, Meal Planner and News fetcher. We wanted this structure to have a high level of scaleability, only needing to do the bare minimum to implement new services.

We were successful in creating an environment in which it was easy to implement external APIs, reformat their responses into discord messages and post them to the user. However, it was not perfect. 

Ideally we would have wanted to refactor the request handler and the message handler, move this message handler onto its own file such that it would be much more organized. There is also room for further refactoring of functions used within the services. With good refactoring, we still believe that all the services can be a part of the same package. But it has become apparent that without further refactoring, there is too much clutter, and each service should be their own package. 

