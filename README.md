# TodoList_Go

![Preview](https://placehold.co/500x300 "Preview of AppName")


## :information_source: About  

This is a back-end project type "Todolist" made with Go.


## :wrench: Tools
- Go 1.22.6
- HTML/CSS
- MySQL 8.4.2 

<!-- 
    SETUP
    Explain using command lines, the steps to follow to setup the project
    At the end show, the expected result with a image   

    Ex: 
    1. Download the whole project `Travel` on your system
    2. Open your terminal in `Travel`
    ```
    cd Travel
    ```
    3. In `Travel` directory, run:
    ```
    go run github.com/loickcherimont/Travel/main
    ```
    4. If there is no error. Go on your favorite browser and use this line in your URL address bar
    ```
    http://localhost:8080/travel
    ```
    5. Here you are! Welcome in the main page of the Web application

    ![Main page of the application](assets/images/readme_images/mainpage.png)
-->

## :inbox_tray: Setup

*It's coming ...*

## :warning: Prerequisites
<!-- Bullet list or simple sentence explaining what contributor needs for this project -->
- Add environment variables with DB logins
```bash
export DBUSER=your_dbuser
export DBPASS=your_dbpass
```

- Complete the DB using bash or directly in MySQL
```sql
CREATE DATABASE todolist;
USE todolist;
CREATE TABLE users (
  id         INT AUTO_INCREMENT NOT NULL,
  username      VARCHAR(255) NOT NULL,
  password     VARCHAR(255) NOT NULL, #! to modify for Security
  PRIMARY KEY (`id`)
);

INSERT INTO users
  (username, password)
VALUES
  ('john.doe', 'test123');
```

## :thinking: How does it run ?
*It's coming ...*

![Preview](https://placehold.co/500x300 "Preview of AppName")

<!-- 
    FEATURES
    List of the main new features, fixes to bring on the project

    Ex:
    - Setup Night/Day mode
    - Add animation when music is playing
-->

## :test_tube: Features (for v1)
- Security for password in DB (Hash)
- Complete the README.md

<!-- 
    LICENSE
    Write Developer name with used license
 -->
 
## :key: License

Developed by Loick Cherimont  

Under MIT License  

Last edition: 2024-08-17
