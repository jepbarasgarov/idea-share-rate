--  BELLI-ERDU DB SCHEME

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


---DROP-CONSTRAINTS---


ALTER TABLE tbl_idea DROP CONSTRAINT fk_worker;
ALTER TABLE tbl_link DROP CONSTRAINT fk_idea;
ALTER TABLE tbl_sketch DROP CONSTRAINT fk_idea;
ALTER TABLE tbl_user_idea_rel DROP CONSTRAINT fk_user;
ALTER TABLE tbl_idea_rate DROP CONSTRAINT fk_idea;
ALTER TABLE tbl_idea_rate DROP CONSTRAINT fk_criteria;
ALTER TABLE tbl_idea_rate DROP CONSTRAINT fk_user;







---DROP-TABLES---

DROP TABLE IF EXISTS tbl_user;
DROP TABLE IF EXISTS tbl_worker;
DROP TABLE IF EXISTS tbl_genre;
DROP TABLE IF EXISTS tbl_position;
DROP TABLE IF EXISTS tbl_mechanic;
DROP TABLE IF EXISTS tbl_link;
DROP TABLE IF EXISTS tbl_sketch;
DROP TABLE IF EXISTS tbl_idea;
DROP TABLE IF EXISTS tbl_idea_rate;
DROP TABLE IF EXISTS tbl_criteria;
DROP TABLE IF EXISTS tbl_user_idea_rel;


 --------USER--------------

CREATE TABLE tbl_user (
	id              UUID			  DEFAULT uuid_generate_v4 (),
	username        VARCHAR(32)  	  NOT NULL,
	password        VARCHAR(64)	      NOT NULL,
    firstname       VARCHAR(32)       NOT NULL,
    lastname        VARCHAR(32)       NOT NULL,
	role	        VARCHAR(16)       NOT NULL,
    status          VARCHAR(32)       NOT NULL,
	create_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    update_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT username_unique UNIQUE (username),

    PRIMARY KEY (id)

);


-------WORKER------------------------

CREATE TABLE tbl_worker(
    id              UUID			  DEFAULT uuid_generate_v4 (),
    firstname        VARCHAR(64)       NOT NULL,
    lastname        VARCHAR(64)       NOT NULL,
    position        VARCHAR(128)       NOT NULL,
    create_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    update_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),

    PRIMARY KEY (id)
);


------------GENRE-------------------------

CREATE TABLE tbl_genre(
    name            VARCHAR(256)      NOT NULL,
    CONSTRAINT genre_unique UNIQUE (name)
);

-----------JOB------------------------------------------

CREATE TABLE tbl_position(
    name            VARCHAR(256)      NOT NULL,
    CONSTRAINT position_unique UNIQUE (name)

);

-------Mechanics---------------------------------------------

CREATE TABLE tbl_mechanic(
    name            VARCHAR(256)      NOT NULL,
    CONSTRAINT mechanic_unique UNIQUE (name)

);

-------------CRITERIA----------------------------------------------

CREATE TABLE tbl_criteria(
    id              UUID			  DEFAULT uuid_generate_v4 (),
    name            VARCHAR(256)      NOT NULL,
    create_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    update_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT criteria_unique UNIQUE (name),
    PRIMARY KEY (id)
);

---------------IDEA---------------------------------------------------

CREATE TABLE tbl_idea(
    id              UUID			  DEFAULT uuid_generate_v4 (),
    name            VARCHAR(256)      NOT NULL,
    worker_id       UUID              NOT NULL,
    date            DATE              NOT NULL,
    genre           VARCHAR(256)      NOT NULL,
    mechanics       VARCHAR(256)[]    DEFAULT ARRAY[]::VARCHAR(256)[],
    description     TEXT              NOT NULL,
    create_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    update_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    
    PRIMARY KEY (id)


);

--------------LINK---------------------------------------------------

CREATE TABLE tbl_link(
    id               UUID			          DEFAULT uuid_generate_v4 (),
    label            VARCHAR(128)             NOT NULL,
    link             TEXT                     NOT NULL,
    idea_id          UUID                     NOT NULL
);

--------------SKETCH---------------------------------------------------

CREATE TABLE tbl_sketch(
    id               UUID			  DEFAULT uuid_generate_v4 (),
    name             VARCHAR(256)             NOT NULL,
    idea_id          UUID                     NOT NULL,
    place            INT                      NOT NULL,
    file_path        VARCHAR(128)             NOT NULL
);

----------------IDEAR RATE---------------------------------------------


