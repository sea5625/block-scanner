ISAAC
============

What's the ISAAC
--------
 The professional monitoring tool of loopchain network. 

Dependencies
---------
- Back-end: Go (over 1.11.5)
- Front-end: Node (over v11.3.0), npm (over 6.4.1), and yarn (1.12.3)


How to build
--------
```
make
```

How to configure and run
--------
Complete build and configure before Run.

1. Setup ISAAC.
    - Configure node, channel, prometheus, etc in ```config/configuration.yaml``` like following.

    ``` yaml
    node :
        - name : node0
          ip : https://int-test-ctz.solidwallet.io
        - name : node1
          ip : https://int-test-ctz.solidwallet.io
        - name : node2
          ip : https://int-test-ctz.solidwallet.io
    
    channel :
        - name : "default"
          nodes : [node0, node1, node2]
    
    prometheus :
        prometheusExternal : http://localhost:9090  # use ISAAC front-end
        prometheusISAAC : http://localhost:9090     # use ISAAC server 
        queryPath: /api/v1/query
        crawlingInterval: 2
        nodeType: loopchain     # Can be used Node type (loopchain, goloop)
        jobNameOfgoloop : goloop  # If using the goloop, should set the job name of prometheus.
    
    etc :
        sessionTimeout: 30
        language: ko
        loglevel: 1
        db:
          - type: sqlite3
            id: ""
            pass: ""
            database: ""
            path: data/isaac.db
        loginLogoImagePath: images/iconloop.png
            
    blockchain:
      crawlingInterval: 0
      db:
        - type: sqlite3
          id: ""
          pass: ""
          database: ""
          path: data/crawling_data.db
          
    authorization:
      thirdPartyUserAPI: [channels, nodes, blocks, txs]
    ```
2. DB Setting
   - target
       - etc : isaac server used database
       - blockchain : block crawling used database   
   - type
       - sqlite3
          - Edit ```/config/configuration.yaml```. file
          
                  ``` yaml
                  ....
                  ....
                  db:                  
                    - type: sqlite3
                      id: ""
                      pass: ""
                      database: ""
                      path: data/isaac.db
                  ```
       - mysql 
          - Edit ```/config/configuration.yaml```. file (ex:)
                
                  ``` yaml
                  ....
                  ....
                  db:                  
                    - type: mysql
                      id: block
                      pass: password
                      database: isaac
                      path:
                  ```
    
3. Setting logo image in login page

   - Edit etc.loginLogoImagePath in ```/config/configuration.yaml```. file (ex:)
    
        !!! Do not change folder path. Do can change only file name.
        
        ```
        yaml
        ....
        etc :
            ....
            loginLogoImagePath: images/iconloop.png
        ....
        ```

4. Setting authorized API list for third party user
   - Can be used API list is channels, nodes, blocks, txs.
        ``` yaml
         ....
         authorization:
               thirdPartyUserAPI: [channels, nodes, blocks, txs]
         ....
        ```
         
5. Run loopchain exporter and  prometheus

   - Run ```run_prom.sh``` below ```prometheus_local```. 

6. Run ISAAC.

    ```
    make run
    ```

7. Open the browser and connect ```http://localhost:6553```.

8. API Document page ```http://localhost:6553/swagger/index.html ```.

9. Log Debugging Mode
    - Edit ```/config/configuration.yaml```. file
    
        ``` yaml
        ....
        ....
        etc :
            ....
            loglevel: 2   //1:Info(default), 2:Debug
            ....
        ```

Using Docker
------

### Build
``` bash
$ make docker-build   #Generate iconloop/isaac:${VERSION} and iconloop/isaac:latest images.
```
To change version, change ```VERSION``` file. 

### Build develop version. 
``` bash
$ make docker-build-dev  #Generate iconloop/isaac:${GIT_SHA1} and iconloop/isaac:dev images.
```

### Run with docker.   
See ```rum_prom_and_isaac.sh``` under ```prometheus_local``` folder.