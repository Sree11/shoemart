# shoemart
Shoe Mart Search Engine
SHOEMART

This is a search engine built with GOLang and elastic search as the backend 
I used a colorlib template called shoppers for frontend UI, mostly modified and stripped off all the features, just for login and product page

The  login page with username "admin" password : "user"

once you are authenticated it will redirect you to localhost:8000/v1/products
This is a UI search where you can search using the HTML page search fields 

First field takes a search string, second field takes filter params, third takes sort , forth is for Limit and fifth field is for offset

The same search can be used through the URL using params as below 
from the endpoint localhost:8000/v1/search
e.g. 
http://localhost:8000/v1/search?q=shoe&limit=20&offset=1&sort:price:desc&brand=adidas

The dataset has around 2000 records of raw data

The elastic search folder has the shoe-prices-QueryResult-mini.csv in CSV and shoedata.json in json format and which contains the data

the index created in elastic search is shoemart  - localhost:9200/shoemart

the Web app looks for the index shoemart 

I have created a docker image for both golang webapp and the elasticsearch app, unfortunately there is a minor issue, as the golang app is  unable to connect to the elasticsearch app in the docker environment, some online forums have referred to a bug, the docker composer will bring up both the images, but I suggest if you could just bring up the elastic search docker 
Create an index with the name shoemart and load the csv date through kibana that would work 

On of the two will work for sure but both apps in docker is not working, I tried fixing it but was out of time, you can run the build on golang app and run it, it should run on port 8000 localhost:8000/

Apologize that I ran into these docker issues, tried really hard to fix it but will try again to get it right.

Please do let me know if you need any info on brining up the application would be glad to help out 
Again thanks a lot for giving me this opportunity to work on the problem solution

Thanks,
Sree




