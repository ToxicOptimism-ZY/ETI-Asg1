CREATE database asg1;

USE asg1;

CREATE TABLE Passenger (
    PassengerID int NOT NULL AUTO_INCREMENT,
    Username varchar(20) NOT NULL,
	`Password` varchar(20) NOT NULL,
	FirstName varchar(50) NOT NULL,
    LastName varchar(50) NOT NULL,
    MobileNo char(8) NOT NULL,
    EmailAddress varchar(100) NOT NULL,
    PRIMARY KEY (PassengerID),
    UNIQUE (Username),
    UNIQUE (EmailAddress)
);

CREATE TABLE Driver (
    DriverID int NOT NULL AUTO_INCREMENT,
    Username varchar(20) NOT NULL,
    `Password` varchar(20) NOT NULL,
	FirstName varchar(50) NOT NULL,
    LastName varchar(50) NOT NULL,
    MobileNo char(8) NOT NULL,
    EmailAddress varchar(100) NOT NULL,
    NRIC char(9) NOT NULL,
    CarLicencePlate char(8) NOT NULL,
    `Status` varchar(50) NOT NULL,
    PRIMARY KEY (DriverID),
    UNIQUE (Username),
    UNIQUE (EmailAddress),
    UNIQUE (NRIC), 
    UNIQUE (CarLicencePlate)
);

CREATE TABLE Trip (
    TripID int NOT NULL AUTO_INCREMENT,
    PickUp char(6) NOT NULL,
    DropOff char(6) NOT NULL,
    DriverID int NOT NULL,
    PassengerID int NOT NULL,
    `Status` varchar(50) NOT NULL,
    PRIMARY KEY (TripID)
);

INSERT INTO Passenger (Username, `Password`, FirstName, LastName, MobileNo, EmailAddress) VALUES 
('Timmy', 'Password', 'Timmy', 'Lim', '82715301', 'timmy@np.edu.sg');

INSERT INTO Driver (Username, `Password`, FirstName, LastName, MobileNo, EmailAddress, NRIC, CarLicencePlate, Status) VALUES 
('Bob', 'Password', 'Bob', 'Lim', '84600000', 'bob@np.edu.sg', 's0123456t','shx0478e', 'Available');

INSERT INTO Trip (PickUp, DropOff, DriverID, PassengerID, `Status`) VALUES 
('670123', '609731', '1', '1', 'Completed');
