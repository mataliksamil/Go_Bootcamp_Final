## Property Finder Go Bootcamp Final Projects
## Shoping Cart REST API 
___________

### Tools 
- gin-gonic framework 
- go-pg ORM 
- Postgresql


The project includes 4 tables which are 

- Baskets : has-one relation with users, has-many relation with basketproducts
Basketstatus field introduced to seperate current basket with the ones belongs to previous orders
Status togles to zero after complete the order command.

- Users : Has many relation with baskets

- Products : 

- BasketProducts:has-one relation with basket, has-many relation with product

The project does not have Order or OrderItem tables so, a user can have many baskets which can be thought as previous orders, 
this decision which I made at the begining of my work, caused the structure to be slightly more intricate inevitably. 


### Test
- The project does not includes unit-tests
- A postman request collection added to main directory
Due to relations in db Requests can be followed as :

- Add Products > Add User > Add Basket > Add products to basket > Show active basket > complete the order 
