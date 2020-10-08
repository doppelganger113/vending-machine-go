# vending-machine-go
Console application that implements vending machine logic written in Go

## Description

Vending machines have an array of buckets where each bucket contains a 
number of (possibly different products). When "vending" a particular 
product from a particular bucket, only the front-most product can be 
vended. When a user generates an order containing a variety of 
different products, the system needs to run an algorithm to decide 
whether the products can be sold and what the buckets look like after 
they have been sold.

We encode a single bucket by a string as follows: 
<int_1>,<int_2>,...,<int_n> where the front-most <int_1> has to be 
vended first. So e.g. 1,2,1,3,4 would be a bucket in which product 1 
can be sold first, then product 2, then again product 1 and so on.

We encode the whole bucket array by concatenating several bucket 
encodings with ; - so e.g. 1,2,1,3,4;5,2,3,3,3 would be a two-bucket 
configuration where product 1 can be vended immediately from the first 
bucket and product 5 can be vended immediately from the second bucket. 
A user order is encoded by <int_1>,<int_2>,...,<int_m> where the order 
of product ids does not matter, e.g. 5,2,2 is the same as 2,2,5.

The algorithm now receives a bucket array as well as a user order as 
an input and either outputs "IMPOSSIBLE" or outputs a new bucket array 
after the order has been vended.

Program receives both inputs via command line and 
output the result to STDOUT.
Please test your program with the following input:
Buckets 1,2,3,5,5;2,5,4,3,1;3,5,4,1,1;5,1,1,1,1
Order 1,2,3,4,5

## Usage

### Build
Build the executable with
```bash
go build
```
### Run

Strict is false by default and refers to strict product order popping
`cmd -strict=boolean <products> <buckets>` strict defaults to false.

```bash
./vending-machine-go "1,2,3,4,5" "1,2,3,5,5;2,5,4,3,1;3,5,4,1,1;5,1,1,1,1"
```

## Testing
```bash
go test ./internal
```