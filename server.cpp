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

class Server {
    int server_socket = -1;
    int server_port;

    std::atomic<bool> is_running = false; // shared variable for workers
    std::vector<std::thread> workers;

    void handle_client(int client_socket) {
        char buffer[128];
        payload client_payload;

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
            close(client_socket);
            return;
        }

        memcpy(&client_payload, buffer, sizeof(client_payload));

        std::cout << std::format("\"{}\" connected to client \"{}\"",
                                 SERVER_NAME, client_payload.name)
                  << std::endl;
        std::cout << std::format("{} (\"{}\" number) + {} (\"{}\" number) = {}",
                                 SERVER_NUMBER, SERVER_NAME, client_payload.number,
                                 client_payload.name,
                                 SERVER_NUMBER + client_payload.number)
                  << std::endl;

        payload server_payload;
        memcpy(server_payload.name, SERVER_NAME, sizeof(payload::name));
        server_payload.number = SERVER_NUMBER;

        if (send(client_socket, &server_payload, sizeof(server_payload), 0) == -1) {
            std::cerr << "Failed to send data to client: " << strerror(errno)
                      << std::endl;
        }

        if (client_payload.number < 1 || client_payload.number > 100) {
            shutdown();
        }

        close(client_socket);
        return;
    }

    void worker_thread() {
        while (is_running) {
            int client_socket = accept(server_socket, nullptr, nullptr);
            if (!is_running) {
                if (client_socket >= 0) {
                    close(client_socket);
                }
                break;
            }
            if (client_socket < 0) {
                if (errno == EBADF || errno == EINVAL || errno == ENOTSOCK || errno == ECONNABORTED) {
                    break;
                }
                std::cerr << "Failed to accept connection: " << strerror(errno) << std::endl;
                continue;
            }
            handle_client(client_socket);
        }
    }

  public:
    Server(int server_port) : server_port(server_port) {}

    int start() {
        // create socket
        if ((server_socket = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        // set SO_REUSEADDR and SO_REUSEPORT to prevent "Address already in use" errors
        int opt = 1;
        if (setsockopt(server_socket, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) == -1) {
            std::cerr << "setsockopt SO_REUSEADDR failed: " << strerror(errno) << std::endl;
        }
#ifdef SO_REUSEPORT
        if (setsockopt(server_socket, SOL_SOCKET, SO_REUSEPORT, &opt, sizeof(opt)) == -1) {
            std::cerr << "setsockopt SO_REUSEPORT failed: " << strerror(errno) << std::endl;
        }
#endif

        sockaddr_in server_address = {.sin_family = AF_INET,
                                .sin_port = htons(server_port),
                                .sin_addr = {.s_addr = INADDR_ANY}};

        // bind address
        if (bind(server_socket, (sockaddr *)&server_address, sizeof(server_address)) ==
            -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        // prepare to accept connection
        if (listen(server_socket, 1024) == -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        std::cout << std::format("Server is listening on port {}...",
                                 server_port)
                  << std::endl;

        is_running = true;

        unsigned int num_threads = std::thread::hardware_concurrency();
        if (num_threads == 0) num_threads = 4;

        for (unsigned int i = 0; i < num_threads; ++i) {
            workers.emplace_back(&Server::worker_thread, this);
        }

        for (auto& worker : workers) {
            if (worker.joinable()) {
                worker.join();
            }
        }

        return 0;
    }

    void shutdown() {
        bool expected = true;
        if (is_running.compare_exchange_strong(expected, false)) {
            if (server_socket != -1) {
                int fd = server_socket;
                server_socket = -1;
                ::shutdown(fd, SHUT_RDWR);
                close(fd);
            }
        }
    }

    ~Server() {
        std::cout << "Shutting down server..." << std::endl;
        shutdown();
        for (auto& worker : workers) {
            if (worker.joinable()) {
                worker.join();
            }
        }
    }
};

int main() {
    Server server(SERVER_PORT);
    server.start();

    server.shutdown();

    return 0;
}