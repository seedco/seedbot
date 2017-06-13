# Seed-Go
A Go client for the Seed API

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/seedco/seed-go) [![Build Status](https://circleci.com/gh/seedco/seed-go.svg?style=shield)](https://circleci.com/gh/seedco/seed-go)

## Usage

```go
// to obtain an access token go to https://api.seed.co/v1/public/auth/token in a browser and enter in your seed username/password

accessToken := "1.iap2H-4qQ-WR9sy55555uaytQ.o5A32LYL5-87a_60kcQiX1Lp878GVbx8xfVvTfp5tpc.orsHbAqao-5KfsH8SdglQFltK7Ii8ktL7xo8tls3HAB"

client := seed.New(accessToken)

getTransactionsReq := TransactionsRequest{
	Client: client,
}
// The two options are between getting all transactions or using an iterator for pagination

// All transcations

transactions, err := getTransactionsReq.Get()

// Using an iterator for pagination

iterator := getTransactionsReq.Iterator()

iterator.SetBatchSize(10)

for iterator.HasNext() {
	var transactions []seed.Transaction
	var err error
	if transactions, err = iterator.Next(); err != nil {
		panic(err.Error())
	}
	fmt.Printf("Transactions:\n%v", transactions)
}

// previous will get the previous page of transactions

previousTransactions, err = iterator.Previous()
```
