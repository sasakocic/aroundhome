# Assignment: Matching Customer & Partner

## Background

At Aroundhome our goal is to propose the right partner (a craftsman) to a customer based
on their project requirements. Matching of customers and partners is a crucial part in our
product. It determines how happy our customers will be with our partners and our partners
with the quality of the customer we connect them with.
The last product category that we reworked was flooring. The goal is to propose the right
partner based on the details of a customer's flooring project.

## Your Task

Your task is to write an API that offers the following functionality:
- Based on a customer-request, return a list of partners that offer the service. The list
should be sorted by “best match”. The quality of the match is determined first on
average rating and second by distance to the customer.
- For a specific partner, return the detailed partner data.
Matching a customer and partner should happen on the following criteria:
- The partner should be experienced with the materials the customer requests for the
project.
- The customer is within the operating radius of the partner.
    
The data in the request from the customer is:
- Material for the floor (wood, carpet, tiles)
- Address (assume that this is the lat/long of the house)
- Square meters of the floor
- Phone number (for the partner to contact the customer)

The structure of the partner data is as follows:
- Experienced in flooring materials (wood, carpet, tiles or any combination)
- Address (assume that this is the lat/long of the office)
- Operating radius (consider the beeline from the address)
- Rating (for this assignment you can assume that you already have a rating for a partner)

Please write the code in Go and generate some partner data for your challenge.
You can decide how you want us to test the solution. Eg. Providing a SwaggerUI for your API
endpoint, developing a simple UI or just a README file that shows how to call your service.

## How to submit the coding challenge

For submitting the coding challenge to us please create a repository in your GitHub account,
and share it with us once you have finished coding. Just send me the link, and I will forward it
to the team.

If you don’t want to have a public repo, you can also make it private and share it with two of
our team members. Just let me know, and I will share their details with you.

If you don’t use GitHub, you can also use another tool of your choice as long as you are able
to give access to us, or you can send a file with the code. You can choose whatever works
best for you as long as our team members have a chance to get to look at your code.

We will review the code, give feedback and ask some follow-up questions in the interviews
that follow. Have fun! We’re looking forward to having a great tech conversation with you! 😀

## System Design

