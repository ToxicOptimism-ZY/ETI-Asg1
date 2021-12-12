<!-- TABLE OF CONTENTS -->
## Table Of Contents
<ol>
  <li>
    Introduction
  </li>
  <li>
    Design considerations
  </li>
  <li>
     Architecture diagram
  </li>
   <li>
    Passenger Microservice
  </li>
  <li>
    Driver Microservice
  </li>
  <li>
    Trip Microservice
  </li>
  <li>
    Setting up database
  </li>
  <li>
    Running the program
  </li>
</ol>



<!-- INTRODUCTION -->
## Introduction
Hi, I am Yap Zhao Yi, the developer of this repository, the codes provided are intended for Ngee Ann Polytechnic's Emerging IT Trends module, october semester 2021. The codes covers the creation of a ride hailing app. Additionally, as per the content page, explanation will be given for the design considerations, architecture diagram and various microservices functionality. 

<!-- DESIGN CONSIDERATIONS -->
## Design Considerations
Due to the nature of the assignment, the microservices are deployed locally without containerization, hence microservices while optimally displayed in as seperate packages are all written under main to allow visual studio code to recognize it as an executable. However, to account for this, go modules have been initalized as part of each microservice. Additionally, due to a lack of containerization and in turn a dedicated DBMS, a single database is also utilized despite seperate databases being more preferable when implementing micro services.

Through tactical domain driven design, persistant entities such as trips, passengers and drivers are recognized. It is also noted that no stand alone value objects have been identified. Due to the simplicity of the assignment, no further aggregation is neccessary. 

The simplificity of the assignment only requires for a single console app, rather than dedicated console micro services for the passenger and driver. Furthermore, it is notable that the all services are designed to be cohesive and loosely coupled, with much of the api http methods being for targeted purposes as well through the use of parameters. All http calls without the neccessary parameters will be rejected as a part of "404 Required parameters not found".

<!-- ARCHITECTURE DIAGRAM -->
## Architecture Diagram
<br />
<div align="center">
  <a href="https://github.com/ToxicOptimism-ZY/ETI-Asg1">
    <img src="architecture.png" alt="Logo" width="755" height="427">
  </a>
</div>
<br/>

<!-- Passenger Microservice -->
## Passenger Microservice

<ol>
    <li>
        POST: /api/v1/passengers <br />
        Create a new passenger account 
    </li>
    <li>
        GET: /api/v1/passengers?username={username}&password={password} <br />
        Gets a passenger account based on username and password provided
    </li>
    <li>
        GET: /api/v1/passengers/{passengerid} <br />
        Gets passenger account with the passenger id provided
    </li>
    <li>
        PUT: /api/v1/passengers/{[passengerid} <br />
        Updates passenger account with the passenger id provided
    </li>
    <li>
        DELETE: /api/v1/passengers/{passengerid} <br />
        Deletes passenger account with the passenger id provided
    </li>
 </ol>
 
 <!-- Driver Microservice -->
## Driver Microservice

<ol>
    <li>
        POST: /api/v1/drivers <br />
        Create a new driver account <br />
    </li>
    <li>
        GET: /api/v1/drivers?username={username}&password={password} <br />
        Gets a driver account based on username and password provided <br />
    </li>
     <li>
        GET: /api/v1/drivers?status={status} <br />
        Gets all driver accounts with the status provided <br />
    </li>
    <li>
        GET: /api/v1/drivers/{driverid} <br />
        Gets driver account with the driver id provided <br />
    </li>
    <li>
        PUT: /api/v1/drivers/{driverid} <br />
        Updates driver account with the driver id provided <br />
    </li>
    <li>
        DELETE: /api/v1/drivers/{driverid} <br />
        Deletes driver account with the driver id provided <br />
    </li>
 </ol>
 
<!-- Trip Microservice -->
## Trip Microservice

<ol>
    <li>
        POST: /api/v1/trips <br />
        Create a new trip record <br />
    </li>
    <li>
        GET: /api/v1/trips?passengerid={passengerid}&status={status} <br />
        Gets a passenger's current active trip record <br />
    </li>
    <li>
        GET: /api/v1/trips?driverid={driverid}&status={status} <br />
        Gets a driver's current active trip record <br />
    </li>
    <li>
        GET: /api/v1/trips?passengerid={passengerid} <br />
        Gets a passenger's trip history (in chronological order)
    </li>
    <li>
        GET: /api/v1/trips/{tripid} <br />
        Gets trip record with the trip id provided <br />
    </li>
    <li>
        PUT: /api/v1/trips/{tripid} <br />
        Updates trip record with the trip id provided <br />
    </li>
    <li>
        DELETE: /api/v1/trips/{tripid} <br />
        Delete trip account with the trip id provided <br />
    </li>
 </ol>
 
 <!-- SETTING UP DB -->
## Setting up database
<ol>
    <li>
        Using sql workbench, create a new connection of Hostname: "127.0.0.1" and Port: "3306"
    </li>
    <li>
        Run the query located in the file "SetupAuthorizedUser".sql
    </li>
    <li>
        Run the query located in the file "CreateDatabase".sql, note that the 3 insert queries are intended to populate the database with sample data and is hence not neccessary.
    </li>
 </ol>
