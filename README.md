# Local Upstash Proxy Server

Upstash redis uses a serverless driver that lets your app connect to it via http instead of tcp to be optimized for serverless environments where you want to be using http to communicate with your redis server, this is all good but the problem is that if you want to use a local development server the basic `node-redis` or equivalent connector uses tcp, causing there to be some potential inconsistencies, this proxy server is meant to sit in front of your redis server and allow you to communicate with your local redis server via http for consistent behaviour between dev & prod 

theres some other solutions to this problem that are probably better and more robust than this one is but I just wanted to build this out, thx
