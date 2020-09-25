# Controllers vrs Services 


https://github.com/System-Glitch/goyave/#controller

This reads of what I call a Service. A Service is reusable buisness logic, for instance the login service. While a controller is something like a front controller where it proxies all requests to services based on the post path. One might think that's a router. Well, the controller would control or inject logic before the router, which calls the service.

Thus I am calling this concept a service.

### Something I found on the web 
e.g. https://stackoverflow.com/questions/3295267/where-does-my-code-go-controller-service-or-model#:~:text=Use%20a%20service%20rather%20than,avoid%20referencing%20application%2Dspecific%20functionality.&text=If%20a%20service%20is%20not,this%20application%20use%20a%20controller.

Model: (from Wikipedia - MVC) "The Model is used to manage information and notify observers when that information changes; it's also a domain-specific representation of the data upon which the application operates." To my mind this implies properties and the like - not methods.

Controller: (from Wikipedia - MVC) "Receives input and initiates a response by making calls on model objects. A controller accepts input from the user and instructs the model and viewport to perform actions based on that input."

Service: There are many different opinions on what a service is, I assume in your context a service is: an externally facing callable point (within the context of the layers of your system) that provides a specific answer to a specific question. (Services being usually based around business concepts not technical ones.)