
Name - Omkar Podey
------------------------------------------------------------------------------------------------------------------------
Assignment - Programming assignment 1 (Http server supporting http/1.0 and http/1.1)
------------------------------------------------------------------------------------------------------------------------
What does the program do?

    1. It starts listening to incoming connections on the port(sent from the terminal) as it's now bound to the server.
    2. Runs an infinite for loop to accept all incoming connections.
    3. Using go routines for func handleConnection(), spawns a new thread for every connection.
    4. Function handleConnection() for every new connection if it's HTTP/1.1 it runs an infinite for loop to keep the
        connection open until the dynamic timeout or else close connection after serving.
    5. Dynamic timeout- It just uses the number of connections to determine the timeouts, basically the higher the
        number of connections the lower the timeout.
    6. It extends the connection timeout every time a new request comes in within the previous timeout to keep the
        persistent tcp connection open.
    7. We call func processRequest() inside the handleConnection to check if the request is properly formed, then pass
        it to func serveFile() to actually transmit the file.
    8. Serve file checks the file permissions, if it exists and send errors with the appropriate headers (also adds
        keep-alive if its http/1.1), if the file exists we send the files in chunks of 8kb using the buffer.
------------------------------------------------------------------------------------------------------------------------
Extra implementations
    1. Buffering and sending chunks of data.
    2. Printing statements to better understand the server working.
    3. Cleaning the file paths before processing.
    4. Check for HOST header.
------------------------------------------------------------------------------------------------------------------------
List of submitted files
    1. Server source file - (http_web_server.go)
    2. Readme.txt (with all the information  and examples)
    3. Screenshots folder with the server being accessed form the browser.
------------------------------------------------------------------------------------------------------------------------
How to run this ?
    1. You need GO on your machine (if you have brew then just "brew install go")
    2. Just run it from the terminal
        The command - "go run http_web_server.go -document_root /Users/omkarpodey/projects/go/distributed_assign/www.sjsu.edu -port 8080"
------------------------------------------------------------------------------------------------------------------------
How did I test it ?

1. For Http/1.0 used this "curl --http1.0 http://localhost:8080/index.html http://localhost:8080/visit/index.php.html"

    The output -
            Handling connection for client: [::1]:60742
            [Thu, 26 Sep 2024 12:35:37 PDT] Client [::1]:60742 requested GET /index.html
            Closing connection for client: [::1]:60742
            The number of connections 1
            Handling connection for client: [::1]:60743
            [Thu, 26 Sep 2024 12:35:37 PDT] Client [::1]:60743 requested GET /visit/index.php.html
            Closing connection for client: [::1]:60743

2. For Http/1.1 I just used the firefox browser (Chrome acts weird at times), just load localhost:8080 (or the port you
    run it on)

    Looking at the output below you can observe that only 6 connections were opened while over 29 requests were made.
    So we can conclude that http/1.1 is working as expected and even times-out in the end as you can see in the output.

    The output -
            Serving files from /Users/omkarpodey/projects/go/distributed_assign/www.sjsu.edu on port 8080...
            The number of connections 1
            The number of connections 2
            The number of connections 3
            The number of connections 4
            The number of connections 5
            The number of connections 6
            Handling connection for client: 127.0.0.1:60776
            Handling connection for client: 127.0.0.1:60772
            Handling connection for client: 127.0.0.1:60771
            Handling connection for client: 127.0.0.1:60774
            Handling connection for client: 127.0.0.1:60775
            Handling connection for client: 127.0.0.1:60773
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60771 requested GET /
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60772 requested GET /aspis/css/foundation.css
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /aspis/css/prototype.css
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60775 requested GET /aspis/css/campus.css
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60772 requested GET /aspis/js/vendor/jquery.js
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60775 requested GET /aspis/js/vendor/what-input.js
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /aspis/js/vendor/foundation.js
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60773 requested GET /aspis/js/vendor/jquery-accessibleMegaMenu.js
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60776 requested GET /aspis/js/vendor/direct-edit.js
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60772 requested GET /_images/news/news-091824-cids-hero.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /_images/sjsu-homepage-hero/2024_wsj_rankings_homepage_3800x1920_v2.mp4
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60775 requested GET /_images/sjsu-rankings/ranking-wsj-overall-university.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60776 requested GET /_images/news/featured_calendar_default-spirit-mark_121520.jpg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60771 requested GET /_images/news/wsq-entrepreneurship-feature-ray-zinn-dschmitz-101518-50_img-web.jpg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60773 requested GET /_images/sjsu-homepage-hero/Homepage__3800x1935-hispanic-heritage_2024.jpg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60773 requested GET /_images/news/news-091824-cob-fortune-mba.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60772 requested GET /_images/news/SJSU_Online_jgensheimer_052523_1396_IMG_WEB.jpg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60775 requested GET /_images/news/homepage_WSQ-S2024_IMG_WEB.jpg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60773 requested GET /_images/sjsu-homepage-hero/2024_wsj_rankings_homepage_3800x1920.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60775 requested GET /aspis/media/brand/logo-we-are-spartans.svg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60772 requested GET /_images/sjsu-rankings/ranking-wsj-public-university.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60773 requested GET /_images/sjsu-rankings/ranking-usn-social-mobility.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /_images/news/news-091824-banned-books.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60773 requested GET /aspis/media/brand/icon-spartan-withbackground.svg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /_images/sjsu-rankings/ranking-usn-public-university.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /_images/sjsu-rankings/ranking-money-transformative-university.png
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /_images/news/news-091824-politics-in-progress.jpg
            [Thu, 26 Sep 2024 12:38:11 PDT] Client 127.0.0.1:60774 requested GET /aspis/media/brand/logo-sjsu.svg
            Read timeout occurred: read tcp 127.0.0.1:8080->127.0.0.1:60776: i/o timeout
            Read timeout occurred: read tcp 127.0.0.1:8080->127.0.0.1:60771: i/o timeout
            Read timeout occurred: read tcp 127.0.0.1:8080->127.0.0.1:60772: i/o timeout
            Read timeout occurred: read tcp 127.0.0.1:8080->127.0.0.1:60775: i/o timeout
            Read timeout occurred: read tcp 127.0.0.1:8080->127.0.0.1:60773: i/o timeout
            Read timeout occurred: read tcp 127.0.0.1:8080->127.0.0.1:60774: i/o timeout
------------------------------------------------------------------------------------------------------------------------



