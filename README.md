# MP3 Web Crawler
Crawl mp3 links from web


#To read
https://jdanger.com/build-a-web-crawler-in-go.html

https://acoustid.org/chromaprint
https://acoustid.org/webservice



Should work like http://beemp3s.org/


64 KB of mp3 file is enought to get all id tags info
 
 
 
 # How it works
 
 1) scan internet for mp3 links
 2) download mp3 files and get their acoustId
 
 
 # Requirements:
 1. Manages request delays and maximum concurrency per domain
 2. Cookie and session handling
 3. Max depth
 