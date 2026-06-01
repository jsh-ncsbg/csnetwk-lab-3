#include <netinet/in.h>
#include <stdint.h>
#include <sys/socket.h>
#include <unistd.h>

#include <atomic>
#include <condition_variable>
#include <cstring>
#include <format>
#include <functional>
#include <iostream>
#include <memory>
#include <mutex>
#include <queue>
#include <thread>
#include <vector>

#define SERVER_PORT 8000

typedef struct {
    char name[64];
    int64_t number;
} payload;

int64_t SERVER_NUMBER = 42;
char SERVER_NAME[64] = "Server of the Sokoboys";

class ThreadPool {
  public:
    explicit ThreadPool(size_t num_threads) {
        for (size_t i = 0; i < num_threads; ++i) {
            threads.emplace_back([this] {
                for (;;) {
                    std::function<void()> task;
                    {
                        std::unique_lock lock(mtx);
                        cv.wait(lock,
                                [this] { return stop || !tasks.empty(); });
                        if (stop && tasks.empty())
                            return;

                        task = std::move(tasks.front());
                        tasks.pop();
                    }

                    task();
                }
            });
        }
    }

    template <typename F, typename... Args>
    void enqueue(F &&f, Args &&...args) {
        auto task = std::bind(std::forward<F>(f), std::forward<Args>(args)...);

        if (stop.load(std::memory_order_relaxed))
            throw std::runtime_error("enqueue on stopped ThreadPool");

        {
            std::scoped_lock lock(mtx);
            tasks.emplace(std::move(task));
        }

        cv.notify_one();
    }

    ~ThreadPool() {
        {
            std::scoped_lock lock(mtx);
            stop = true;
        }
        cv.notify_all();
        for (auto &t : threads) {
            if (t.joinable())
                t.join();
        }
    }

  private:
    std::vector<std::thread> threads;
    std::queue<std::function<void()>> tasks;
    std::mutex mtx;
    std::condition_variable cv;
    std::atomic<bool> stop = false;
};

class Server {
    int server_socket = -1;
    int server_port;

    std::atomic<bool> is_running = false; // shared variable for workers
    std::vector<std::thread> workers;

    ThreadPool thread_pool;

    void handle_client(int client_socket) {
        char buffer[128];
        payload cli_pload;

        ssize_t bytes_received =
            recv(client_socket, buffer, sizeof(buffer) - 1, 0);

        if (!is_running) {
            // exit if we reached here from ::shutdown()
            return;
        }

        if (bytes_received <= 0) {
            if (bytes_received == 0)
                std::cout << "Client disconnected." << std::endl;
            else
                std::cerr << "Failed to receive data from client: "
                          << strerror(errno) << std::endl;
        }

        memcpy(&cli_pload, buffer, sizeof(cli_pload));

        // shutdown server if out of range
        if (cli_pload.number < 1 || cli_pload.number > 100) {
            shutdown();
            close(client_socket);
            return;
        }

        std::cout << std::format("\"{}\" connected to client \"{}\"",
                                 SERVER_NAME, cli_pload.name)
                  << std::endl;
        std::cout << std::format("{} (\"{}\" number) + {} (\"{}\" number) = {}",
                                 SERVER_NUMBER, SERVER_NAME, cli_pload.number,
                                 cli_pload.name,
                                 SERVER_NUMBER + cli_pload.number)
                  << std::endl;

        payload ser_pload;
        memcpy(ser_pload.name, SERVER_NAME, sizeof(payload::name));
        ser_pload.number = SERVER_NUMBER;

        if (send(client_socket, &ser_pload, sizeof(ser_pload), 0) == -1) {
            std::cerr << "Failed to send data to client: " << strerror(errno)
                      << std::endl;
        }

        close(client_socket);
        return;
    }

  public:
    Server(int server_port)
        : server_port(server_port),
          thread_pool(std::thread::hardware_concurrency()) {}

    int start() {
        // create socket
        if ((server_socket = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        sockaddr_in ser_addr = {.sin_family = AF_INET,
                                .sin_port = htons(server_port),
                                .sin_addr = {.s_addr = INADDR_ANY}};

        // bind address
        if (bind(server_socket, (sockaddr *)&ser_addr, sizeof(ser_addr)) ==
            -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        // prepare to accept connection
        if (listen(server_socket, 5) == -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        std::cout << std::format("Server is listening on port {}...",
                                 server_port)
                  << std::endl;

        is_running = true;

        while (is_running) {
            // wait for client connections
            int client_socket = accept(server_socket, nullptr, nullptr);
            if (!is_running)
                break;

            if (client_socket < 0) {
                std::cerr << strerror(errno) << std::endl;
                continue;
            }

            thread_pool.enqueue(&Server::handle_client, this, client_socket);
        }

        return 0;
    }

    void shutdown() {
        if (server_socket != -1) {
            close(server_socket);
        }
        ::shutdown(server_socket, SHUT_RDWR);
        is_running = false;
    }

    ~Server() {
        std::cout << "Shutting down server..." << std::endl;
        shutdown();
    }
};

int main() {
    Server server(SERVER_PORT);
    server.start();

    server.shutdown();

    return 0;
}