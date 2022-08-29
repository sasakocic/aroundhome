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
and share it with us once you have finished coding. Just send me the link and I will forward it
to the team.

If you don’t want to have a public repo, you can also make it private and share it with two of
our team members. Just let me know and I will share their details with you.

If you don’t use GitHub, you can also use another tool of your choice as long as you are able
to give access to us or you can send a file with the code. You can choose whatever works
best for you as long as our team members have a chance to get to look at your code.

We will review the code, give feedback and ask some follow up questions in the interviews
that follow. Have fun! We’re looking forward to having a great tech conversation with you! 😀