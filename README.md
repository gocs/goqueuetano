# goqueuetano
 queueing web app


addform:
    creates new customer
    assigns new customer to available worker
    sets customer state to running
    sets incoming states to paused

editform:
    sets updated customer to paused
    modifies or updates customer
    updates customer worker
    sets updated customer to running

deleteform:
    sets the customer state

## Running

```
make run
```