CREATE TABLE tbl_idea_rate(
    idea_id              UUID                     NOT NULL,
    user_id              UUID                     NOT NULL,
    criteria_id          UUID                     NOT NULL,
    rate                 INT                      NOT NULL,
    create_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    update_ts  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    


    CONSTRAINT user_rate  UNIQUE (idea_id, user_id, criteria_id)
);
----------------------------------------------------------------------------------------------------------------------------------------------------------------


CREATE TABLE tbl_user_idea_rel(
  idea_id              UUID                     NOT NULL,
  user_id              UUID                     NOT NULL,
  mark                 VARCHAR(16)              DEFAULT 'RATED',

  CONSTRAINT user_idea_rel  UNIQUE (idea_id, user_id)


);

--CONTRAINTS

ALTER TABLE tbl_user_idea_rel ADD CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES tbl_user(id) ON DELETE CASCADE;
ALTER TABLE tbl_user_idea_rel ADD CONSTRAINT fk_idea FOREIGN KEY(idea_id) REFERENCES tbl_idea(id) ON DELETE CASCADE;
ALTER TABLE tbl_idea ADD CONSTRAINT fk_worker FOREIGN KEY(worker_id) REFERENCES tbl_worker(id) ON DELETE RESTRICT;
ALTER TABLE tbl_link ADD CONSTRAINT fk_idea FOREIGN KEY(idea_id) REFERENCES tbl_idea(id) ON DELETE CASCADE;
ALTER TABLE tbl_sketch ADD CONSTRAINT fk_idea FOREIGN KEY(idea_id) REFERENCES tbl_idea(id) ON DELETE CASCADE;
ALTER TABLE tbl_idea_rate ADD CONSTRAINT fk_idea FOREIGN KEY(idea_id) REFERENCES tbl_idea(id) ON DELETE CASCADE;
ALTER TABLE tbl_idea_rate ADD CONSTRAINT fk_criteria FOREIGN KEY(criteria_id) REFERENCES tbl_criteria(id) ON DELETE CASCADE;
ALTER TABLE tbl_idea_rate ADD CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES tbl_user(id) ON DELETE CASCADE;








-----INSERTIONS     

--User
INSERT INTO tbl_user(firstname, lastname, username, password, role, status) VALUES 
('Dayanch', 'Bababyev', 'dayanch.b', '$2a$12$BWKzPWJXBqhHuvhKlCCYQ.tbmBCJq3IULfqBPjMHSrZeJnLxk/Fhq', 'USER', 'ACTIVE'),
('Ybrayym' ,'Dathudayev', 'ybrayym.d' ,'$2a$12$BWKzPWJXBqhHuvhKlCCYQ.tbmBCJq3IULfqBPjMHSrZeJnLxk/Fhq', 'USER' , 'ACTIVE'),
('Tilla' ,'Bozaganova', 'tilla.b' ,'$2a$12$BWKzPWJXBqhHuvhKlCCYQ.tbmBCJq3IULfqBPjMHSrZeJnLxk/Fhq', 'USER' , 'ACTIVE'),
('Kerim' ,'Aganiyazov', 'kerim.a' ,'$2a$12$BWKzPWJXBqhHuvhKlCCYQ.tbmBCJq3IULfqBPjMHSrZeJnLxk/Fhq', 'USER' , 'ACTIVE'),
('admin' ,'admin', 'admin' ,'$2a$12$BWKzPWJXBqhHuvhKlCCYQ.tbmBCJq3IULfqBPjMHSrZeJnLxk/Fhq', 'ADMIN' , 'ACTIVE');





INSERT INTO tbl_position(name) VALUES('Game developer'),('Junior developer'),('3D designer'),('Artist'), ('Project manager'), ('Game designer');
INSERT INTO tbl_genre(name) VALUES('Hyper casual'),('Sandbox'), ('Real-time startegy');
INSERT INTO tbl_mechanic(name) VALUES('Turning Mechanics'),('Dexterity mechanincs'),('Swerve mechanics'), ('Rising/Falling mechanics'), ('Merging mechanics');
INSERT INTO tbl_criteria(name) VALUES('Time management'),('Simillarity to trends'),('Adding different features'), ('Our knowledge to make a game');

INSERT INTO tbl_worker(firstname, lastname, position) VALUES
('Kadyrberdi', 'Annabayev','Game developer'),
('Yhlas', 'Kerimov','Game developer'),
('Dowletgeldi', 'Rozymuradov','Game developer'),
('Tilla', 'Bozaganova','Project manager'),
('Ybrayym', 'Dathudayev','Game designer'),
('Bayramsoltan', 'Soyunova','3D designer'),
('Rahmet', 'Muhtarkulyyev','3D designer'),
('Salyh', 'Babayev','Game developer'),
('Shanazar', 'Geldiyev','Game developer'),
('Oraz', 'Esenov','Game developer');







