# docker-graph
Show dependency graph of docker images/containers like this:

![image](https://user-images.githubusercontent.com/20103016/152722948-52340f5c-acad-4db3-aecb-48f558e8bcea.png)

Orange is images and green is containers.

# Features
Collect docker images, containers information, and show dependency graph using by Cytoscape.js.

This program calls docker command to get info, user need permission to use docker command but that's all. You don't need any other settings/configs.

# Install/usage
Simply copy docker-graph command to docker host.

$ ./docker-graph

access from browser http://SERVER-IP:9091/

Caution!: There are no access restrictions. It must not to use production.
  
# Build

Need golang 1.16 or later.
  
$ go build .

# License
GPLv3
# Author
Tomohisa Hirami(hirami@io-architect.com)
