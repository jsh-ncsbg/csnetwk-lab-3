import asyncio
import struct
import time
import argparse

# payload struct format: 64s (char name[64]) + q (int64_t number)
PAYLOAD_FORMAT = '<64sq'
PAYLOAD_SIZE = struct.calcsize(PAYLOAD_FORMAT)

async def send_request(client_id, host, port, name, number):
    start_time = time.time()
    try:
        reader, writer = await asyncio.open_connection(host, port)
    except Exception as e:
        return False, 0.0, f"Connect error: {e}"

    try:
        # Prepare payload
        name_bytes = name.encode('utf-8')[:63]
        payload = struct.pack(PAYLOAD_FORMAT, name_bytes, number)
        
        # Send
        writer.write(payload)
        await writer.drain()
        
        # Receive response (72 bytes)
        data = await reader.readexactly(PAYLOAD_SIZE)
        if len(data) == PAYLOAD_SIZE:
            # Unpack response if needed
            srv_name, srv_number = struct.unpack(PAYLOAD_FORMAT, data)
            latency = time.time() - start_time
            return True, latency, None
        else:
            return False, time.time() - start_time, f"Incomplete read: {len(data)} bytes"
    except Exception as e:
        return False, time.time() - start_time, f"Communication error: {e}"
    finally:
        try:
            writer.close()
            await writer.wait_closed()
        except Exception:
            pass

async def worker(queue, host, port, results):
    while True:
        task = await queue.get()
        if task is None:
            queue.task_done()
            break
        client_id, name, number = task
        success, latency, err = await send_request(client_id, host, port, name, number)
        results.append((success, latency, err))
        queue.task_done()

async def main():
    parser = argparse.ArgumentParser(description="Benchmark client for Sokoboys server")
    parser.add_argument("--host", default="127.0.0.1", help="Server host")
    parser.add_argument("--port", type=int, default=8000, help="Server port")
    parser.add_argument("--concurrency", type=int, default=100, help="Number of concurrent clients")
    parser.add_argument("--total", type=int, default=1000, help="Total number of requests to send")
    parser.add_argument("--name", default="Benchmarker", help="Client name prefix")
    parser.add_argument("--number", type=int, default=42, help="Client number (1-100 to avoid shutdown)")
    
    args = parser.parse_args()
    
    print(f"Starting benchmark on {args.host}:{args.port}")
    print(f"Total requests: {args.total}, Concurrency limit: {args.concurrency}")
    
    queue = asyncio.Queue()
    for i in range(args.total):
        queue.put_nowait((i, f"{args.name}-{i}", args.number))
        
    results = []
    
    start_time = time.time()
    
    # Start workers
    workers = []
    for _ in range(args.concurrency):
        workers.append(asyncio.create_task(worker(queue, args.host, args.port, results)))
        
    # Wait for queue to be empty
    await queue.join()
    
    # Stop workers
    for _ in range(args.concurrency):
        queue.put_nowait(None)
    await asyncio.gather(*workers)
    
    total_time = time.time() - start_time
    
    successes = [r for r in results if r[0]]
    failures = [r for r in results if not r[0]]
    
    print("\n--- Benchmark Results ---")
    print(f"Total Time:      {total_time:.2f} seconds")
    print(f"Successful:      {len(successes)} / {args.total} ({len(successes)/args.total*100:.1f}%)")
    print(f"Failed:          {len(failures)} / {args.total} ({len(failures)/args.total*100:.1f}%)")
    if total_time > 0:
        print(f"Throughput:      {len(successes)/total_time:.2f} successful req/sec")
        
    if successes:
        latencies = [r[1] for r in successes]
        print(f"Min Latency:     {min(latencies)*1000:.2f} ms")
        print(f"Max Latency:     {max(latencies)*1000:.2f} ms")
        print(f"Avg Latency:     {(sum(latencies)/len(latencies))*1000:.2f} ms")
        
    if failures:
        print("\nSample Errors:")
        err_counts = {}
        for _, _, err in failures:
            err_counts[err] = err_counts.get(err, 0) + 1
        for err, count in list(err_counts.items())[:5]:
            print(f"  - {err} ({count} occurrences)")

if __name__ == "__main__":
    asyncio.run(main())
