<a id="readme-top"></a>


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <img src="images/logo.png" alt="Logo" width="80" height="80">

  <h3 align="center">Anonymous Telegram Chat Bot</h3>

  <p align="center">
    An anonymous chat for Telegram.
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#features">Features</a></li>
  </ol>
</details>



## About The Project

This is a fun little Telegram bot that lets people chat anonymously. Messages are sent without showing who wrote them, so you can message freely without anyone knowing your identity. It’s just a simple project made for casual and open conversations—nothing too serious, just a way to chat without names attached!



## Usage 

To run the bot, clone this repository and start it using either Go or Docker.  

---

### Run using Docker  

1. **Build the Docker image:**  
   ```sh
   docker build -t anonymous-chat-bot .
   ```  

2. **Run the container:**  
   ```sh
   docker run -d --name chat-bot \
     -e TOKEN=<your-bot-token> \
     -e DB_URI=<your-database-uri> \
     -e ADMIN_ID=<your-admin-id> \
     anonymous-chat-bot
   ```  

   - The `-d` flag runs the container in detached mode.  
   - Replace `<your-bot-token>`, `<your-database-uri>`, and `<your-admin-id>` with actual values.  

---

### Run using Go  

```sh
go run main.go $TOKEN $DB_URI $ADMIN_ID
```  

#### Arguments / Environment Variables:  
- **TOKEN** – Telegram Bot Token (Get it from [@BotFather](https://t.me/BotFather))  
- **DB_URI** – PostgreSQL connection string  
- **ADMIN_ID** – The bot admin’s Telegram ID (Find yours via [@userinfobot](https://t.me/userinfobot))  



## Features     

- **Fully Anonymous Chat** – Messages are copied and broadcasted to all users without revealing the sender’s identity.  
- **Root Users** – Only root users can see the sender’s name (but not their username).  
- **Admin Control** – The admin (set via CLI argument) can manage root users with:  
  - `/give_root [username]` – Grants root access to a user.  
  - `/remove_root [username]` – Revokes root access from a user.  
- **User Count Command** – All users can check how many people are in the bot’s database with `/count`.  
