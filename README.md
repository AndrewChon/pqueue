# pqueue

[![GoDoc](https://godoc.org/github.com/AndrewChon/pqueue?status.png)](https://godoc.org/github.com/AndrewChon/pqueue)

A collection of concurrency-safe generic priority queue implementations.

> [!NOTE]
> This package is under development, and I plan to add more implementations in the future.

## About

Every priority queue in this package is a generic struct, where key _K_ has constraint `cmp.Ordered` and value _V_ has
constraint `any`. All priority queues are min-heap-based (i.e., the smallest key has the highest priority) and allow for
duplicate keys/priorities.

## Contributing

I am just one guy working on this in my free time for fun. So, if you have any suggestions or issues, please feel free
to open an issue or pull request!

## Implementations

Currently, this package provides the following queues.

### Time Complexities

| Type          | findMin | removeMin    | insert       | meld         |
|---------------|---------|--------------|--------------|--------------|
| Binary        | Θ(1)    | Θ(log n)     | Θ(log n)     | Θ(n)         |
| Pairing       | Θ(1)    | O(log n) am. | Θ(1)         | Θ(1)         |
| Skew          | Θ(1)    | O(log n) am. | O(log n) am. | O(log n) am. |
| Skew Binomial | Θ(1)    | Θ(log n)     | Θ(1)         | Θ(log n)     |

### Benchmarks

Apple MacBook Pro (M4 Pro, 24GB RAM)

| Type          | push             | meld          | pop              |
|---------------|------------------|---------------|------------------|
| Binary        | 47.72 ns/op/node | 107.4 ns/node | 523.2 ns/op/node |
| Pairing       | 41.91 ns/op/node | 69.96 ns/node | 792.2 ns/op/node |
| Skew          | 184.3 ns/op/node | 381.5 ns/node | 562.1 ns/op/node |
| Skew Binomial | -                | -             | -                |
