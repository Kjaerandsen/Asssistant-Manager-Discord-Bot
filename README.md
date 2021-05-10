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

### Setup

- Download the git repo.
- Compile the program using go.
- Create a discord bot with message and webhook rights and add it to your discord server.
- Set the system env BOT_TOKEN to the discord bot's token.
- Set up the database as shown under "Database setup" further down.
- Run the program.

#### Database setup: 
This project uses google's firebase firestore as the database. First set up a firebase account and project.
Go to project settings, under the "Service accounts" tab click "Generate new private key".
Accept the prompt and download the json key file.  This file will be used by the application to connect to the database.
Rename this file to "service-account.json" and move it to the DB folder of the git repo.

### Dependencies:

Discordgo https://github.com/bwmarrin/discordgo