Since the tasks requires calculating distance between two points given in (latitude,longitude)
we could use variants of [Haversine formula - Wikipedia](https://en.wikipedia.org/wiki/Haversine_formula) like

```sql
SELECT
    (
        6371 * acos(
            cos(radians(lat2)) * cos(radians(lat1)) * cos(radians(lng1) - radians(lng2)) + sin(radians(lat2)) * sin(radians(lat1))
        )
    ) as distance
from
    partners
```

If we use MySQL 5.7 then we have a function for this
```sql
select
  ST_Distance_Sphere(
    point(-87.6770458, 41.9631174),
    point(-73.9898293, 40.7628267)
  )
```

and we can always create a function to stay flexible on database choice like

```sql
DELIMITER $$ 
CREATE FUNCTION   `getDistance`(
    `lat1` VARCHAR(200),
    `lng1` VARCHAR(200),
    `lat2` VARCHAR(200),
    `lng2` VARCHAR(200)
) RETURNS varchar(10) CHARSET utf8
BEGIN
DECLARE distance varchar(10);
set
distance = (
    select 
      (
        6371 * acos(
          cos(
            radians(lat2)
          ) * cos(
            radians(lat1)
          ) * cos(
            radians(lng1) - radians(lng2)
          ) + sin(
            radians(lat2)
          ) * sin(
            radians(lat1)
          )
        )
      ) as distance
  );
if (distance is null) then return '';
else return distance;
end if;
END $$ DELIMITER;
```

So, what we need is to select all partners within the radius distance of a given customer point

## Model

To define the models we see clearly the format of customer request and partner data which should be in a database.
As we will query by distance, relational database can be used like MySQL, MariaDB or PostgreSQL.
MySQL has the aforementioned advantage or having a finished function, but we will use a function 
so that we are future-proof to replace MySQL with more robust PostgreSQL.
Materials on the other hand is clearly and array. For simplicity, we could have had separate columns for wood, carpet and tiles
which would made queries faster and simpler, but in practical usage these are expected to change.
In MySQL world sets can be used like [MySQL :: The MySQL SET Datatype](http://download.nust.na/pub6/mysql/tech-resources/articles/mysql-set-datatype.html)
In PostgreSQL this is a TEXT[] field which can be used as described in https://www.postgresql.org/docs/9.1/functions-array.html
For us, the contains operator @> is interesting as we want to search for all partners that contain the materials in the request.
Some potential usages of arrays are described here https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-array/

So we will have table partners with the minimal number of fields that are required to work
- id INTEGER
- name VARCHAR(256)
- lat FLOAT
- lng FLOAT
- radius
- rating
- flooring_experience TEXT []

As we will order the results by rating, it should be indexed on this field
id would be autonumber, not null primary id.

```sql
CREATE TABLE
    public.partners (
                        id serial NOT NULL,
                        name character varying(255) NOT NULL,
                        lat numeric NOT NULL,
                        lng numeric NOT NULL,
                        radius numeric NOT NULL DEFAULT 0,
                        rating double precision NULL,
                        flooring_experience text [] NULL
);

ALTER TABLE
    public.partners
    ADD
        CONSTRAINT partners_pkey PRIMARY KEY (id)
```

The function for PostgreSQL looks like

```sql
CREATE FUNCTION
  getDistance(
    lat1 DECIMAL,
    lng1 DECIMAL,
    lat2 DECIMAL,
    lng2 DECIMAL
  ) RETURNS DECIMAL AS $$

declare distance DECIMAL;

begin 
    select
      (
        6371 * acos(
          cos(radians(lat2)) * cos(radians(lat1)) * cos(radians(lng1) - radians(lng2)) + sin(radians(lat2)) * sin(radians(lat1))
        )
      ) INTO distance;
return distance;
end;

$$ language plpgsql
```

Dummy data generated with [Mockaroo - Random Data Generator and API Mocking Tool | JSON / CSV / SQL / Excel](https://www.mockaroo.com/) and imported with

```sql
COPY partners(id,name,lat,lng,radius,rating,flooring_experience)
FROM './db/partners.csv'
DELIMITER ','
CSV HEADER;
```

Or we can use CREATE and INSERT statements from `db/seed.sql`

This model would allow a queries like

```sql
select
    Id, Name, Lat, Lng, Radius, Rating, flooring_experience AS FlooringExperience,
    getDistance(40.076762, 113.300129, Lat, Lng) AS Distance
from
    partners
where
    getDistance(40.076762, 113.300129, Lat, Lng) < Radius AND flooring_experience @> ARRAY['carpet','tiles','wood']
order by
    Rating DESC,
    Distance;
```

Retrieving the partners along with the distance require the following types

```go
type Partner struct {
    Id                 int16
    Name               string
    Lat                float32
    Lng                float32
    Radius             float32
    Rating             float32
    FlooringExperience string
}

type PartnerWithDistance struct {
    Partner  Partner
    Distance float32
}
```

For the given sample data, result of the query 
http://127.0.0.1:3000/query/?address=40.076762,113.300129&material=carpet,tiles&phone=0160153700132&sqm=35
is

```json
{
  "partners":[
    {
      "Partner":{
        "Id":883,
        "Name":"Meevee",
        "Lat":39.296173,
        "Lng":113.6907,
        "Radius":127.98,
        "Rating":5.25,
        "FlooringExperience":"{carpet,tiles,wood}"
      },
      "Distance":93.00927
    },
    {
      "Partner":{
        "Id":1,
        "Name":"Lazz",
        "Lat":40.076763,
        "Lng":113.30013,
        "Radius":108.83,
        "Rating":0.96,
        "FlooringExperience":"{carpet,tiles}"
      },
      "Distance":0.0003288655
    }
  ],
  "phone":"0160153700132",
  "sqm":"35"
}
```

We see that carpet and tiles is required and more is welcomed.
The walking distance between 40.076762,113.300129 and 40.076762,113.300129 is 129km which is within 127.98 bee-line. 
according to https://www.google.com/maps/dir/'40.076762,113.300129'/39.296173,113.6907/@39.684718,112.8772452,9z/data=!3m1!4b1!4m7!4m6!1m3!2m2!1d113.300129!2d40.076762!1m0!3e2
[Pingcheng, Datong, Shanxi, China, 037048 to Fanshi County, Xinzhou, Shanxi, China, 034303 - Google Maps](https://www.google.com/maps/dir/'40.076762,113.300129'/39.296173,113.6907/@39.684718,112.8772452,9z/data=!3m1!4b1!4m7!4m6!1m3!2m2!1d113.300129!2d40.076762!1m0!3e2)

and the result of http://127.0.0.1:3000/partners/883 is

```json
{
   "Id":883,
   "Name":"Meevee",
   "Lat":39.296173,
   "Lng":113.6907,
   "Radius":127.98,
   "Rating":5.25,
   "FlooringExperience":"{carpet,tiles,wood}"
}
```

## Dependencies

We will use Fiber because of the extreme performance according to benchmarks [Fiber](https://gofiber.io/)
For production low memory footprint and rate limiter are especially important.

```sh
go get -u golang.org/x/lint/golint
go get -u github.com/gofiber/fiber/v2
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/arsmn/fiber-swagger/v2
go get -u github.com/lib/pq
```

## Swagger OpenAPI documentation

We use [GitHub - swaggo/swag: Automatically generate RESTful API documentation with Swagger 2.0 for Go.](https://github.com/swaggo/swag#declarative-comments-format)
The reason is that the declarative format is kept and updated within the code with less chances to diverge with time
It can be seen by using http://127.0.0.1:3000/swagger/index.html

## Environment

The following environment variables are used with defaults in parentheses:

- PORT to specify webserver port (3000)
- PG_HOSTNAME for PostgreSQL host (localhost)
- PG_USER for PostgreSQL user (postgres)
- PG_PASSWORD for PostgreSQL password (postgres)
- PG_DATABASE for PostgreSQL database name (aroundhome)

The environment for Docker can be copied from the .env.example to .env and adjusted.

Use [Use Docker Compose | Docker Documentation](https://docs.docker.com/get-started/08_using_compose/)

    docker-compose up

