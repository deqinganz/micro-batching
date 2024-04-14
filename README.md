# Micro-Batching Library

This is a micro-batching library that processes jobs in batches.

It takes a job and put it in the queue. It calls BatchProcessor by the frequency and batch size set by the user.

## Design

### Queue

The queue is a FIFO queue implemented using a slice. The time complexity of `Enqueue()` is O(1), amortized constant time. The time complexity of `Dequeue(k)` is O(N) because it removes elements from the beginning of a slice.

There are several other possible ways to implement this queue:
- Using a linked list
  - This would make `Dequeue(k)` O(k) but both `Enqueue()` and `Dequeue()` are likely slower than a slice. Probably only suitable when the `Dequeue(k)` is called frequently.
- Using a ring buffer
  - `container/ring` in Go provides a ring buffer implementation. It will have better performance than a slice because it doesn't need to resize the buffer or move elements. But in this implementation, the size of the ring is fixed. If the buffer is full, we need to decide whether to drop the new job or to remove the oldest job.
- Implementing a ring buffer by a slice with pointers to the head and tail
  - `Dequeue(k)` will have better performance because it doesn't need to move elements. It will be more complex to implement and maintain.

When performance comes to play, we should consider the trade-offs between these possible implementations.

### Preprocessing

## Limitations

### Multithreading

This library is not thread-safe, `Enqueue()` and `Dequeue()` can cause race condition. If you want to use this library in a multi-threaded environment, you need to add a lock to the queue.
