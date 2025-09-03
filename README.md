# pqueue

[![GoDoc](https://godoc.org/github.com/AndrewChon/pqueue?status.png)](https://godoc.org/github.com/AndrewChon/pqueue)

A collection of concurrency-safe generic priority queue implementations.

> [!NOTE]
> This package is under development, and I plan to add more implementations in the future.

## About

Every priority queue (except Circular FIFO, obviously) in this package is a generic struct, where key _K_ has constraint
`cmp.Ordered` and value _V_ has constraint `any`. All priority queues are min-heap-based (i.e., the smallest key has the
highest priority) and allow for duplicate keys/priorities.

## Contributing

I am just one guy working on this in my free time for fun. So, if you have any suggestions or issues, please feel free
to open an issue or pull request!

## Implementations

Currently, this package provides the following queues.

### Time Complexities

| Type          | findMin | removeMin    | insert       | meld         |
|---------------|---------|--------------|--------------|--------------|
| Binary        | Θ(1)    | Θ(log n)     | Θ(log n)     | Θ(n)         |
| Circular FIFO | Θ(1)    | Θ(1)         | Θ(1)         | Θ(1)         |
| Pairing       | Θ(1)    | O(log n) am. | Θ(1)         | Θ(1)         |
| Skew          | Θ(1)    | O(log n) am. | O(log n) am. | O(log n) am. |
| Skew Binomial | Θ(1)    | Θ(log n)     | Θ(1)         | Θ(log n)     |

### Benchmarks

Apple MacBook Pro (M4 Pro, 24GB RAM)

| Type          | push        | meld                          | pop         |
|---------------|-------------|-------------------------------|-------------|
| Binary        | 49.40 ns/op | 17.72 ns/node                 | 502.7 ns/op |
| Circular FIFO | 24.46 ns/op | 3.1×10<sup>-6</sup> ns/node   | 5.945 ns/op |
| Pairing       | 48.17 ns/op | 4.83×10<sup>-4</sup> ns/node  | 728.0 ns/op |
| Skew          | 193.6 ns/op | 3.103×10<sup>-4</sup> ns/node | 587.9 ns/op |
| Skew Binomial | 62.05 ns/op | 6.102×10<sup>-4</sup> ns/node | 1070 ns/op  |
