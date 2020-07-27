# Payment

## Introduction

This package is built for the sole purpose of card payments.

## Usage

### Install Package

```bash
go get github.com/ichtrojan/payment
```

### Paystack

* Initialise Paystack

```go
package main

import (
	"github.com/ichtrojan/payment"
)

var paystack = payment.Paystack("sk_test_000000000000000000000000000000")
...
```

>**NOTE**<br/>
>Ensure you pass your paystack secret key

* Initiate charge

```go
...

func main () {
    initiateCharge, err := paystack.InitiateCharge("exmaple@domain.com", "chop_life_01")
    	
    if err != nil {
        log.Println(err)
    }

    fmt.Printf("%+v\n", initiateCharge)
}
```

* Verify charge

```go
...

func main () {
    verifyCharge, err := paystack.VerifyCharge("chop_life_01")
    
    if err != nil {
        log.Println(err)
    }

    fmt.Printf("%+v\n", verifyCharge)
}
```

* Charge charge

```go
...

func main () {
    chargeCard, err := paystack.ChargeCard("AUTH_qeut4h3xfn", "exmaple@domain.com", 9000)
    
    if err != nil {
        log.Println(err)
    }

    fmt.Printf("%+v\n", chargeCard)
}
```
* Fetch Transaction

```go
...

func main () {
    transaction, err := paystack.FetchTransaction(292584114)
    
    if err != nil {
        log.Println(err)
    }

    fmt.Printf("%+v\n", transaction)
}
```


>**NOTE**<br/>
>Check the `example` directory to see a sample implementation

### Flutterwave

comming soon....

## Contributors

* Deji Ajibola - [Twitter](https://twitter.com/damndeji) [GitHub](https://github.com/youthtrouble)

